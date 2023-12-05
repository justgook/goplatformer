package domain

import "github.com/justgook/goplatformer/pkg/util"

type RoomNavigation = util.Bits

const (
	RoomNavigationExitN RoomNavigation = 1 << iota
	RoomNavigationExitE
	RoomNavigationExitS
	RoomNavigationExitW
	RoomNavigationStart
	RoomNavigationGoal
)
