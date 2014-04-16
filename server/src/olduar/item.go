package olduar

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"math/rand"
)

// Loader for item templates

func LoadItems(path string) bool {

	files, err := ioutil.ReadDir(path);
	if(err != nil) {
		fmt.Println("Unable to load items from \""+path+"\"")
		return false
	}

	//Load locations
	fmt.Println("Loading item files:")
	for _, file := range files {
		data, err := ioutil.ReadFile(path+"/"+file.Name());
		if(err == nil) {
			items := make(ItemTemplates,0)
			err := json.Unmarshal(data,&items)
			if(err == nil) {
				fmt.Println("\t" + file.Name() + ": loaded "+strconv.Itoa(len(items))+" items")
				for _, item := range items {
					item.Prepare()
					ItemTemplateDirectory[item.Id] = item
				}
			} else {
				fmt.Println("\t" + file.Name() + ": Failed to load")
			}
		} else {
			fmt.Println("\t" + file.Name() + ": Failed to load")
		}
	}

	fmt.Println()

	return len(ItemTemplateDirectory) > 0
}

// Loot definition

func GetItemsFromLootTable(player *Player, amount int, table ItemLootTable) []*Item {
	//Loot table has no items or amount = 0? return none
	if(len(table) == 0 || amount == 0) {
		return nil
	}

	//Prepare loot bag
	loot := make([]*Item,amount)

	//Loot table has only 1 item? Return that item @amount times
	if(len(table) == 1) {
		for i:=0; i<amount; i++ {
			loot[i] = ItemFromTemplate(table[0].Template)
		}
		return loot
	}

	//Prepare selection of loot pointers
	minChance := 1.0
	for _, item := range table {
		if(item.Chance < minChance) {
			minChance = item.Chance
		}
	}
	selectionAmount, s := 0, 0
	for _, item := range table {
		selectionAmount += (int)(item.Chance/minChance)
	}
	selection, shuffledSelection := make(ItemLootTable,selectionAmount), make(ItemLootTable,selectionAmount)
	for _, item := range table {
		cnt := (int)(item.Chance/minChance)
		for i:=0;i<cnt;i++ {
			selection[s] = item
			s++
		}
	}

	//Shuffle selection
	perm := rand.Perm(selectionAmount)
	for i, v := range perm {
		shuffledSelection[v] = selection[i]
	}

	//Pick items
	for i:=0;i<amount;i++ {
		loot[i] = ItemFromTemplate(shuffledSelection[rand.Intn(selectionAmount)].Template)
	}

	//Return looted items
	return loot
}

func ItemFromTemplate(template *ItemTemplate) *Item {
	return &Item{
		Id: template.Id,
		Equipped: false,
		Attributes: template,
	}
}

type ItemLootTable []*ItemLoot
type ItemLoot struct {
	Template *ItemTemplate
	Chance float64
	MessageParty, MessagePlayer string
}

// Item Template definition

var ItemTemplateDirectory map[string]*ItemTemplate = make(map[string]*ItemTemplate)

type ItemTemplates []*ItemTemplate

type ItemTemplate struct {
	Id string 				`json:"id"`
	Name string 			`json:"name"`
	Description string		`json:"desc"`
	Type string				`json:"type"`
	Weight float64			`json:"weight"`

	//Stats
	DamageMin int64			`json:"damage_min"`
	DamageMax int64			`json:"damage_max"`

	//Prepared response object
	Response ResponseItemDetail `json:"-"`
}

func (i *ItemTemplate) Prepare() {
	i.Response.Name = &i.Name
	i.Response.Description = &i.Description
}

type Item struct {
	Id string						`json:"id"`
	Attributes *ItemTemplate		`json:"-"`
	Equipped bool					`json:"equipped"`
}
