package olduar

import (
	"strings"
	"fmt"
)

type ActionFunction func (state *GameState, player *Player, config map[string]interface{})

func AppendVariablesToString(str string, player *Player, config map[string]interface{}) string {
	str = strings.Replace(str,"%player%",player.Name,-1)
	for key, value := range config {
		str = strings.Replace(str,"%"+key+"%",fmt.Sprint(value),-1)
	}
	return str
}

var ActionsDirectory = make(map[string]ActionFunction)

func InitializeActions() {

	ActionsDirectory["message"] = func(state *GameState,player *Player,config map[string]interface{}) {} //Automatically processed - just a placeholder

	ActionsDirectory["location"] = func(state *GameState,player *Player,config map[string]interface{}) {
		actionType, found := config["type"]
		if(found) {
			location := state.CurrentLocation
			switch (actionType){
			case "use":
				actionName, found := config["value"]
				if(found) {
					location.DoAction(state,player,fmt.Sprint(actionName))
				}
			}
		}
	}

	ActionsDirectory["give"] = func(state *GameState,player *Player,config map[string]interface{}) {
		//Amount of looted items
		amount := 1
		value, found := config["amount"]
		if (found) {
			amount = (int)(value.(float64))
		}

		//Prepare loot table
		table := ItemLootTable{}
		for _, itemConfig := range config["items"].([]interface{}) {
			config := itemConfig.(map[string]interface{})
			item := &ItemLoot{}
			value, found := config["id"]
			if (found) {
				item.Template = ItemTemplateDirectory[value.(string)]

				value, found = config["chance"]
				if (found) {
					item.Chance = value.(float64)
				} else {
					item.Chance = 1.0
				}

				value, found = config["msg_party"]
				if (found) {
					item.MessageParty = AppendVariablesToString(value.(string),player,config)
				}

				value, found = config["msg_player"]
				if (found) {
					item.MessagePlayer = AppendVariablesToString(value.(string),player,config)
				}

				table = append(table, item)
			}
		}

		//Get looted items
		items := GetItemsFromLootTable(amount, table)
		for _, item := range items {
			if (item.MessagePlayer != "") {
				state.Tell(item.MessagePlayer, player)
			}
			if (item.MessageParty != "") {
				state.TellAllExcept(item.MessageParty, player)
			}
			if (item.Template != nil) {
				player.Inventory = append(player.Inventory, item.Template.GenerateItem())
			}
		}
	}

	ActionsDirectory["effect"] = func(state *GameState,player *Player,config map[string]interface{}) {
		//Process effect
		fxType, found := config["type"]
		if(found) {
			switch(fxType){
			case "hurt":
				value, found := config["value"]
				if(found) {
					player.Damage((int64)(value.(float64)))
				}
				break
			case "heal":
				value, found := config["value"]
				if(found) {
					player.Heal((int64)(value.(float64)))
				}
				break
			}
		}
	}

}

