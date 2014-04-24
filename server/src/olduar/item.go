package olduar

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"math/rand"
)

// Loader for item templates

func LoadItems() bool {

	files, err := ioutil.ReadDir(MainServerConfig.DirItems);
	if(err != nil) {
		fmt.Println("Unable to load items from \""+MainServerConfig.DirItems+"\"")
		return false
	}

	//Load locations
	fmt.Println("Loading item files:")
	for _, file := range files {
		data, err := ioutil.ReadFile(MainServerConfig.DirItems+"/"+file.Name());
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

	return len(ItemTemplateDirectory) > 0
}

// Loot definition

func GetItemsFromLootTable(amount int, table ItemLootTable) ItemLootTable {
	//Loot table has no items or amount = 0? return none
	if(len(table) == 0 || amount == 0) {
		return nil
	}

	//Prepare loot bag
	loot := make(ItemLootTable,amount)

	//Loot table has only 1 item? Return that item @amount times
	if(len(table) == 1) {
		for i:=0; i<amount; i++ {
			loot[i] = table[0]
		}
		return loot
	}

	//Prepare selection of loot pointers
	minChance := 1.0
	for _, item := range table {
		if(item.Chance < minChance && item.Chance != 0.0) {
			minChance = item.Chance
		}
	}
	selectionAmount, s := 0, 0
	for _, item := range table {
		if(item.Chance != 0.0) {
			selectionAmount += (int)(item.Chance/minChance)
		}
	}
	selection, shuffledSelection := make(ItemLootTable,selectionAmount), make(ItemLootTable,selectionAmount)
	for _, item := range table {
		if(item.Chance != 0.0) {
			cnt := (int)(item.Chance / minChance)
			for i := 0; i < cnt; i++ {
				selection[s] = item
				s++
			}
		}
	}

	//Shuffle selection
	perm := rand.Perm(selectionAmount)
	for i, v := range perm {
		shuffledSelection[v] = selection[i]
	}

	//Pick items
	for i:=0;i<amount;i++ {
		loot[i] = shuffledSelection[rand.Intn(selectionAmount)]
	}

	//Return looted items
	return loot
}

func (template *ItemTemplate) GenerateItem() *Item {
	if(template == nil) {
		return nil
	}

	//Set action charges to unlimited for 0 value
	actions := template.Actions
	for index, action := range actions {
		if(action.Charges == 0) {
			actions[index].Charges = -1 //-1 = unlimited
		}
	}

	return &Item{
		Id: template.Id,
		Equipped: false,
		Attributes: template,
		Actions: actions,
	}
}

type ItemLootTable []*ItemLoot
type ItemLoot struct {
	Id string 					`json:"id"`
	Template *ItemTemplate 		`json:"-"`
	Chance float64 				`json:"chance"`
	MessageParty string			`json:"msg_party"`
	MessagePlayer string		`json:"msg_player"`
}

// Item Template definition

var ItemTemplateDirectory map[string]*ItemTemplate = make(map[string]*ItemTemplate)

type ItemTemplates []*ItemTemplate

type ItemTemplate struct {
	Id string 				`json:"id"`
	Quality int8 			`json:"quality"`
	Name string 			`json:"name"`
	Description string		`json:"desc"`
	Class string			`json:"class"`
	Type string				`json:"type"`
	Weight float64			`json:"weight"`
	Actions Actions			`json:"actions,omitempty"`
	Stats AttributeList		`json:"stats"`

	//Prepared response object
	Response ResponseItemDetail `json:"-"`
}

func (i *ItemTemplate) Prepare() {
	res := &i.Response

	res.Quality = i.Quality
	res.Name = &i.Name
	res.Description = &i.Description
	res.Class = &i.Class
	res.Type = &i.Type
	res.Weight = i.Weight
	res.Usable = len(i.Actions) > 0
	res.Stats = &i.Stats
}

type Item struct {
	Id string						`json:"id"`
	Attributes *ItemTemplate		`json:"-"`
	Actions Actions					`json:"actions,omitempty"`
	Equipped bool					`json:"equipped"`
}

func (item *Item) Load() bool {
	if(item.Attributes == nil) {
		attr, found := ItemTemplateDirectory[item.Id]
		if(found) {
			item.Attributes = attr
		}
		return found
	}
	return true
}

func (item *Item) GenerateResponse() ResponseItem {
	return ResponseItem{
		Quality: item.Attributes.Quality,
		Id: &item.Id,
		Name: &item.Attributes.Name,
		Description: &item.Attributes.Description,
		Equipped: item.Equipped,
		Usable: len(item.Actions) > 0,
	}
}

func (item *Item) Use(player *Player) {
	for index, _ := range item.Attributes.Actions {
		player.Room.DoAction(player,&item.Actions[index])
	}
}
