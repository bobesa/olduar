package olduar

type Inventory []*Item
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

func (p *Player) Give(entry string) {
	template, found := ItemTemplateDirectory[entry]
	if(found) {
		p.Inventory = append(p.Inventory,ItemFromTemplate(template))
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
