package olduar

import (
	"strconv"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"math/rand"
	"time"
)

// Create Game State from save / scratch

var AllRooms map[string]*Room = make(map[string]*Room)

func GetRoomList() []string {
	roomList := make([]string,len(AllRooms))
	roomId := 0
	for room := range AllRooms {
		roomList[roomId] = room
		roomId++
	}
	return roomList
}

func CreateRoomWithName(name string) *Room {
	room := CreateRoomFromSave(name+".json")
	if(room == nil) {
		loc := CreateLocationFromRegion("start")
		room = &Room{
			Id: name,
			CurrentLocation: loc,
			StartingLocation: loc,
			Players: make(Players,0),
		}
		room.CurrentLocation.Visit()
		room.Prepare()
	}
	return room
}

func CreateRoomFromScratch() *Room {
	return CreateRoomWithName("game_" + strconv.Itoa(len(AllRooms)+1))
}

func CreateRoomFromSave(filename string) *Room {
	room := Room{}
	data, err := ioutil.ReadFile("save/rooms/"+filename);
	if(err == nil) {
		err := json.Unmarshal(data, &room)
		if(err == nil) {
			room.Prepare()
			return &room
		}
	}
	return nil
}

// Message Object

type MessageObjects []*MessageObject

type MessageObject struct {
	Id int64					`json:"id"`
	Message string				`json:"text"`
	IgnoredPlayer *Player		`json:"-"`
	OnlyForPlayer *Player		`json:"-"`
}

// Command object

type Command struct {
	Player *Player
	Command, Parameter string
	Response chan []byte
}

type Response struct {
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
	History MessageObjects 		`json:"history"`
	Exits map[string]string		`json:"exits"`
	Actions map[string]string	`json:"actions"`
	Items []ResponseItem		`json:"items,omitempty"`
}

type ResponseItem struct {
	Id *string `json:"id"`
	Name *string `json:"name"`
	Description *string `json:"desc"`
}

type ResponseItemDetail struct {
	//All properties must be pointers as we are just reusing something from item template
	Name *string `json:"name"`
	Description *string `json:"desc"`
}

// Game State and functions

type Rooms []*Room

type Room struct {
	Id string 					`json:"id"`
	CurrentLocation *Location 	`json:"-"`
	StartingLocation *Location 	`json:"location"`
	Players Players				`json:"-"`
	History MessageObjects		`json:"-"`
	LastMessageId int64			`json:"message_count"`

	queue chan *Command
	voting bool
	votingTime time.Time
}

func (room *Room) Save() {
	data, err := json.Marshal(room)
	if(err == nil) {
		err = ioutil.WriteFile("save/rooms/"+room.Id+".json", data, 0644)
		if err != nil {
			fmt.Println("Failed to save room \""+room.Id+"\":",err)
		} else {
			fmt.Println("Room \""+room.Id+"\" has been saved")
		}
	} else {
		fmt.Println("Failed to serialize room \""+room.Id+"\":",err)
	}
}

func (room *Room) cycleLocations(location *Location) {
	//Items
	for _, item := range location.Items {
		item.Load()
	}

	//Check location
	if(location.Current) {
		room.CurrentLocation = location
	}
	for _, exit := range location.Exits {
		exit.Target.Parent = location
		room.cycleLocations(exit.Target)
	}
}

func (room *Room) Prepare() {
	//Put room to list of rooms
	AllRooms[room.Id] = room

	//Check for current location
	if(room.CurrentLocation == nil) {
		room.cycleLocations(room.StartingLocation)
	}

	//Prepare variables
	room.queue = make(chan *Command,0)
	room.voting = false

	//Message worker
	go func(){
		for {
			cmd := <- room.queue
			var resp []byte = nil

			//Check for voting timeout
			if(room.voting && cmd.Command != "go" && room.votingTime.Before(time.Now())) {
				room.CheckVoting()
			}

			//Process commands
			switch(cmd.Command) {
			case "save":
				room.Save()
				resp = []byte("null")
			case "look":
				resp = room.GetPlayerResponse(cmd.Player)
			case "do":
				if(cmd.Parameter != "") {
					room.CurrentLocation.DoAction(room,cmd.Player,cmd.Parameter)
				}
				resp = room.GetPlayerResponse(cmd.Player)
			case "go":
				if(cmd.Parameter != "") {
					room.GoTo(cmd.Parameter,cmd.Player)
				}
				resp = room.GetPlayerResponse(cmd.Player)
			case "inventory":
				inventory := make([]ResponseItem,len(cmd.Player.Inventory))
				for index, item := range cmd.Player.Inventory {
					inventory[index] = item.GenerateResponse()
				}
				resp, _ = json.Marshal(inventory)
			case "inspect":
				item := cmd.Player.Inventory.Get(cmd.Parameter)
				if(item != nil) {
					resp, _ = json.Marshal(item.Attributes.Response)
				} else {
					item := room.CurrentLocation.Items.Get(cmd.Parameter)
					if(item != nil) {
						resp, _ = json.Marshal(item.Attributes.Response)
					}
				}
			case "pickup":
				if(cmd.Parameter != "" && cmd.Player.Pickup(cmd.Parameter)) {
					resp = room.GetPlayerResponse(cmd.Player)
				}
			case "drop":
				if(cmd.Player.Drop(cmd.Parameter)) {
					resp = room.GetPlayerResponse(cmd.Player)
				}
			case "use":
				if(cmd.Player.Use(cmd.Parameter)) {
					resp = room.GetPlayerResponse(cmd.Player)
				}
			}
			if(resp == nil) {
				resp = []byte("null")
			}
			cmd.Response <- resp
		}
	}()

	fmt.Println("Game room \""+room.Id+"\" is ready")
}

func (room *Room) AddMessage(message *MessageObject) {
	room.LastMessageId++
	message.Id = room.LastMessageId
	room.History = append(room.History, message)
}

func (room *Room) TellAll(str string) {
	room.AddMessage(&MessageObject{Message:str})
}

func (room *Room) TellAllExcept(str string, player *Player) {
	room.AddMessage(&MessageObject{Message:str,IgnoredPlayer:player})
}

func (room *Room) Tell(str string, player *Player) {
	room.AddMessage(&MessageObject{Message:str,OnlyForPlayer:player})
}

func (room *Room) GetPlayerResponse(player *Player) []byte {
	from := player.LastResponseId

	res := Response{
		Name: room.CurrentLocation.Name,
		Description: room.CurrentLocation.Description,
		History: make(MessageObjects,0),
		Exits: make(map[string]string),
		Actions: make(map[string]string),
		Items: make([]ResponseItem,len(room.CurrentLocation.Items)),
	}

	//Append items
	for index, item := range room.CurrentLocation.Items {
		res.Items[index] = item.GenerateResponse()
	}

	//Append history
	for _, entry := range room.History {
		if(entry.Id > from && ((entry.IgnoredPlayer == nil && entry.OnlyForPlayer == nil) || (entry.IgnoredPlayer != player && entry.OnlyForPlayer == nil) || (entry.IgnoredPlayer == nil && entry.OnlyForPlayer == player))) {
			res.History = append(res.History,entry)
			player.LastResponseId = entry.Id
		}
	}

	//Append exits
	for _, exit := range room.CurrentLocation.Exits {
		res.Exits[exit.Id] = exit.Target.DescriptionShort
	}
	if(room.CurrentLocation.Parent != nil) {
		res.Exits["back"] = room.CurrentLocation.Parent.DescriptionShort
	}

	//Append actions
	for _, action := range room.CurrentLocation.Actions {
		if(action.Charges != 0) {
			res.Actions[action.Id] = action.Description
		}
	}

	//Prepare JSON
	data, error := json.Marshal(res)
	if(error != nil) {
		return nil
	}
	return data
}

func (room *Room) DoAction(player *Player, action *Action) {
	//Check for requirements
	if(len(action.Requirements)>0 && (action.Charges == -1 || action.Charges > 0)) {
		for _, requirement := range action.Requirements {
			switch(requirement.Type){
			case "item":
				if(!player.Owns(requirement.Value)) {
					if(requirement.ErrorMessage != "") {
						room.Tell(requirement.ErrorMessage,player)
					}
					return
				}
			}
		}
	}

	//Charges
	if(action.Charges > -1) {
		if(action.Charges > 0) {
			action.Charges--;
		} else {
			return //No charges left = no loot
		}
	}

	//Do actual action
	action.Do(room,player)

	//Messages
	message, found := action.Config["msg_player"]
	if(found) {
		room.Tell(AppendVariablesToString(message.(string),player,action.Config),player)
	}
	message, found = action.Config["msg_party"]
	if(found) {
		room.TellAllExcept(AppendVariablesToString(message.(string),player,action.Config),player)
	}
	message, found = action.Config["msg"]
	if(found) {
		room.TellAll(AppendVariablesToString(message.(string),player,action.Config))
	}
}

func (room *Room) Travel(location *Location) {
	room.voting = false

	//Reset voting state
	for _, player := range room.Players {
		player.VotedLocation = nil
	}

	//Tell players
	room.TellAll("You went to "+location.DescriptionShort)

	//Add "back" exit if location was not visited before
	if(location.Visit()) {
		room.CurrentLocation.Current = false;
		location.Current = true;
		location.Parent = room.CurrentLocation
	}

	//Set new location
	room.CurrentLocation = location
}

func (room *Room) CheckVoting() {
	//Voting not in progress? skip
	if(!room.voting) {
		return
	}

	//Check if all players voted
	proceedToNewLocation := true
	votes := make(map[*Location]int)
	for _, player := range room.Players {
		if(player.VotedLocation == nil) {
			proceedToNewLocation = false
		} else {
			votes[player.VotedLocation]++
		}
	}

	if(room.votingTime.Before(time.Now())) {
		proceedToNewLocation = true
	}

	//Select voted location
	if(proceedToNewLocation) {
		var votedLocation *Location = nil
		votedLocationVotes := 0
		for location, votes := range votes {
			if(votes > votedLocationVotes || (votes == votedLocationVotes && rand.Float64() > 0.5)) {
				votedLocation = location
				votedLocationVotes = votes
			}
		}
		//Travel
		room.Travel(votedLocation)
	}
}

func (room *Room) GoTo(way string, player *Player) {
	oldLocation := room.CurrentLocation
	var newLocation *Location = nil

	if(way == "back") {
		newLocation = room.CurrentLocation.Parent
	} else {
		for _, exit := range oldLocation.Exits {
			if(exit.Id == way) {
				newLocation = exit.Target
			}
		}
	}

	if(newLocation != nil) {
		if(len(room.Players) == 1) {
			//One player: instant travel
			room.Travel(newLocation)
		} else {
			//More players: voting
			room.voting = true
			room.votingTime = time.Now().Add(time.Second * 10)
			player.VotedLocation = newLocation
			//Count players who voted
			votes, maxVotes := 0, len(room.Players)
			for _, player := range room.Players {
				if(player.VotedLocation != nil) {
					votes++
				}
			}
			//Send voting messages
			voteStatus := "("+strconv.Itoa(votes)+" of "+strconv.Itoa(maxVotes)+" players voted)"
			room.TellAllExcept(player.Name+" wants to go to "+newLocation.DescriptionShort+" "+voteStatus,player)
			room.Tell("You want to go to "+newLocation.DescriptionShort+" "+voteStatus,player)
			//Check if voting has been completed
			room.CheckVoting()
		}
	}
}

func (room *Room) Leave(player *Player) {
	if(player.Room == room) {
		player.Room = nil
		count := 0
		newPlayers := make(Players,len(room.Players)-1)
		for _, p := range room.Players {
			if(p != player) {
				newPlayers[count] = p
				count++
			}
		}
		room.Players = newPlayers
		player.Save()
	}
	if(len(room.Players) == 0) {
		room.Save()
		delete(AllRooms,room.Id)
	}
}

func (room *Room) Join(player *Player) {
	room.Players = append(room.Players,player)
	player.Room = room
	player.LastResponseId = room.LastMessageId
}
