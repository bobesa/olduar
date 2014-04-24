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

	files, err := ioutil.ReadDir(MainServerConfig.DirCharacters);
	if(err != nil) {
		fmt.Println("Unable to load locations from \""+MainServerConfig.DirCharacters+"\"")
		return false
	}
	//Load locations
	fmt.Println("Loading character files:")
	for _, file := range files {
		data, err := ioutil.ReadFile(MainServerConfig.DirCharacters + "/" + file.Name());
		if (err == nil) {
			npcs := make([]Npc,0)
			err := json.Unmarshal(data,&npcs)
			if(err == nil) {
				fmt.Println("\t"+file.Name()+": loaded "+strconv.Itoa(len(npcs))+" characters")
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

func (npc *Npc) Heal(value float64) {
	npc.Health += value
	if(npc.Health > npc.MaxHealth) {
		npc.Health = npc.MaxHealth
	}
}

func (npc *Npc) Damage(value float64) {
	npc.Health -= value
	if(npc.Health <= 0) {
		npc.Health = 0
		npc.Die()
	}
}

func (npc *Npc) Die() {
	fmt.Println(npc.Name+" died!")
}
