package olduar

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

// Loader for location templates

func LoadLocations(path string) bool {

	files, err := ioutil.ReadDir(path);
	if(err != nil) {
		fmt.Println("Unable to load locations from \""+path+"\"")
		return false
	}

	type LoadingRegion struct {
		Region string 					`json:"region"`
		Description string 				`json:"desc"`
		Locations LocationTemplates 	`json:"locations"`
	}

	//Load locations
	fmt.Println("Loading location files:")
	for _, file := range files {
		data, err := ioutil.ReadFile(path+"/"+file.Name());
		if(err == nil) {
			region := LoadingRegion{}
			err := json.Unmarshal(data,&region)
			if(err == nil && region.Region != "") {
				fmt.Println("\t" + file.Name() + ": (Region: "+region.Region+") " + region.Description);
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
				fmt.Println("\t" + file.Name() + ": Failed to load")
			}
		} else {
			fmt.Println("\t" + file.Name() + ": Failed to load")
		}
	}

	//Check for amount of regions (must be > 2) & "start" region
	_, found := LocationTemplateDirectoryRegions["start"];
	if(!found) {
		fmt.Println("Error: \"start\" region not found!")
		return false
	}

	fmt.Println()

	return true

}

func CreateLocationFromTemplate(template *LocationTemplate) *Location {
	loc := Location{}

	loc.Name = template.Name
	loc.Description = template.Description
	loc.DescriptionShort = template.DescriptionShort
	loc.Actions = template.Actions
	loc.Exits = template.Exits
	loc.Region = template.Region
	loc.Visited = false

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

type LocationTemplates []*LocationTemplate
type LocationTemplate struct {
	Id string 					`json:"id"`
	Name string 				`json:"name"`
	Region string				`json:"region,omitempty"`
	Description string 			`json:"desc"`
	DescriptionShort string 	`json:"desc_short"`
	Actions Actions				`json:"actions,omitempty"`
	Exits LocationExits			`json:"exits,omitempty"`
}

