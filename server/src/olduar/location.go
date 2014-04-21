package olduar

import (
	"math/rand"
)

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
	Actions map[string]*Action		`json:"actions,omitempty"`
	Items Inventory					`json:"items,omitempty"`
	Visited bool					`json:"visited"`
	Current bool					`json:"current"`
	Parent *Location				`json:"-"`
}

func (loc *Location) DoAction(room *Room, player *Player, actionName string) {
	action, found := loc.Actions[actionName]
	if(found) {
		room.DoAction(player,action)
	}
}

func (loc *Location) Visit() bool {
	loc.Current = true;

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
