package olduar

import (
	"testing"
)

const (
	MAX_ITERATIONS = 500
)

func TestAttributes(t *testing.T) {
	attrDamage := &Attribute{
		Id:          "damage",
		Name:        "Damage",
		Type:        "damage",
		Description: "Basic Damage",
		Groups:      &[]int{1},
	}
	attrDefense := &Attribute{
		Id:          "defense",
		Name:        "Defense",
		Type:        "resistance",
		Description: "Basic Defense",
		Groups:      &[]int{1},
	}
	attrLifeSteal := &Attribute{
		Id:          "life_steal",
		Name:        "Life Steal",
		Type:        "mod",
		Description: "Basic Defense",
		Groups:      &[]int{1},
		Config: map[string]interface{}{
			"heal": 1.0, //Some config must be present - otherwise Prepare() will fail
		},
	}

	AllAttributes["damage"] = attrDamage
	AllAttributes["defense"] = attrDefense
	AllAttributes["life_steal"] = attrLifeSteal

	if !attrDamage.Prepare() {
		t.Error("Unable to prepare DMG attribute")
	} else if !attrDefense.Prepare() {
		t.Error("Unable to prepare DEF attribute")
	} else if !attrLifeSteal.Prepare() {
		t.Error("Unable to prepare MOD attribute (Config must be present)")
	}

}

func TestCombatSimple(t *testing.T) {
	//Create room, player & enemy
	room := &Room{}
	player := &Player{
		Name:      "Dummy player",
		Health:    50,
		MaxHealth: 50,
		Stats:     AttributeList{"damage": MakeAttributeValue(10)},
		Room:      room,
		Inventory: make(Inventory, 0),
	}
	enemy := &Npc{
		Name:      "Dummy enemy",
		Health:    50,
		MaxHealth: 50,
		Stats:     AttributeList{"damage": MakeAttributeValue(5)},
		Friendly:  false,
	}

	//Create & start queue
	queue := MakeCombatQueue(room)
	queue.Add(player)
	queue.Add(enemy)
	queue.Start()

	//Process the queue
	i := MAX_ITERATIONS //Number of maximum combat iterations
	for queue.InProgress && i > 0 {

		if enemy.GetTeam() == queue.GetCurrentFighter().GetTeam() {
			//Enemies are attacking
			if player.IsAlive() {
				queue.Attack(player) //Player should be attacked
			}

		} else {
			//Players are attacking
			if enemy.IsAlive() {
				queue.Attack(enemy) //Enemy should be attacked
			}

		}

		i-- //Decrese number of iterations
	}

	if i == 0 {
		t.Error("Combat was in infinite loop")
	}
}

func TestCombatDefense(t *testing.T) {
	//Create room, player & enemy
	room := &Room{}
	player := &Player{
		Name:      "Dummy player",
		Health:    50,
		MaxHealth: 50,
		Stats:     AttributeList{"damage": MakeAttributeValue(10), "defense": MakeAttributeValue(5)},
		Room:      room,
		Inventory: make(Inventory, 0),
	}
	enemy := &Npc{
		Name:      "Dummy enemy",
		Health:    50,
		MaxHealth: 50,
		Stats:     AttributeList{"damage": MakeAttributeValue(5), "defense": MakeAttributeValue(0)},
		Friendly:  false,
	}

	//Create & start queue
	queue := MakeCombatQueue(room)
	queue.Add(player)
	queue.Add(enemy)
	queue.Start()

	//Check for player to start
	if queue.GetCurrentFighter().GetTeam() == 1 {
		queue.NextTurn()
	}

	//Process the queue
	i := MAX_ITERATIONS //Number of maximum combat iterations
	for queue.InProgress && i > 0 {

		if enemy.GetTeam() == queue.GetCurrentFighter().GetTeam() {
			//Enemies are attacking
			if player.IsAlive() {
				queue.Attack(player) //Player should be attacked
			}

		} else {
			//Players are attacking
			if enemy.IsAlive() {
				queue.Attack(enemy) //Enemy should be attacked
			}

		}

		i-- //Decrese number of iterations
	}

	//Turns check
	if i == 0 {
		t.Error("Combat was in infinite loop")
	} else if i < MAX_ITERATIONS-9 { //This should take 9 turns
		t.Error("Combat was longer than 9 turns")
	} else if i > MAX_ITERATIONS-9 { //This should take 9 turns
		t.Error("Combat was shorter than 9 turns")
	}

	//Health check
	if player.Health != player.MaxHealth {
		t.Error("Player should have 100% health")
	} else if enemy.Health != 0 {
		t.Error("Enemy should have 0% health")
	}
}

func TestCombatHeal(t *testing.T) {
	//Create room, player & enemy
	room := &Room{}
	player := &Player{
		Name:      "Dummy player",
		Health:    20, //Lower health - it should be at max after the combat
		MaxHealth: 50,
		Stats:     AttributeList{"damage": MakeAttributeValue(10), "defense": MakeAttributeValue(5), "life_steal": MakeAttributeValue(100)},
		Room:      room,
		Inventory: make(Inventory, 0),
	}
	enemy := &Npc{
		Name:      "Dummy enemy",
		Health:    50,
		MaxHealth: 50,
		Stats:     AttributeList{"damage": MakeAttributeValue(5)},
		Friendly:  false,
	}

	//Create & start queue
	queue := MakeCombatQueue(room)
	queue.Add(player)
	queue.Add(enemy)
	queue.Start()

	//Check for player to start
	if queue.GetCurrentFighter().GetTeam() == 1 {
		queue.NextTurn()
	}

	//Process the queue
	i := MAX_ITERATIONS //Number of maximum combat iterations
	for queue.InProgress && i > 0 {

		if enemy.GetTeam() == queue.GetCurrentFighter().GetTeam() {
			//Enemies are attacking
			if player.IsAlive() {
				queue.Attack(player) //Player should be attacked
			}

		} else {
			//Players are attacking
			if enemy.IsAlive() {
				queue.Attack(enemy) //Enemy should be attacked
			}

		}

		i-- //Decrese number of iterations
	}

	//Turns check
	if i == 0 {
		t.Error("Combat was in infinite loop")
	} else if i < MAX_ITERATIONS-9 { //This should take 9 turns
		t.Error("Combat was longer than 9 turns")
	} else if i > MAX_ITERATIONS-9 { //This should take 9 turns
		t.Error("Combat was shorter than 9 turns")
	}

	//Health check
	if player.Health != player.MaxHealth {
		t.Error("Player should have 100% health")
	} else if enemy.Health != 0 {
		t.Error("Enemy should have 0% health")
	}
}
