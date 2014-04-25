package olduar

import (
	"fmt"
)

func RunTests() {
	return //Comment/Remove this to run tests
	fmt.Println("\nRunning tests")
	testResult("Combat",testCombat())
	fmt.Println()
}

func testResult(name string, result bool) {
	if(result) {
		fmt.Println("\t"+name+": success")

	} else {
		fmt.Println("\t"+name+": failed")
	}
}

func testCombat() bool {
	//Create room, player & enemy
	room := &Room{}
	player1 := &Player{
		Name: "Tester 1",
		Health: 50,
		MaxHealth: 50,
		Stats: AttributeList{ "damage":MakeAttributeValue(10), "defense":MakeAttributeValueMinMax(0,10)},
		Room: room,
		Inventory: make(Inventory,0),
	}
	player2 := &Player{
		Name: "Tester 2",
		Health: 50,
		MaxHealth: 50,
		Stats: AttributeList{ "damage":MakeAttributeValue(10), "defense":MakeAttributeValueMinMax(0,10)},
		Room: room,
		Inventory: make(Inventory,0),
	}
	player3 := &Player{
		Name: "Tester 3",
		Health: 50,
		MaxHealth: 50,
		Stats: AttributeList{ "damage":MakeAttributeValue(10), "defense":MakeAttributeValueMinMax(0,10)},
		Room: room,
		Inventory: make(Inventory,0),
	}
	enemy1 := &Npc{
		Name: "Dummy enemy 1",
		Health: 50,
		MaxHealth: 50,
		Stats: AttributeList{ "damage":MakeAttributeValue(10), "defense":MakeAttributeValueMinMax(0,10)},
		Friendly: false,
	}
	enemy2 := &Npc{
		Name: "Dummy enemy 2",
		Health: 50,
		MaxHealth: 50,
		Stats: AttributeList{ "damage":MakeAttributeValue(10), "defense":MakeAttributeValueMinMax(0,10)},
		Friendly: false,
	}

	//Create & start queue
	queue := MakeCombatQueue(room)
	queue.Add(player1)
	queue.Add(player2)
	queue.Add(player3)
	queue.Add(enemy1)
	queue.Add(enemy2)
	queue.Start()

	for queue.InProgress {
		if(enemy1.GetTeam() == queue.GetCurrentFighter().GetTeam()) {
			//Enemies are attacking
			if(player1.IsAlive()) {
				queue.Attack(player1)
			} else if(player1.IsAlive()) {
				queue.Attack(player2)
			} else {
				queue.Attack(player3)
			}

		} else {
			//Players are attacking
			if(enemy1.IsAlive()) {
				queue.Attack(enemy1)
			} else {
				queue.Attack(enemy2)
			}

		}
	}

	return true
}
