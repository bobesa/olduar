package olduar

import (
	"strconv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
	"math/rand"
	"time"
)

// Create Game State from save / scratch

var AllGameStates GameStates = make(GameStates,0)

func CreateGameStateFromName(name string) *GameState {
	gs := GameState{
		Id: name,
		CurrentLocation: CreateLocationFromRegion("start"),
		Players: make(Players,0),
	}
	gs.CurrentLocation.Visit()
	gs.Prepare()
	return &gs
}

func CreateGameStateFromScratch() *GameState {
	return CreateGameStateFromName("game_" + strconv.Itoa(len(AllGameStates)+1))
}

func CreateGameStateFromSave(filename string) *GameState {
	gs := GameState{}
	data, err := ioutil.ReadFile("save/"+filename);
	if(err == nil) {
		err := json.Unmarshal(data, &gs)
		if(err == nil) {
			gs.Prepare()
			return &gs
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

type GameStates []*GameState

type GameState struct {
	Id string 					`json:"id"`
	CurrentLocation *Location 	`json:"location"`
	Players Players				`json:"-"`
	History MessageObjects		`json:"-"`
	LastMessageId int64			`json:"message_count"`

	queue chan *Command
	voting bool
	votingTime time.Time
}

func (state *GameState) Prepare() {

	//Prepare variables
	state.queue = make(chan *Command,0)
	state.voting = false

	//Append REST api
	apiPath := "/"+state.Id+"/"
	apiPathLen := len(apiPath)
	MainServerMux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request){
			//Authentication
			authToken, found := r.Header["Authorization"]
			if(!found) {
				http.NotFound(w,r)
				return
			}

			//Search for player in state
			player, active := ActivePlayers[authToken[0]]
			if(!active) {
				http.NotFound(w,r)
				return
			}

			//Process command
			params := strings.Split(r.URL.Path[apiPathLen:],"/")
			if(len(params)>0) {
				resp := make(chan []byte)
				command := Command{Player:player, Command: strings.ToLower(params[0]), Response: resp}
				if(len(params)>1) {
					command.Parameter = strings.ToLower(params[1])
				}
				state.queue <- &command
				data := <- resp
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
			}
		})

	//Message worker
	go func(){
		for {
			cmd := <- state.queue
			var resp []byte = nil

			//Check for voting timeout
			if(state.voting && cmd.Command != "go" && state.votingTime.Before(time.Now())) {
				state.CheckVoting()
			}

			//Process commands
			switch(cmd.Command) {
			case "look":
				resp = state.GetPlayerResponse(cmd.Player)
			case "do":
				if(cmd.Parameter != "") {
					state.CurrentLocation.DoAction(state,cmd.Player,cmd.Parameter)
				}
				resp = state.GetPlayerResponse(cmd.Player)
			case "go":
				if(cmd.Parameter != "") {
					state.GoTo(cmd.Parameter,cmd.Player)
				}
				resp = state.GetPlayerResponse(cmd.Player)
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
					item := state.CurrentLocation.Items.Get(cmd.Parameter)
					if(item != nil) {
						resp, _ = json.Marshal(item.Attributes.Response)
					}
				}
			case "pickup":
				if(cmd.Parameter != "" && cmd.Player.Pickup(cmd.Parameter)) {
					resp = state.GetPlayerResponse(cmd.Player)
				}
			case "drop":
				if(cmd.Player.Drop(cmd.Parameter)) {
					resp = state.GetPlayerResponse(cmd.Player)
				}
			case "use":
				if(cmd.Player.Use(cmd.Parameter)) {
					resp = state.GetPlayerResponse(cmd.Player)
				}
			}
			if(resp == nil) {
				resp = []byte("null")
			}
			cmd.Response <- resp
		}
	}()

	fmt.Println("Game room \""+state.Id+"\" is ready")
}

func (state *GameState) AddMessage(message *MessageObject) {
	state.LastMessageId++
	message.Id = state.LastMessageId
	state.History = append(state.History, message)
}

func (state *GameState) TellAll(str string) {
	state.AddMessage(&MessageObject{Message:str})
}

func (state *GameState) TellAllExcept(str string, player *Player) {
	state.AddMessage(&MessageObject{Message:str,IgnoredPlayer:player})
}

func (state *GameState) Tell(str string, player *Player) {
	state.AddMessage(&MessageObject{Message:str,OnlyForPlayer:player})
}

func (state *GameState) GetPlayerResponse(player *Player) []byte {
	from := player.LastResponseId

	res := Response{
		Name: state.CurrentLocation.Name,
		Description: state.CurrentLocation.Description,
		History: make(MessageObjects,0),
		Exits: make(map[string]string),
		Actions: make(map[string]string),
		Items: make([]ResponseItem,len(state.CurrentLocation.Items)),
	}

	//Append items
	for index, item := range state.CurrentLocation.Items {
		res.Items[index] = item.GenerateResponse()
	}

	//Append history
	for _, entry := range state.History {
		if(entry.Id > from && ((entry.IgnoredPlayer == nil && entry.OnlyForPlayer == nil) || (entry.IgnoredPlayer != player && entry.OnlyForPlayer == nil) || (entry.IgnoredPlayer == nil && entry.OnlyForPlayer == player))) {
			res.History = append(res.History,entry)
			player.LastResponseId = entry.Id
		}
	}

	//Append exits
	for _, exit := range state.CurrentLocation.Exits {
		res.Exits[exit.Id] = exit.Target.DescriptionShort
	}

	//Append actions
	for _, action := range state.CurrentLocation.Actions {
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

func (state *GameState) DoAction(player *Player, action *Action) {
	//Check for requirements
	if(len(action.Requirements)>0 && (action.Charges == -1 || action.Charges > 0)) {
		for _, requirement := range action.Requirements {
			switch(requirement.Type){
			case "item":
				if(!player.Owns(requirement.Value)) {
					if(requirement.ErrorMessage != "") {
						state.Tell(requirement.ErrorMessage,player)
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
	action.Do(state,player)

	//Messages
	message, found := action.Config["msg_player"]
	if(found) {
		state.Tell(AppendVariablesToString(message.(string),player,action.Config),player)
	}
	message, found = action.Config["msg_party"]
	if(found) {
		state.TellAllExcept(AppendVariablesToString(message.(string),player,action.Config),player)
	}
	message, found = action.Config["msg"]
	if(found) {
		state.TellAll(AppendVariablesToString(message.(string),player,action.Config))
	}
}

func (state *GameState) Travel(location *Location) {
	state.voting = false

	//Reset voting state
	for _, player := range state.Players {
		player.VotedLocation = nil
	}

	//Tell players
	state.TellAll("You went to "+location.DescriptionShort)

	//Add "back" exit if location was not visited before
	if(location.Visit()) {
		location.Exits = append(
			location.Exits,
			LocationExit{
				Id:"back",
				Target: state.CurrentLocation,
			},
		)
	}

	//Set new location
	state.CurrentLocation = location
}

func (state *GameState) CheckVoting() {
	//Voting not in progress? skip
	if(!state.voting) {
		return
	}

	//Check if all players voted
	proceedToNewLocation := true
	votes := make(map[*Location]int)
	for _, player := range state.Players {
		if(player.VotedLocation == nil) {
			proceedToNewLocation = false
		} else {
			votes[player.VotedLocation]++
		}
	}

	if(state.votingTime.Before(time.Now())) {
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
		state.Travel(votedLocation)
	}
}

func (state *GameState) GoTo(way string, player *Player) {
	oldLocation := state.CurrentLocation
	var newLocation *Location = nil
	for _, exit := range oldLocation.Exits {
		if(exit.Id == way) {
			newLocation = exit.Target
		}
	}

	if(newLocation != nil) {
		if(len(state.Players) == 1) {
			//One player: instant travel
			state.Travel(newLocation)
		} else {
			//More players: voting
			state.voting = true
			state.votingTime = time.Now().Add(time.Second * 10)
			player.VotedLocation = newLocation
			//Count players who voted
			votes, maxVotes := 0, len(state.Players)
			for _, player := range state.Players {
				if(player.VotedLocation != nil) {
					votes++
				}
			}
			//Send voting messages
			voteStatus := "("+strconv.Itoa(votes)+" of "+strconv.Itoa(maxVotes)+" players voted)"
			state.TellAllExcept(player.Name+" wants to go to "+newLocation.DescriptionShort+" "+voteStatus,player)
			state.Tell("You want to go to "+newLocation.DescriptionShort+" "+voteStatus,player)
			//Check if voting has been completed
			state.CheckVoting()
		}
	}
}

func (state *GameState) Leave(player *Player) {
	count := 0
	newPlayers := make(Players,len(state.Players)-1)
	for _, p := range state.Players {
		if(p != player) {
			newPlayers[count] = p
			count++
		}
	}
}

func (state *GameState) Join(player *Player) {
	state.Players = append(state.Players,player)
	player.GameState = state
}
