package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const (
	VERSION = "0.01a"
	SERVER = "http://localhost:8080"
	GAME = "test"
)

// Types

type Message struct {
	Id int64					`json:"id"`
	Message string				`json:"text"`
}

type State struct {
	Name string 				`json:"name"`
	Description string 			`json:"desc"`
	History []*Message	 		`json:"history"`
	Exits map[string]string		`json:"exits"`
	Actions map[string]string	`json:"actions"`
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
	fmt.Println("\tlook \t\t\t\t See current location description etc.")
	//fmt.Println("\tinventory \t\t\t See your inventory")
	fmt.Println("\tgo [direction] \t\t Walk to a new place (example: go west)")
	fmt.Println("\tdo [task] \t\t\t Do action (example: do drink)")
	//fmt.Println("\tpickup [item] \t\t Pickup object on ground (example: pickup shield)")
	//fmt.Println("\tdrop [item] \t\t Drop object from inventory on ground (example: drop axe)")
	//fmt.Println("\tequip [item] \t\t Equip item from inventory (example: equip sword)")
	//fmt.Println("\tinspect [item] \t\t Displays complete info about item (example: inspect fishing_pole)")
	//fmt.Println("\tsay [message] \t\t Send message to party (example: say Hello guys)")
}

func Process(command string, param string) {
	switch(command){
	case "do", "go", "look":
		state := &State{}
		json.Unmarshal(Fetch(command,param),&state)

		//History events
		if(len(state.History)>0) {
			fmt.Println("----------------------------------------")
			for _, event := range state.History {
				fmt.Println(event.Message)
			}
		}

		//Location info
		fmt.Println("----------------------------------------")
		fmt.Println(state.Name)
		fmt.Println(state.Description)

		//Actions
		actionCount := len(state.Actions)
		if(actionCount > 0) {
			if(actionCount == 1) {
				fmt.Print("There is only 1 action: ")
			} else {
				fmt.Print("Possible actions are: ")
			}
			actionId := 0
			for action, desc := range state.Actions {
				fmt.Print(action+" ("+desc+")")
				if(actionCount-2 == actionId) {
					fmt.Print(" & ")
				} else if(actionCount-1 != actionId) {
					fmt.Print(", ")
				}
				actionId++
			}
			fmt.Println();
		}

		//Exits
		exitCount := len(state.Exits)
		if(exitCount > 0) {
			if(exitCount == 1) {
				fmt.Print("There is only 1 exit: ")
			} else {
				fmt.Print("Directions are: ")
			}
			exitId := 0
			for command, desc := range state.Exits {
				fmt.Print(command+" ("+desc+")")
				if(exitCount-2 == exitId) {
					fmt.Print(" & ")
				} else if(exitCount-1 != exitId) {
					fmt.Print(", ")
				}
				exitId++
			}
			fmt.Println();
		} else {
			fmt.Println("There is no exit from this place")
		}
	case "help":
		fmt.Println("Available commands")
		PrintHelp()
	default:
		fmt.Println("Unknown command "+command+", see available commands below")
		PrintHelp()
	}
}

func main() {
	fmt.Println("OLDUAR Client "+VERSION+"\n")

	Process("look","")
	for {
		command, param := "", ""
		fmt.Print("> ")
		fmt.Scanln(&command,&param)
		Process(command,param)
	}
}
