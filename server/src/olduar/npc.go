package olduar

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

var CharacterTemplateDirectory map[string]Npc

func LoadCharacters() bool {
	CharacterTemplateDirectory = make(map[string]Npc)

	files := GetFilesFromDirectory(MainServerConfig.DirCharacters);
	if(len(files) == 0) {
		fmt.Println("Unable to load locations from \""+MainServerConfig.DirCharacters+"\"")
		return false
	}
	//Load locations
	fmt.Println("Loading character files:")
	for _, filename := range files {
		data, err := ioutil.ReadFile(filename);
		if (err == nil) {
			npcs := make([]Npc,0)
			err := json.Unmarshal(data,&npcs)
			if(err == nil) {
				fmt.Println("\t" + filename + ": loaded "+strconv.Itoa(len(npcs))+" characters")
				for _, npc := range npcs {
					npc.Health = npc.MaxHealth
					CharacterTemplateDirectory[npc.Id] = npc
				}
			}
		}
	}

	return len(CharacterTemplateDirectory) > 0
}

type Npc struct {
	Id string 				`json:"id"`
	Guid GUID				`json:"guid"`
	Name string 			`json:"name"`
	Description string 		`json:"desc"`
	Stats AttributeList		`json:"stats"`
	Health float64 			`json:"health"`
	MaxHealth float64 		`json:"health_max"`
	Money int64	 			`json:"money"`
	Friendly bool 			`json:"friendly"`
	//TODO: Implement loot drop
}

func (npc Npc) MakeInstance() *Npc {
	fmt.Println("Making instance of \""+npc.Name+"\"")
	npcCopy := npc
	npcCopy.Guid = GenerateGUID()
	return &npcCopy
}

func (npc *Npc) GenerateResponse() ResponseNpc {
	return ResponseNpc{
		Id: &npc.Id,
		Name: &npc.Name,
		Description: &npc.Description,
		Health: npc.Health,
		HealthMax: npc.MaxHealth,
		Friendly: npc.Friendly,
	}
}

func (npc *Npc) GetStats() AttributeList {
	return npc.Stats
}

func (npc *Npc) Heal(value float64, log *CombatQueue) {
	npc.Health += value
	if(npc.Health > npc.MaxHealth) {
		npc.Health = npc.MaxHealth
	}
}

func (npc *Npc) Damage(value float64, log *CombatQueue, attacker Fighter) {
	npc.Health -= value
	if(npc.Health <= 0) {
		npc.Health = 0
		npc.Die(log,attacker)
	}
}

func (npc *Npc) Die(log *CombatQueue, attacker Fighter) {
	if(log != nil) {
		log.Log(attacker,"You killed "+npc.GetName(),attacker.GetName() + " killed " + npc.GetName())
	}
}

func (npc *Npc) GetId() string {
	return npc.Id
}

func (npc *Npc) GetGUID() GUID {
	return npc.Guid
}

func (npc *Npc) IsPlayer() bool {
	return false
}

func (npc *Npc) IsAlive() bool {
	return npc.Health > 0.0
}

func (npc *Npc) GetName() string {
	return npc.Name
}

func (npc *Npc) GetTeam() CombatTeam {
	//TODO: implement teams (currently player + friendly npc is 0, enemy npc is 1)
	if(npc.Friendly) {
		return 0
	}
	return 1
}
