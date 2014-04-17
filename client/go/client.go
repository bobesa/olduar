package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
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
		fmt.Println(formatColor(item.Description,COLOR_PURPLE))
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
		fmt.Println("----------------------------------------")
	}

	//Location info
	fmt.Println(formatBold(state.Name))
	fmt.Println(formatColor(state.Description,COLOR_PURPLE))

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

func Fetch(command string, param string) []byte {
	if(param != "") {
		command = command +"/"
	}
	resp, err := http.Get(SERVER+"/"+GAME+"/"+command+param)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if(err != nil) {
		return nil
	}
	return body
}

func PrintHelp() {
	fmt.Println("\thelp \t\t\t\t List of available commands")
	fmt.Println("\tlook \t\t\t\t See current location description etc.")
	fmt.Println("\tgo [direction] \t\t Walk to a new place (example: go west)")
	fmt.Println("\tdo [task] \t\t\t Do action (example: do drink)")
	fmt.Println("\tinventory \t\t\t See your inventory")
	fmt.Println("\tpickup [item] \t\t Pickup object on ground (example: pickup shield)")
	fmt.Println("\tdrop [item] \t\t Drop object from inventory on ground (example: drop axe)")
	//fmt.Println("\tequip [item] \t\t Equip item from inventory (example: equip sword)")
	fmt.Println("\tinspect [item] \t\t Displays complete info about item (example: inspect fishing_pole)")
	//fmt.Println("\tsay [message] \t\t Send message to party (example: say Hello guys)")
}

func Process(command string, param string) {
	fmt.Println("----------------------------------------")
	switch(command){
	case "inventory":
		inventory := Inventory{}
		json.Unmarshal(Fetch(command,param),&inventory)
		inventory.Print()

	case "inspect":
		var item *ItemDetail = nil
		json.Unmarshal(Fetch(command,param),&item)
		item.Print()

	case "pickup", "drop":
		if(param != "") {
			var state *State = nil
			json.Unmarshal(Fetch(command,param),&state)
			if(state != nil) {
				state.Print()
			} else {
				fmt.Print(formatBold("You cannot " + command + " ") + formatCommand(param))
				if(command == "drop") {
					fmt.Println(formatBold(" because it's not in your inventory"))
				} else {
					fmt.Println(formatBold(" because it's not here"))
				}
			}
		} else {
			fmt.Println(formatBold("You need to say what you want to " + command))
		}

	case "do", "go", "look":
		var state *State = nil
		json.Unmarshal(Fetch(command,param),&state)
		state.Print()

	case "help":
		fmt.Println("Available commands")
		PrintHelp()
	default:
		fmt.Println("Unknown command "+command+", see available commands below")
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

	Process("look","")
	for {
		command, param := "", ""
		fmt.Print("> \033[32m")
		fmt.Scanln(&command,&param)
		fmt.Print("\033[0m")
		Process(strings.ToLower(command),strings.ToLower(param))
	}
}
