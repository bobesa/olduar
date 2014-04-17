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

type Message struct {
	Id int64					`json:"id"`
	Message string				`json:"text"`
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
			fmt.Println(formatCommand(item.Id)+": "+item.Name+" ("+item.Description+")")
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

type State struct {
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
	History []*Message	 		`json:"history"`
	Exits map[string]string		`json:"exits"`
	Actions map[string]string	`json:"actions"`
	Items Inventory				`json:"items"`
}

func (state *State) Print() {
	//History events
	if(len(state.History)>0) {
		for _, event := range state.History {
			fmt.Println(event.Message)
		}
		PrintLine()
	}

	//Location info
	fmt.Println(formatBold(state.Name))
	fmt.Println(formatInfo(state.Description))

	//Items
	count := len(state.Items)
	if(count > 0) {
		if(count == 1) {
			fmt.Print("Item on ground: ")
		} else {
			fmt.Print("Items on ground: ")
		}
		for index, item := range state.Items {
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
	count = len(state.Actions)
	if(count > 0) {
		if(count == 1) {
			fmt.Print("There is only 1 action: ")
		} else {
			fmt.Print("Possible actions are: ")
		}
		actionId := 0
		for action, desc := range state.Actions {
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
	count = len(state.Exits)
	if(count > 0) {
		if(count == 1) {
			fmt.Print("There is only 1 exit: ")
		} else {
			fmt.Print("Directions are: ")
		}
		exitId := 0
		for command, desc := range state.Exits {
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
		fmt.Println("There is no exit from this place")
	}
}

func RequireAuth() {
	fmt.Println(formatBoldColor("Wrong Username or Password",COLOR_RED))
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
	PrintLine()
}

func Fetch(command string, param string) []byte {
	//Check for param
	if(param != "") {
		command = command +"/"
	}

	//Do the request
	client := &http.Client{}
	req, err := http.NewRequest("GET", SERVER+"/"+GAME+"/"+command+param, nil)
	req.SetBasicAuth(Username, Password)
	resp, err := client.Do(req)
	if(err != nil) {
		fmt.Println(formatBoldColor("Network issue: Retrying in 5 seconds",COLOR_RED))
		time.Sleep(time.Second*5)
		return Fetch(command,param)
	}

	if(resp.StatusCode == 404) {
		RequireAuth()
		return Fetch(command,param)
	}

	//Read body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return nil
	}
	return body
}

func PrintLine() {
	fmt.Println("----------------------------------------")
}

func PrintHelp() {
	fmt.Println("\t"+formatCommand("help")+" \t\t\t\t List of available commands")
	fmt.Println("\t"+formatCommand("look")+" \t\t\t\t See current location description etc.")
	fmt.Println("\t"+formatCommand("go")+" [direction] \t\t Walk to a new place (example: "+formatInfo("go west")+")")
	fmt.Println("\t"+formatCommand("do")+" [task] \t\t\t Do action (example: "+formatInfo("do drink")+")")
	fmt.Println("\t"+formatCommand("inventory")+" \t\t\t See your inventory")
	fmt.Println("\t"+formatCommand("pickup")+" [item] \t\t Pickup object on ground (example: "+formatInfo("pickup shield")+")")
	fmt.Println("\t"+formatCommand("drop")+" [item] \t\t Drop object from inventory on ground (example: "+formatInfo("drop axe")+")")
	//fmt.Println("\tequip [item] \t\t Equip item from inventory (example: equip sword)")
	fmt.Println("\t"+formatCommand("inspect")+" [item] \t\t Displays complete info about item (example: "+formatInfo("inspect fishing_pole")+")")
	//fmt.Println("\tsay [message] \t\t Send message to party (example: say Hello guys)")
}

func Process(command string, param string) {
	PrintLine()
	switch(command){
	case "inventory":
		inventory := Inventory{}
		json.Unmarshal(Fetch(command,param),&inventory)
		inventory.Print()

	case "inspect":
		var item *ItemDetail = nil
		json.Unmarshal(Fetch(command,param),&item)
		item.Print()

	case "pickup", "drop", "use":
		if(param != "") {
			var state *State = nil
			json.Unmarshal(Fetch(command,param),&state)
			if(state != nil) {
				state.Print()
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

	case "do", "go", "look":
		var state *State = nil
		json.Unmarshal(Fetch(command,param),&state)
		if(state != nil) {
			state.Print()
		} else {
			fmt.Println(formatBold("Something went wrong"))
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

	Process("look","")
	for {
		command, param := "", ""
		fmt.Print("> \033[32m")
		fmt.Scanln(&command,&param)
		fmt.Print("\033[0m")
		Process(strings.ToLower(command),strings.ToLower(param))
	}
}
