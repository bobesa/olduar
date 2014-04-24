package olduar

type Fighter interface {
	GetStats() AttributeList
	Damage(float64)
	Heal(float64)
	Die()
}

func CombatAttack(room *Room, attacker Fighter, enemy Fighter) {
	damage, heal := attacker.GetStats().Attack(enemy.GetStats(),room)
	attacker.Heal(heal)
	enemy.Damage(damage)
}
