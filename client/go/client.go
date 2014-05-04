package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"os"
	"time"
)

const (
	VERSION = "0.02a"
	SERVER = "http://localhost:8080"
	GAME = "test"

	COLOR_BLACK = 30
	COLOR_RED = 31
	COLOR_GREEN = 32
	COLOR_YELLOW = 33
	COLOR_BLUE = 34
	COLOR_PURPLE = 35
	COLOR_CYAN = 36
	COLOR_GREY = 37

	COLOR_BG_BLACK = 40
	COLOR_BG_RED = 41
	COLOR_BG_GREEN = 42
	COLOR_BG_YELLOW = 43
	COLOR_BG_BLUE = 44
	COLOR_BG_PURPLE = 45
	COLOR_BG_CYAN = 46
	COLOR_BG_GREY = 47
)

var Username, Password string = "", ""

// Types

type LogEvent struct {
	Type int					`json:"type"`
	Data string					`json:"data"`
}

type Inventory []Item
type Item struct {
	Id string 					`json:"id"`
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
}

func (inventory Inventory) Print() {
	if(len(inventory) > 0) {
		fmt.Println(formatBold("Your inventory"))
		for _, item := range inventory {
			fmt.Println(item.Name+" ("+formatCommand(item.Id)+") "+formatInfo(item.Description))
		}
	} else {
		fmt.Println(formatBold("Your inventory is empty"))
	}
}

type ItemDetail struct {
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
}

func (item *ItemDetail) Print() {
	if(item != nil) {
		fmt.Println(formatBold(item.Name))
		fmt.Println(formatInfo(item.Description))
	} else {
		fmt.Println(formatBold("You don't own this item"))
	}
}

type Npcs []Npc
type Npc struct {
	Id string 					`json:"id"`
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
	Health float64				`json:"health`
	HealthMax float64			`json:"health_max`
	Friendly bool				`json:"friendly`
}

func (npc *Npc) Print() {
	if(!npc.Friendly) {
		fmt.Println(formatBold(npc.Name + " (" + formatColor(fmt.Sprint(npc.Health),COLOR_RED) + "/" + formatColor(fmt.Sprint(npc.HealthMax),COLOR_RED) + ")") )
	} else {
		fmt.Println(formatBold(npc.Name))
	}
}

func (npcs Npcs) InCombat() bool {
	for _, npc := range npcs {
		if(!npc.Friendly) {
			return true
		}
	}
	return false
}

type Room struct {
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
	Log []*LogEvent		 		`json:"log"`
	Exits map[string]string		`json:"exits"`
	Actions map[string]string	`json:"actions"`
	Items Inventory				`json:"items"`
	Npcs Npcs					`json:"npcs"`
}

func (room *Room) Print() {
	//History events
	if(len(room.Log)>0) {
		for _, event := range room.Log {
			//TODO: Handle color by event .Type
			fmt.Println(event.Data)
		}
		PrintLine()
	}

	//Location info
	fmt.Println(formatBold(room.Name))
	fmt.Println(formatInfo(room.Description))

	//Items
	count := len(room.Items)
	if(count > 0) {
		if(count == 1) {
			fmt.Print("Item on ground: ")
		} else {
			fmt.Print("Items on ground: ")
		}
		for index, item := range room.Items {
			fmt.Print(formatCommand(item.Id)+" ("+item.Name+")")
			if(count-2 == index) {
				fmt.Print(" & ")
			} else if(count-1 != index) {
				fmt.Print(", ")
			}
		}
		fmt.Println();
	}

	//Actions
	count = len(room.Actions)
	if(count > 0) {
		if(count == 1) {
			fmt.Print("There is only 1 action: ")
		} else {
			fmt.Print("Possible actions are: ")
		}
		actionId := 0
		for action, desc := range room.Actions {
			fmt.Print(formatCommand(action)+" ("+desc+")")
			if(count-2 == actionId) {
				fmt.Print(" & ")
			} else if(count-1 != actionId) {
				fmt.Print(", ")
			}
			actionId++
		}
		fmt.Println();
	}

	//Exits
	count = len(room.Exits)
	if(count > 0) {
		if(count == 1) {
			fmt.Print("There is only 1 exit: ")
		} else {
			fmt.Print("Directions are: ")
		}
		exitId := 0
		for command, desc := range room.Exits {
			fmt.Print(formatCommand(command)+" ("+desc+")")
			if(count-2 == exitId) {
				fmt.Print(" & ")
			} else if(count-1 != exitId) {
				fmt.Print(", ")
			}
			exitId++
		}
		fmt.Println();
	} else {
		if(room.Npcs.InCombat()) {
			fmt.Println("There is combat going on around you")
			for _, npc := range room.Npcs {
				npc.Print()
			}
		} else {
			fmt.Println("There is no exit from this place")
		}
	}


	//Npcs
	if(!room.Npcs.InCombat()) {
		count = len(room.Npcs)
		if(count > 0) {
			if(count == 1) {
				fmt.Print("There is only 1 npc: ")
			} else {
				fmt.Print("Npcs around you are: ")
			}
			for index, npc := range room.Npcs {
				fmt.Print(formatCommand(npc.Id)+" ("+npc.Name+")")
				if(count-2 == index) {
					fmt.Print(" & ")
				} else if(count-1 != index) {
					fmt.Print(", ")
				}
			}
			fmt.Println();
		}
	}
}

type AttributeValue struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
type AttributeList map[string]AttributeValue

func (list AttributeList) Print() {
	for stat, values := range list {
		if(values.Min == values.Max) {
			fmt.Println(formatBold(stat) + ": " + formatInfo(fmt.Sprint(values.Min)))
		} else {
			fmt.Println(formatBold(stat) + ": " + formatInfo(fmt.Sprint(values.Min)) + formatBold(" - ") + formatInfo(fmt.Sprint(values.Max)))
		}
	}
}

//Authorization + Registration

func RequireAuth() {
	fmt.Println(formatBoldColor("Wrong Username (",COLOR_RED)+formatBold(Username)+formatBoldColor(") or Password ",COLOR_RED)+" ("+formatInfo("enter same credentials again to register new user")+")")
	oUser, oPass := Username, Password
	Username, Password = "", ""
	for(Username == "") {
		fmt.Print(formatBold("Enter Username:")+" \033[32m")
		fmt.Scanln(&Username)
		fmt.Print("\033[0m")
	}
	for(Password == "") {
		fmt.Print(formatBold("Enter Password:")+" \033[32m")
		fmt.Scanln(&Password)
		fmt.Print("\033[0m")
	}

	//Handle registration
	if(Password == oPass && Username == oUser) {
		if(string(Fetch("POST","register","")) != "true") {
			fmt.Println(formatBoldColor("Username is invalid or already taken. Please try again.",COLOR_RED))
			PrintLine()
			RequireAuth()
		} else {
			fmt.Println(formatBoldColor("Registration was successful.",COLOR_GREEN))
		}
	}

	PrintLine()
}


func Fetch(method string, command string, param string) []byte {
	return FetchWithBody(method,command,param,"")
}

func FetchWithBody(method string, command string, param string, body string) []byte {
	//Check for param
	if(param != "") {
		command = command +"/"
	}

	//Do the request
	client := &http.Client{}
	req, err := http.NewRequest(method, SERVER+"/api/"+command+param, strings.NewReader(body))
	req.SetBasicAuth(Username, Password)
	resp, err := client.Do(req)
	if(err != nil) {
		fmt.Println(formatBoldColor("Network issue: Retrying in 5 seconds",COLOR_RED))
		time.Sleep(time.Second*5)
		return FetchWithBody(method,command,param,body)
	}

	if(resp.StatusCode == 404) {
		RequireAuth()
		return FetchWithBody(method,command,param,body)
	}

	//Read body
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return nil
	}
	return data
}

func PrintLine() {
	fmt.Println("----------------------------------------")
}

func PrintHelp() {
	fmt.Println("\t"+formatCommand("help")+" \t\t\t\t\t\t List of available commands")

	fmt.Println("\t"+formatCommand("look")+" \t\t\t\t\t\t See current location description etc.")
	fmt.Println("\t"+formatCommand("go")+" [direction] \t\t\t\t Walk to a new place (example: "+formatInfo("go west")+")")
	fmt.Println("\t"+formatCommand("do")+" [task] \t\t\t\t\t Do action (example: "+formatInfo("do drink")+")")

	fmt.Println("\t"+formatCommand("inventory")+" \t\t\t\t\t See your inventory")
	fmt.Println("\t"+formatCommand("pickup")+" [item] \t\t\t\t Pickup object on ground (example: "+formatInfo("pickup shield")+")")
	fmt.Println("\t"+formatCommand("drop")+" [item] \t\t\t\t Drop object from inventory on ground (example: "+formatInfo("drop axe")+")")
	fmt.Println("\t"+formatCommand("inspect")+" [item] \t\t\t\t Displays complete info about item (example: "+formatInfo("inspect fishing_pole")+")")
	fmt.Println("\t"+formatCommand("equip")+" [item] \t\t\t\t Equips item (example: "+formatInfo("equip fishing_pole")+")")
	fmt.Println("\t"+formatCommand("stats")+" \t\t\t\t\t\t Current player's stats")

	fmt.Println("\t"+formatCommand("tell")+" [player] [message] \t Send message to player (example: "+formatInfo("tell noam Hi!")+")")
	fmt.Println("\t"+formatCommand("say")+" [message] \t\t\t\t Send message to party (example: "+formatInfo("say Hello guys")+")")
	fmt.Println("\t"+formatCommand("rename")+" [name] \t\t\t\t Rename yourself to new name (example: "+formatInfo("rename Bugmaster3000")+")")

	fmt.Println("\t"+formatCommand("join")+" [room] \t\t\t\t Join a room (example: "+formatInfo("join room_of_horrors")+")")
	fmt.Println("\t"+formatCommand("leave")+" \t\t\t\t\t\t Leave current room")
	fmt.Println("\t"+formatCommand("rooms")+" \t\t\t\t\t\t Show list of rooms on server")

	fmt.Println("\t"+formatCommand("players")+" \t\t\t\t\t Show list of players on server")
	fmt.Println("\t"+formatCommand("party")+" \t\t\t\t\t\t Show list of current players in your party")
}

func Process(command string, param string, param2 string) {
	if(command == "") {
		return
	}
	PrintLine()
	switch(command){
	case "inventory":
		inventory := Inventory{}
		json.Unmarshal(Fetch("GET",command,param),&inventory)
		inventory.Print()

	case "inspect":
		var item *ItemDetail = nil
		json.Unmarshal(Fetch("GET",command,param),&item)
		item.Print()

	case "stats","equip":
		var list *AttributeList = nil
		json.Unmarshal(Fetch("GET",command,param),&list)
		if(list != nil) {
			list.Print()
		} else if(command == "stats") {
			fmt.Println(formatBold("You posses no stats"))
		} else if(command == "equip") {
			fmt.Println(formatBoldColor("You must own the item to equip it!",COLOR_RED))
		}

	case "pickup", "drop", "use":
		if(param != "") {
			var room *Room = nil
			json.Unmarshal(Fetch("POST",command,param),&room)
			if(room != nil) {
				room.Print()
			} else {
				fmt.Print(formatBold("You cannot " + command + " ") + formatCommand(param))
				if(command == "pickup") {
					fmt.Println(formatBold(" because it's not here"))
				} else {
					fmt.Println(formatBold(" because it's not in your inventory"))
				}
			}
		} else {
			fmt.Println(formatBold("You need to say what you want to " + command))
		}

	case "rename": //POST "rename" postbody name
		if(param != "") {
			FetchWithBody("POST",command,"",param)
		} else {
			fmt.Println(formatBold("You need to enter a new name!"))
		}

	case "join": //POST "join/{room}"
		var room *Room = nil
		json.Unmarshal(Fetch("POST",command,param),&room)
		if(room == nil) {
			fmt.Println(formatBold("You cannot create this room, maximum amount of rooms reached."))
		} else {
			room.Print()
		}

	case "tell":
		if(string(FetchWithBody("POST",command,param,param2)) != "true") {
			fmt.Println(formatBoldColor("Player "+param+" could not get your message",COLOR_RED))
		}

	case "say":
		if(string(FetchWithBody("POST",command,"",param)) != "true") {
			fmt.Println(formatBoldColor("Sending of message failed",COLOR_RED))
		}

	case "leave":
		if(strings.Index(string(Fetch("POST",command,param)),"[") == 0) {
			fmt.Println(formatBold("You have left the room"))
		} else {
			fmt.Println(formatBold("You cannot leave the room right now"))
		}

	case "rooms", "players", "party":
		list := []string{}
		json.Unmarshal(Fetch("GET",command,param),&list)
		switch(command){
		case "rooms","leave":
			if(len(list) == 0) {
				fmt.Println(formatBold("No rooms are currently active"))
			} else {
				fmt.Println(formatBold("List of active rooms"))
			}
		case "players":
			if(len(list) == 0) {
				fmt.Println(formatBold("No players are currently playing"))
			} else {
				fmt.Println(formatBold("List of players currently playing"))
			}
		case "party":
			if(len(list) == 0) {
				fmt.Println(formatBold("No players are in your party"))
			} else {
				fmt.Println(formatBold("List of players in your party"))
			}
		}
		for _, entry := range list {
			fmt.Println(formatBold("- ")+formatCommand(entry))
		}

	case "do", "go", "look":
		var room *Room = nil
		method := "POST"
		if(command == "look") {
			method = "GET"
		}
		json.Unmarshal(Fetch(method,command,param),&room)
		if(room != nil) {
			room.Print()
		} else {
			fmt.Println(formatBold("You need to join the room first!"));
			fmt.Println(formatBold("Type \"")+formatCommand("rooms")+formatBold("\" to get the list of rooms you can join."))
			fmt.Println(formatBold("Type \"")+formatCommand("join ")+formatBoldColor("{name of room}",COLOR_PURPLE)+formatBold("\" to join selected room, or create unexisting room with ")+formatBoldColor("{name of room}",COLOR_PURPLE)+formatBold("."))
		}

	case "help":
		fmt.Println(formatBold("Available commands"))
		PrintHelp()
	default:
		fmt.Println(formatBold("Unknown command ")+formatCommand(command)+formatBold(", see available commands below"))
		PrintHelp()
	}
}

//Formatting
func formatBold(str string) string {
	return "\033[1m"+str+"\033[0m"
}

func formatCommand(action string) string {
	return formatBoldColor(action,COLOR_BLUE)
}

func formatInfo(info string) string {
	return formatColor(info,COLOR_PURPLE)
}

func formatBoldColor(str string, color int) string {
	colorStr := strconv.Itoa(color)
	return "\033[1m\033["+colorStr+"m"+str+"\033[0m"
}

func formatColor(str string, color int) string {
	colorStr := strconv.Itoa(color)
	return "\033["+colorStr+"m"+str+"\033[0m"
}

//Main
func main() {
	fmt.Println("\033[0m"+formatBold("OLDUAR Client "+VERSION+"\n"))

	if(len(os.Args) >= 3) {
		//User + pass provided in command line arguments
		Username = os.Args[1]
		Password = os.Args[2]
	}

	Process("look","","")
	for {
		command, param1, param2 := "", "", ""
		fmt.Print("> \033[32m")
		fmt.Scanln(&command,&param1,&param2)
		fmt.Print("\033[0m")
		Process(strings.ToLower(command),param1,param2)
	}
}
