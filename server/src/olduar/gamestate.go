package olduar

import (
	"strconv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"fmt"
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
}

func (state *GameState) Prepare() {
	//Append REST api
	apiPath := "/"+state.Id+"/"
	apiPathLen := len(apiPath)
	MainServerMux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request){
			//TODO: Proper handling of authentication

			//Search for player in state
			var player *Player = nil
			for _, p := range state.Players {
				//if(p.Username ==)
				player = p
			}
			if(player == nil) {
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

	//Prepare channels
	state.queue = make(chan *Command,0)

	//Message worker
	go func(){
		for {
			cmd := <- state.queue
			var resp []byte = nil
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
					state.GoTo(cmd.Parameter)
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
					resp, _ = json.Marshal(item.Attributes.Response) //decrease index as we are increasing it in inventory command
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

	fmt.Println("Game room \""+state.Id+"\" is ready.")
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
		if(entry.Id > from && (entry.IgnoredPlayer != player || entry.OnlyForPlayer == player)) {
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

func (state *GameState) GoTo(way string) {
	oldLocation := state.CurrentLocation
	var newLocation *Location = nil
	for _, exit := range oldLocation.Exits {
		if(exit.Id == way) {
			newLocation = exit.Target
		}
	}

	if(newLocation != nil) {
		state.TellAll("You went to "+newLocation.DescriptionShort)
		if(newLocation.Visit()) {
			newLocation.Exits = append(
				newLocation.Exits,
				LocationExit{
					Id:"back",
					Target: oldLocation,
				},
			)
		}
		state.CurrentLocation = newLocation
	}
}

func (state *GameState) Describe() {
	state.CurrentLocation.Describe()
}

func (state *GameState) Leave(player *Player) {

}

func (state *GameState) Join(player *Player) {
	state.Players = append(state.Players,player)
	player.GameState = state
}
