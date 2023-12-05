package domain

import "github.com/justgook/goplatformer/pkg/resolv/v2"

type ObjectTag = int64
type ObjectSpace = resolv.Space[ObjectTag]
type Object = resolv.Object[ObjectTag]

const (
	ObjectTagSolid            ObjectTag = 1
	ObjectTagOneWayUp         ObjectTag = 3
	ObjectTagLava             ObjectTag = 4
	ObjectTagRamp             ObjectTag = 6
	ObjectTagEnemyTrigger     ObjectTag = 7
	ObjectTagExitTriggerNorth ObjectTag = 9
	ObjectTagExitTriggerEast  ObjectTag = 10
	ObjectTagExitTriggerSouth ObjectTag = 11
	ObjectTagExitTriggerWest  ObjectTag = 12
	ObjectTagGoal             ObjectTag = 13
)
