package layers

import "github.com/yohamta/donburi/ecs"

const (
	Background ecs.LayerID = iota
	Enemy
	PlayerBullets
	Player
	EnemyBullets
	Foreground
	HUD
)
