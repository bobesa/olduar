package olduar

import (
	"fmt"
	"math/rand"
)

type Actions []Action
type Action struct {
	Id string							`json:"id"`
	Description string					`json:"desc,omitempty"`
	Action string						`json:"action"`
	Charges int							`json:"charges"`
	Config map[string]interface{}		`json:"config,omitempty"`
	Requirements ActionRequirements		`json:"requirements"`
}

type ActionRequirements []*ActionRequirement
type ActionRequirement struct {
	Type string 				`json:"type"`
	Value string 				`json:"value"`
	ErrorMessage string 		`json:"error_msg"`
}

func (action *Action) Do(state *GameState, player *Player) {
	entry, found := ActionsDirectory[action.Action]
	if(found) {
		entry(state,player,action.Config)
	}
}

type LocationExits []LocationExit
type LocationExit struct{
	Id string						`json:"id"`
	Region string					`json:"region,omitempty"`
	Entry string					`json:"entry,omitempty"`
	Target *Location				`json:"target"`
}

type Location struct {
	Name string 					`json:"name"`
	Region string					`json:"region,omitempty"`
	Description string 				`json:"desc"`
	DescriptionShort string 		`json:"desc_short"`
	Exits LocationExits 			`json:"exits"`
	Actions Actions					`json:"actions,omitempty"`
	Items Inventory					`json:"items,omitempty"`
	Visited bool					`json:"visited"`
}

func (loc *Location) Describe() {
	fmt.Println("----------------------------")
	fmt.Println(loc.Name)
	fmt.Println(loc.Description+"\n")

	//Actions
	actionCount := len(loc.Actions)
	if(actionCount > 0) {
		if(actionCount == 1) {
			fmt.Print("There is only 1 action: ")
		} else {
			fmt.Print("Possible actions are: ")
		}
		for index, action := range loc.Actions {
			fmt.Print(action.Id)
			if(action.Description != "") {
				fmt.Print(" ("+action.Description+")")
			}
			if(actionCount-2 == index) {
				fmt.Print(" & ")
			} else if(actionCount-1 != index) {
				fmt.Print(", ")
			}
		}
		fmt.Println();
	}

	//Exits
	exitCount := len(loc.Exits)
	if(exitCount > 0) {
		if(exitCount == 1) {
			fmt.Print("There is only 1 exit: ")
		} else {
			fmt.Print("Directions are: ")
		}
		for index, exit := range loc.Exits {
			fmt.Print(exit.Id)
			if(exit.Target != nil) {
				fmt.Print(" ("+exit.Target.DescriptionShort+")")
			}
			if(exitCount-2 == index) {
				fmt.Print(" & ")
			} else if(exitCount-1 != index) {
				fmt.Print(", ")
			}
		}
		fmt.Println();
	} else {
		fmt.Println("There is no exit from this place")
	}

	fmt.Println("----------------------------")
}

func (loc *Location) Visit() bool {

	if(len(loc.Exits)==0) {
		//No exits? Generate exits
		exitNames := []string{
			"west",
			"east",
			"south",
			"north",
			"northeast",
			"northwest",
			"southwest",
			"southeast",
		}
		for i:=0; i<2; i++ {
			index := rand.Intn(len(exitNames))
			exitName := exitNames[index]
			newExitNames := []string{}
			for i, exit := range exitNames {
				if(i != index) {
					newExitNames = append(newExitNames,exit)
				}
			}
			exitNames = newExitNames

			loc.Exits = append(loc.Exits,LocationExit{
					Id:exitName,
					Target:CreateLocationFromRegion(loc.Region),
				})
		}
	} else {
		//Cycle trough exits and generate Targets if not shown
		for index, exit := range loc.Exits {
			if(exit.Target == nil && exit.Region != "") {
				loc.Exits[index].Target = CreateLocationFromRegion(exit.Region)
				loc.Exits[index].Region = ""
			} else if(exit.Target == nil && exit.Entry != "") {
				loc.Exits[index].Target = CreateLocationFromEntry(exit.Entry)
				loc.Exits[index].Entry = ""
			}
		}

	}

	if(!loc.Visited) {
		loc.Visited = true
		return true
	}
	return false
}
