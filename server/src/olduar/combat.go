package olduar

import (
	"fmt"
	"math/rand"
)

type Fighter interface {
	GetStats() AttributeList
	Damage(float64)
	Heal(float64)
	Die()
	GetTeam() CombatTeam
	GetName() string
	GetId() string
	GetGUID() GUID
	IsAlive() bool
	IsPlayer() bool
}

type CombatTeam int8

type CombatQueue struct {
	room *Room
	combatants map[Fighter]CombatTeam
	queue []Fighter
	queuePosition int
	InProgress bool
}

func MakeCombatQueue(room *Room) *CombatQueue {
	return &CombatQueue{
		room: room,
		combatants: make(map[Fighter]CombatTeam),
		queue: make([]Fighter,0),
		queuePosition: 0,
		InProgress: false,
	}
}

func (q *CombatQueue) Add(combatant Fighter) {
	_, found := q.combatants[combatant]
	if(!found) {
		q.combatants[combatant] = combatant.GetTeam()
		if(q.InProgress) {
			q.recomputeQueue()
		}
	}
}

func (q *CombatQueue) Available() bool {
	//Check team status
	teams := make(map[CombatTeam]int)
	for combatant, team := range q.combatants {
		if(combatant.IsAlive()) {
			teams[team]++
		}
	}

	//If only last team is standing - combat ends
	if(len(teams) <= 1) {
		return false
	}
	return true
}

func (q *CombatQueue) Start() {
	q.InProgress = true
	q.recomputeQueue()
	fmt.Println("Combat started!")
	fmt.Println("-----------------------------")
}

func (q *CombatQueue) End() {
	q.InProgress = false
	fmt.Println("-----------------------------")
	fmt.Println("Combat ended!")
}

func (q *CombatQueue) recomputeQueue() {
	//Count teams
	max, total, index := 0, 0, 0
	teams := make(map[CombatTeam]int)
	for _, team := range q.combatants {
		teams[team]++
		if(teams[team] > max) {
			max = teams[team]
		}
	}

	//Prepare ratio and total count
	teamIds := make([]CombatTeam,0)
	for team, count := range teams {
		teams[team] = max/count
		teamIds = append(teamIds,team)
	}
	teamCombatants := make(map[CombatTeam]map[Fighter]int)
	for combatant, team := range q.combatants {
		total += teams[team]
		_, found := teamCombatants[team]
		if(!found) {
			teamCombatants[team] = make(map[Fighter]int)
		}
		teamCombatants[team][combatant] = teams[team]
	}

	//Prepare queue
	teamCount := len(teamIds)
	remaining, team := total, rand.Intn(teamCount)
	index = 0
	q.queue = make([]Fighter,total)
	for remaining > 0 {
		//Get list of fighters in group
		teamId := teamIds[team]
		var fighters []Fighter
		for k, count := range teamCombatants[teamId] {
			if(count > 0){
				fighters = append(fighters, k)
			}
		}

		team++
		if(team == teamCount) {
			team = 0
		}

		var fightersLength = len(fighters)
		if(fightersLength > 0) {
			fighter := fighters[rand.Intn(fightersLength)]
			teamCombatants[teamId][fighter]--
			q.queue[index] = fighter
			remaining--
			index++
		}

	}
}

func (q *CombatQueue) GetCurrentFighter() Fighter {
	if(!q.InProgress) {
		return nil
	}
	return q.queue[q.queuePosition]
}

func (q *CombatQueue) NextTurn() {
	if(q.InProgress) {
		fmt.Println(q.GetCurrentFighter().GetName() + "> Next turn")

		//Advance
		q.queuePosition++
		if (q.queuePosition >= len(q.queue)) {
			q.queuePosition = 0
		}

		if(!q.Available()) {
			q.End()
		}

	}
}

func (q *CombatQueue) Defend() {
	if(q.InProgress) {
		fmt.Println(q.GetCurrentFighter().GetName() + "> Defend")
		q.NextTurn()
	}
}

func (q *CombatQueue) Attack(enemy Fighter) bool {
	if(!q.InProgress) {
		return false
	}

	//Get attacker
	attacker := q.GetCurrentFighter()

	//Friendly fire check
	if(enemy == nil || attacker.GetTeam() == enemy.GetTeam()) {
		return false
	}

	//Do actual attack
	damage, heal := attacker.GetStats().Attack(enemy.GetStats(),q.room)
	fmt.Println(attacker.GetName() + "> Attacked ("+enemy.GetName()+") for",damage," damage")
	attacker.Heal(heal)
	enemy.Damage(damage)

	//Advance to next turn
	q.NextTurn()
	return true
}

func (q *CombatQueue) MakeAutoTurn() {
	attacker := q.GetCurrentFighter();
	if(attacker.IsPlayer()) {
		q.Defend() //Player will always defend only on auto-turn (which is timeout in this case)

	} else {
		for combatant, team := range q.combatants {
			if(attacker.GetTeam() != team) {
				//TODO: Actually think about what you do
				q.Attack(combatant)
				return
			}
		}
	}
}
