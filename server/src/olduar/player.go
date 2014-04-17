package olduar

type Players []*Player

type Player struct {
	Username string			`json:"username"`
	HashPass string			`json:"hash"`

	Name string 			`json:"name"`
	Health int64 			`json:"health"`
	MaxHealth int64 		`json:"health_max"`
	Inventory Inventory		`json:"inventory"`

	GameState *GameState 	`json:"-"`
	LastResponseId int64	`json:"-"`
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
		p.GameState.CurrentLocation.Items.Remove(item)
		p.Inventory.Add(item)
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
