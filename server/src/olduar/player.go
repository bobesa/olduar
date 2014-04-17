package olduar

import (
	"encoding/base64"
	"time"
)

var ActivePlayers map[string]*Player = make(map[string]*Player)

type Players []*Player

type Player struct {
	Username string			`json:"username"`
	Password string			`json:"password"`
	AuthToken string		`json:"-"`

	Name string 			`json:"name"`
	Health int64 			`json:"health"`
	MaxHealth int64 		`json:"health_max"`
	Inventory Inventory		`json:"inventory"`

	GameState *GameState 	`json:"-"`
	LastResponseId int64	`json:"-"`
	LastResponse time.Time	`json:"-"`
}

func (p *Player) Activate() {
	p.AuthToken = "Basic "+base64.StdEncoding.EncodeToString([]byte(p.Username+":"+p.Password))
	ActivePlayers[p.AuthToken] = p
}

func (p *Player) Deactivate() {

}

func (p *Player) Ability(target *Npc, skill string) {
	//TODO: Add ability functionality
}

func (p *Player) Attack(target *Npc) {
	//TODO: Add attack functionality
}

func (p *Player) Pickup(entry string) bool {
	item := p.GameState.CurrentLocation.Items.Get(entry)
	if(item != nil) {
		var weight float64 = item.Attributes.Weight
		for _, invItem := range p.Inventory {
			weight += invItem.Attributes.Weight
		}
		if(weight <= 1.0) {
			p.GameState.CurrentLocation.Items.Remove(item)
			p.Inventory.Add(item)
			p.GameState.Tell("You picked up "+item.Attributes.Name+" from ground",p)
			p.GameState.TellAllExcept(p.Name+" picked up "+item.Attributes.Name+" from ground",p)
		} else {
			p.GameState.Tell("You cannot keep more items in inventory",p)
		}
		return true
	}
	return false
}

func (p *Player) Use(entry string) bool {
	item := p.Inventory.Get(entry)
	if(item != nil) {
		item.Use(p)
		if(item.Attributes.Type == "consumable") {
			p.Inventory.Remove(item)
		}
		return true
	}
	return false
}

func (p *Player) Drop(entry string) bool {
	item := p.Inventory.Get(entry)
	if(item != nil) {
		p.GameState.CurrentLocation.Items.Add(item)
		p.Inventory.Remove(item)
		return true
	}
	return false
}

func (p *Player) Owns(entry string) bool {
	return p.Inventory.Get(entry) != nil
}

func (p *Player) Give(entry string) {
	template, found := ItemTemplateDirectory[entry]
	if(found) {
		p.Inventory.Add(template.GenerateItem())
	}
}

func (p *Player) Heal(value int64) {
	p.Health += value
	if(p.Health > p.MaxHealth) {
		p.Health = p.MaxHealth
	}
}

func (p *Player) Damage(value int64) {
	p.Health -= value
	if(p.Health <= 0) {
		p.Health = 0
	}
}
