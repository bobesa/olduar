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
	Actions Actions			`json:"actions,omitempty"`

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
	Actions Actions					`json:"actions,omitempty"`
	Equipped bool					`json:"equipped"`
}

func (item *Item) GenerateResponse() ResponseItem {
	return ResponseItem{
		Id: &item.Id,
		Name: &item.Attributes.Name,
		Description: &item.Attributes.Description,
	}
}

func (item *Item) Use(player *Player) {
	for index, _ := range item.Attributes.Actions {
		player.GameState.DoAction(player,&item.Attributes.Actions[index])
	}
}
