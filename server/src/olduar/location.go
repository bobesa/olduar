package olduar

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
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
	DescriptionShort string 		`json:"descShort"`
	Exits LocationExits 			`json:"exits"`
	Actions map[string]*Action		`json:"actions,omitempty"`
	Npcs []*Npc						`json:"npcs,omitempty"`
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

// Loader for location templates

func LoadLocations() bool {

	files := GetFilesFromDirectory(MainServerConfig.DirLocations);
	if(len(files) == 0) {
		fmt.Println("Unable to load locations from \""+MainServerConfig.DirLocations+"\"")
		return false
	}

	type LoadingRegion struct {
		Region string 					`json:"region"`
		Description string 				`json:"desc"`
		Locations LocationTemplates 	`json:"locations"`
	}

	//Load locations
	fmt.Println("Loading location files:")
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename);
		if(err == nil) {
			region := LoadingRegion{}
			err := json.Unmarshal(data,&region)
			if(err == nil && region.Region != "") {
				fmt.Println("\t" + filename + ": (Region: "+region.Region+") " + region.Description);
				_, found := LocationTemplateDirectoryRegions[region.Region];
				if(!found) {
					LocationTemplateDirectoryRegions[region.Region] = make(LocationTemplates,0)
				}
				for _, location := range region.Locations {
					location.Region = region.Region
					//Set action charges to unlimited for 0 value
					for index, action := range location.Actions {
						if(action.Charges == 0) {
							location.Actions[index].Charges = -1 //-1 = unlimited
						}
					}
					//Add location as entry
					if(location.Id != "") {
						LocationTemplateDirectoryEntries[location.Id] = location
					}
					//Add location to region
					LocationTemplateDirectoryRegions[region.Region] = append(LocationTemplateDirectoryRegions[region.Region],location)
				}
			} else {
				fmt.Println("\t" + filename + ": Failed to load")
			}
		} else {
			fmt.Println("\t" + filename + ": Failed to load")
		}
	}

	//Check for amount of regions (must be > 2) & "start" region
	_, found := LocationTemplateDirectoryRegions["start"];
	if(!found) {
		fmt.Println("Error: \"start\" region not found!")
		return false
	}

	return true

}

func CreateLocationFromTemplate(template *LocationTemplate) *Location {
	loc := Location{}

	loc.Name = template.Name
	loc.Description = template.Description
	loc.DescriptionShort = template.DescriptionShort
	loc.Exits = template.Exits
	loc.Region = template.Region
	loc.Visited = false

	//Cycle trough actions and create map
	loc.Actions = make(map[string]*Action)
	for _, action := range template.Actions {
		loc.Actions[action.Id] = &Action{
			Id: action.Id,
			Description: action.Description,
			Action: action.Action,
			Charges: action.Charges,
			Config: action.Config,
			Requirements: action.Requirements,
		}
	}

	//Generate items on ground
	loc.Items = make(Inventory,0)
	for _, item := range template.Items {
		if(item.Chance > 0 && item.Chance < rand.Float64()) {
			continue
		}
		if(item.Id != "") {
			entry := ItemTemplateDirectory[item.Id]
			if(entry != nil) {
				loc.Items = append(loc.Items,entry.GenerateItem())
			}
		} else if(item.Group != "") {
			//Give any template from item group
		}
	}

	//Generate characters
	loc.Npcs = make([]*Npc,0)
	for _, npc := range template.Npcs {
		if(npc.Chance > 0 && npc.Chance < rand.Float64()) {
			continue
		}
		if(npc.Id != "") {
			entry, found := CharacterTemplateDirectory[npc.Id]
			if(found) {
				loc.Npcs = append(loc.Npcs,entry.MakeInstance())
			}
		} else if(npc.Group != "") {
			//Give any template from item group
		}
	}

	return &loc
}

func CreateLocationFromRegion(region string) *Location {
	templateBank, found := LocationTemplateDirectoryRegions[region]
	if(!found) {
		return nil
	}
	template := templateBank[rand.Intn(len(templateBank))]
	return CreateLocationFromTemplate(template)
}

func CreateLocationFromEntry(entry string) *Location {
	template, found := LocationTemplateDirectoryEntries[entry]
	if(!found) {
		return nil
	}
	return CreateLocationFromTemplate(template)
}

// Locations types and functions

var LocationTemplateDirectoryEntries = make(map[string]*LocationTemplate)
var LocationTemplateDirectoryRegions = make(map[string]LocationTemplates)

type LocationEntryTemplates []*LocationEntryTemplate
type LocationEntryTemplate struct {
	Id string 		`json:"id"`
	Group string 	`json:"group"`
	Chance float64	`json:"chance"`
}

type LocationTemplates []*LocationTemplate
type LocationTemplate struct {
	Id string 						`json:"id"`
	Name string 					`json:"name"`
	Region string					`json:"region,omitempty"`
	Description string 				`json:"desc"`
	DescriptionShort string 		`json:"descShort"`
	Actions Actions					`json:"actions,omitempty"`
	Exits LocationExits				`json:"exits,omitempty"`
	Items LocationEntryTemplates 	`json:"items"`
	Npcs LocationEntryTemplates 	`json:"npcs"`
}
