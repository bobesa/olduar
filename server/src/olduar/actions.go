package olduar

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Actions []Action
type Action struct {
	Id           string                 `json:"id,omitempty"`
	Description  string                 `json:"desc,omitempty"`
	Action       string                 `json:"action"`
	Charges      int                    `json:"charges,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	Requirements ActionRequirements     `json:"requirements,omitempty"`

	worker Actioner
}

func (a *Action) Prepare() bool {
	if a.worker == nil {
		data, err := json.Marshal(a.Config)
		if err == nil {
			switch a.Action {
			case "message":
				a.worker = new(ActionTypeMessage)
			case "location":
				a.worker = new(ActionTypeLocation)
			case "effect":
				a.worker = new(ActionTypeEffect)
			case "give":
				a.worker = new(ActionTypeGive)
			default:
				return false
			}
			err = json.Unmarshal(data, &a.worker)
			if err != nil || a.worker.Prepare() {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (action *Action) Do(room *Room, player *Player) {
	if action.worker == nil {
		action.Prepare()
	}
	action.worker.Do(player, room, action)
}

//Action requirements

type ActionRequirements []*ActionRequirement
type ActionRequirement struct {
	Type         string `json:"type"`
	Value        string `json:"value"`
	ErrorMessage string `json:"errorMsg"`
}

//Actioner interface

type Actioner interface {
	Prepare() bool
	Do(*Player, *Room, *Action)
}

func AppendVariablesToString(str string, player *Player, config map[string]interface{}) string {
	str = strings.Replace(str, "%player%", player.Name, -1)
	for key, value := range config {
		str = strings.Replace(str, "%"+key+"%", fmt.Sprint(value), -1)
	}
	return str
}

//Message action type

type ActionTypeMessage struct {
	MessageAll    string `json:"msgAll"`
	MessageParty  string `json:"msgParty"`
	MessagePlayer string `json:"msgPlayer"`
}

func (a *ActionTypeMessage) Prepare() bool {
	return a.MessageAll != "" || a.MessageParty != "" || a.MessagePlayer != ""
}

func (a *ActionTypeMessage) Do(player *Player, room *Room, action *Action) {
	if a.MessageAll != "" {
		room.TellAll(AppendVariablesToString(a.MessageAll, player, action.Config))
	}
	if a.MessageParty != "" {
		room.TellAllExcept(AppendVariablesToString(a.MessageParty, player, action.Config), player)
	}
	if a.MessagePlayer != "" {
		room.Tell(AppendVariablesToString(a.MessagePlayer, player, action.Config), player)
	}
}

//Location action type

type ActionTypeLocation struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (a *ActionTypeLocation) Prepare() bool {
	return a.Type != "" && a.Value != ""
}

func (a *ActionTypeLocation) Do(player *Player, room *Room, action *Action) {
	switch a.Type {
	case "use":
		room.CurrentLocation.DoAction(room, player, a.Value)
	}
}

//Effect action type

type ActionTypeEffect struct {
	Type          string  `json:"type"`
	Value         float64 `json:"value"`
	MessageAll    string  `json:"msgAll"`
	MessageParty  string  `json:"msgParty"`
	MessagePlayer string  `json:"msgPlayer"`
}

func (a *ActionTypeEffect) Prepare() bool {
	return a.Type != ""
}

func (a *ActionTypeEffect) Do(player *Player, room *Room, action *Action) {
	//Do effect
	switch a.Type {
	case "damage":
		player.Damage(a.Value, nil, nil)
	case "heal":
		player.Heal(a.Value, nil)
	}
	//Send messages
	if a.MessageAll != "" {
		room.TellAll(AppendVariablesToString(a.MessageAll, player, action.Config))
	}
	if a.MessageParty != "" {
		room.TellAllExcept(AppendVariablesToString(a.MessageParty, player, action.Config), player)
	}
	if a.MessagePlayer != "" {
		room.Tell(AppendVariablesToString(a.MessagePlayer, player, action.Config), player)
	}
}

//Give action type

type ActionTypeGive struct {
	Amount int           `json:"amount"`
	Items  ItemLootTable `json:"items"`
}

func (a *ActionTypeGive) Prepare() bool {
	//No items in table = fail
	if len(a.Items) == 0 {
		return false
	}

	//Cycle trough items and bind template
	for _, item := range a.Items {
		item.Template = ItemTemplateDirectory[item.Id]
	}

	return true
}

func (a *ActionTypeGive) Do(player *Player, room *Room, action *Action) {
	if a.Amount > 0 {
		//Give looted items
		items := GetItemsFromLootTable(a.Amount, a.Items)
		for _, item := range items {
			if item.MessagePlayer != "" {
				room.Tell(AppendVariablesToString(item.MessagePlayer, player, action.Config), player)
			}
			if item.MessageParty != "" {
				room.TellAllExcept(AppendVariablesToString(item.MessageParty, player, action.Config), player)
			}
			if item.Template != nil {
				player.Inventory = append(player.Inventory, item.Template.GenerateItem())
			}
		}
	}
}
