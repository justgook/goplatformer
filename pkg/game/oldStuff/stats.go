package oldStuff

type Stats struct {
	Power                        int // Increase weapon damage on attack. (%)
	TrueDamage                   int // Ignores all enemy defense, dealing fixed damage.
	DamageAgainstStrongerEnemies int // Increase percentage damage to Bosses and Mini-Bosses.
	DamageAgainstWeakerEnemies   int // Increase percentage damage to all other enemies.
	DashPower                    int // Percentage of weapon damage that a dash attack deals.
	CriticalChance               int // Increases percentage chance to deal critical damage on enemies up to 75%.
	CriticalDamage               int // The percentage of additional critical damage.
	AimAccuracy                  int // Increases ranged weapon accuracy.
	Defense                      int // Decrease damage from enemies and traps up to 75%. (Excluding debuffs - Poison and Fire)
	Toughness                    int // Reduces a fixed amount of damage from enemy or trap damage. (Excluding debuffs - Poison and Burn)
	Block                        int // Increases percentage chance to block enemy or trap damage.
	BlockRangeDamage             int //Increases percentage chance to block enemy Range attack.
	Evasion                      int // Increases percentage chance to dodge enemy or traps damage up to 75%.
	AttackSpeed                  int // Number of attacks possible in one second.
	ReloadSpeed                  int // The amount of time it takes to reload.
	MoveSpeed                    int // Increases movement speed.
	Dash                         int // Increase the number of dashes.
	Satiety                      int // Decrease satiety at the entrance to the room.
	EnhanceBurnDamage            int // Increase damage from Burn debuff on the enemy.
	EnhancePoisonDamage          int // Increase damage from Poison debuff on the enemy.
	EnhanceColdDamage            int // Slows down enemies, increases debuff time.
	EnhanceShockDamage           int // Renders enemy's defense useless, increases debuff time.
	EnhanceStunDamage            int // Temporarily disables the enemy, increases debuff time.

}
