package olduar

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"math/rand"
)

const (
	PARAM_DEFAULT = 0

	PARAM_DAMAGE = 0
	PARAM_HEALING = 1
	PARAM_ABILITY = 2
)

var AllAttributes map[string]*Attribute = make(map[string]*Attribute)

type DamageGroups map[int]float64

type Attribute struct {
	Id string 						`json:"id"`
	Name string 					`json:"name"`
	Type string 					`json:"type"`
	Description string 				`json:"desc"`
	Groups *[]int					`json:"groups"`

	Config map[string]interface{} 	`json:"config"`
	worker AttributeWorker 			`json:"-"`
}

func LoadAttributes() {
	data, err := ioutil.ReadFile(MainServerConfig.DirOther+"/attributes.json");
	if(err == nil) {
		attributes := make([]*Attribute, 0)
		err := json.Unmarshal(data, &attributes)
		if(err == nil) {
			for _, attribute := range attributes {
				if(attribute.Prepare()) {
					AllAttributes[attribute.Id] = attribute
				}
			}
			fmt.Println("Loaded "+strconv.Itoa(len(AllAttributes))+" attributes")
		} else {
			fmt.Println("Failed to parse attributes ("+MainServerConfig.DirOther+"/attributes.json)")
		}
	} else {
		fmt.Println("Failed to load attributes ("+MainServerConfig.DirOther+"/attributes.json)")
	}
}

func (a *Attribute) MatchGroup(b *Attribute) bool {
	if(a.Groups == nil || b.Groups == nil) {
		return true
	}
	for _, groupA := range *a.Groups {
		for _, groupB := range *b.Groups {
			if(groupA == groupB) {
				return true
			}
		}
	}
	return false
}

func (a *Attribute) Prepare() bool {
	switch(a.Type) {
	case "damage":
		a.worker = new(DamageAttribute)
		return a.worker.Prepare(a.Config)
	case "mod":
		a.worker = new(ModAttribute)
		return a.worker.Prepare(a.Config)
	case "resistance":
		a.worker = new(ResistanceAttribute)
		return a.worker.Prepare(a.Config)
	}
	return false
}

type AttributeWorker interface {
	Prepare(map[string]interface{}) bool
	Compute(float64,float64,*Room,int) (float64)
}

//Attribute value
type AttributeValue struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
func (value AttributeValue) Add(stat AttributeValue) {
	value.Min += stat.Min
	value.Max += stat.Max
}
func (value AttributeValue) Value() float64 {
	return value.Min + (rand.Float64() * (value.Max - value.Min))
}
func MakeAttributeValue(val float64) AttributeValue {
	return AttributeValue{Min:val,Max:val}
}
func MakeAttributeValueMinMax(min float64, max float64) AttributeValue {
	return AttributeValue{Min:min,Max:max}
}

//List of attributes & values
type AttributeList map[string]AttributeValue

func (list *AttributeList) Reset() {
	*list = make(AttributeList)
}

func (list *AttributeList) Append(list2 AttributeList) {
	for key, value := range list2 {
		_, found := (*list)[key]
		if(found) {
			(*list)[key].Add(value)
		} else {
			(*list)[key] = value
		}
	}
}

func (attacker AttributeList) Attack(target AttributeList, room *Room) (float64, float64) {
	dmgTarget, healAttacker := 0.0, 0.0

	//Go trough all attributes
	for nameDamage, value := range attacker {
		attributeDamage, found := AllAttributes[nameDamage]
		healValue := 0.0
		damageValue := value.Value()

		//Check if found & attributeA is "damage" type
		if(found && attributeDamage.Type == "damage") {

			//Cycle trough my damage-mod attributes
			for nameMod, value := range attacker {
				modValue := value.Value()
				attributeMod, found := AllAttributes[nameMod]
				if(found && attributeMod.Type == "mod" && attributeDamage.MatchGroup(attributeMod)) {
					damageValue = attributeMod.worker.Compute(modValue,damageValue,room,PARAM_DAMAGE)
				}
			}

			//Cycle trough resistance attributes of target for damage reduction
			for nameRes, value := range target {
				resistanceValue := value.Value()
				attributeRes, found := AllAttributes[nameRes]
				if(found && attributeRes.Type == "resistance" && attributeDamage.MatchGroup(attributeRes)) {
					damageValue = attributeRes.worker.Compute(resistanceValue,damageValue,room,PARAM_DEFAULT)
				}
			}

			//Cycle trough my damage-mod attributes (again with param 1 = After resistance)
			for nameMod, value := range attacker {
				modValue := value.Value()
				attributeMod, found := AllAttributes[nameMod]
				if(found && attributeMod.Type == "mod" && attributeDamage.MatchGroup(attributeMod)) {
					healValue += attributeMod.worker.Compute(modValue,damageValue,room,PARAM_HEALING)
				}
			}

			//Add damageValue to total damage
			dmgTarget += damageValue
			healAttacker += healValue
		}

	}

	return dmgTarget, healAttacker
}

//Damage attribute type
type DamageAttribute struct {
	Message string
}

func (w *DamageAttribute) Prepare(config map[string]interface{}) bool {
	data, found := config["msg"]
	msg := fmt.Sprint(data)
	if(found && msg != "") {
		w.Message = msg
	}
	return true
}

func (w *DamageAttribute) Compute(damage float64, _ float64, room *Room, param int) (float64) {
	return damage
}

//Damage-mod attribute type
type ModAttribute struct {
	Message string
	DamageValue, HealingValue, AbilityValue float64
}

func (w *ModAttribute) Prepare(config map[string]interface{}) bool {
	//Message
	data, found := config["msg"]
	msg := fmt.Sprint(data)
	if(found && msg != "") {
		w.Message = msg
	}

	//Damage
	data, found = config["damage"]
	if(found) {
		w.DamageValue = data.(float64)
	}

	//Healing
	data, found = config["heal"]
	if(found) {
		w.HealingValue = data.(float64)
	}

	//Damage
	data, found = config["ability"]
	if(found) {
		w.AbilityValue = data.(float64)
	}

	//This attribute must have at least healing or damage value above 0 (default)
	return w.HealingValue > 0.0 || w.DamageValue > 0.0 || w.AbilityValue > 0.0
}

func (w *ModAttribute) Compute(mod float64, damage float64, room *Room, param int) (float64) {
	switch(param) {
	case PARAM_DAMAGE:
		return damage + (mod * w.DamageValue)
	case PARAM_HEALING:
		return damage * (mod * 0.01) * w.HealingValue
	case PARAM_ABILITY:
		return damage * mod * w.AbilityValue
	}
	return 0
}

//Resistance attribute type
type ResistanceAttribute struct {
	Message string
}

func (w *ResistanceAttribute) Prepare(config map[string]interface{}) bool {
	data, found := config["msg"]
	msg := fmt.Sprint(data)
	if(found && msg != "") {
		w.Message = msg
	}
	return true
}

func (w *ResistanceAttribute) Compute(resistance float64, damage float64, room *Room, param int) (float64) {
	damage -= resistance
	if(damage < 0) {
		return 0
	}
	return damage
}
