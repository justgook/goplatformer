package gen

import (
	"image"

	. "github.com/justgook/goplatformer/pkg/core/domain"
	"github.com/justgook/goplatformer/pkg/pcg"
	"github.com/justgook/goplatformer/pkg/util"
)

type RoomPoint = image.Point

type LayoutData struct {
	RoomsLayout []RoomMetaData
	Goal        image.Point
	Start       image.Point
}
type RoomMetaData struct {
	RoomPoint
	RoomNavigation
}

type TmpRoom struct {
	*RoomMetaData
	depth       uint
	branchDepth uint
}

// RoomsLayoutFN call for each room in level in random order, but only one time for each room
type RoomsLayoutFN = func(data *RoomMetaData, start image.Point, goal image.Point)

// RoomsLayout generates random map,
// `goalDistance` - distance between start and end of level
// `branchLength`  - how far away path can branch from main road
// `callback` - function that is called for each room
func RoomsLayout(
	rnd *pcg.PCG32,
	goalDistance int,
	branchLength int,
	callback RoomsLayoutFN,
) {
	goalDepth := uint(goalDistance) + 2
	maxBrunchDepth := uint(branchLength)

	minX := 0
	minY := 0
	fromRoom := util.Bits(0)

	nextXY := RoomPoint{X: 0, Y: 0}

	lastLeft := rnd.Bounded(2) > 0
	unfinished := make(map[RoomPoint]*TmpRoom)
	roomMap := make(map[RoomPoint]*TmpRoom)
	lastXY := nextXY

	// First pass - from start to finish
	for i := uint(0); i < goalDepth; i++ {
		exits := util.Bits(0)
		room := &TmpRoom{
			RoomMetaData: &RoomMetaData{
				RoomPoint: nextXY,
			},
			depth: i,
		}
		minX = min(minX, nextXY.X)
		minY = min(minY, nextXY.Y)
		lastXY = nextXY

		lastLeft, nextXY, fromRoom, exits = firstPass(rnd, lastLeft, nextXY, fromRoom)
		// add random exits
		room.RoomNavigation = addRandExits(rnd, goalDepth, room.depth, exits)
		roomMap[RoomPoint{X: room.X, Y: room.Y}] = room
		if room.RoomNavigation != exits {
			unfinished[RoomPoint{X: room.X, Y: room.Y}] = room
		}
	}

	startPoint := RoomPoint{X: 0, Y: 0}
	roomMap[startPoint].RoomNavigation = roomMap[startPoint].RoomNavigation.Set(RoomNavigationStart)
	roomMap[lastXY].RoomNavigation = roomMap[lastXY].RoomNavigation.Set(RoomNavigationGoal)

	// Trim goal exits
	trimExits(lastXY, roomMap)
	delete(unfinished, lastXY)

	// Second pass - all brunches
	for len(unfinished) > 0 {
		for k := range unfinished {
			forEachExit(unfinished[k].RoomNavigation, func(e util.Bits) {
				newXY := neighbor(k, e)
				v, ok := roomMap[newXY]
				if ok {
					if !v.RoomNavigation.Has(e) {
						v.RoomNavigation = v.RoomNavigation.Set(oppositeExit(e))
					}

					return
				}
				depth := unfinished[k].depth + 1
				branchDepth := uint(1)
				if roomMap[k].branchDepth > 0 {
					branchDepth = roomMap[k].branchDepth + 1
				}
				room := &TmpRoom{
					RoomMetaData: &RoomMetaData{RoomPoint: newXY},
					depth:        depth,
					branchDepth:  branchDepth,
				}
				minX = min(minX, newXY.X)
				minY = min(minY, newXY.Y)
				room.RoomNavigation = oppositeExit(e)
				room.RoomNavigation = addRandExits(rnd, goalDepth, room.depth, room.RoomNavigation)
				roomMap[newXY] = room
				if branchDepth >= maxBrunchDepth {
					return
				}
				unfinished[RoomPoint{X: room.X, Y: room.Y}] = room
			})
			delete(unfinished, k)
		}
	}

	minXY := RoomPoint{X: minX, Y: minY}
	startPint := RoomPoint{X: 0, Y: 0}.Sub(minXY)
	goalPoint := lastXY.Sub(minXY)

	for k := range roomMap {
		trimExits(k, roomMap)
		item := roomMap[k]
		item.RoomMetaData.RoomPoint = item.RoomMetaData.RoomPoint.Sub(minXY)
		callback(item.RoomMetaData, startPint, goalPoint)
	}
}

func trimExits(p RoomPoint, roomMap map[RoomPoint]*TmpRoom) {
	forEachExit(roomMap[p].RoomNavigation, func(e util.Bits) {
		newXY := neighbor(p, e)
		if v, ok := roomMap[newXY]; ok && v.RoomNavigation.Has(oppositeExit(e)) {
			return
		}
		roomMap[p].RoomNavigation = roomMap[p].RoomNavigation.Clear(e)
	})
}

func neighbor(p RoomPoint, now util.Bits) RoomPoint {
	switch now {
	case RoomNavigationExitN:
		p.Y -= 1
	case RoomNavigationExitE:
		p.X += 1
	case RoomNavigationExitS:
		p.Y += 1
	case RoomNavigationExitW:
		p.X -= 1
	}
	return p
}

func forEachExit(exits util.Bits, fn func(util.Bits)) {
	if exits.Has(RoomNavigationExitN) {
		fn(RoomNavigationExitN)
	}
	if exits.Has(RoomNavigationExitE) {
		fn(RoomNavigationExitE)
	}
	if exits.Has(RoomNavigationExitS) {
		fn(RoomNavigationExitS)
	}
	if exits.Has(RoomNavigationExitW) {
		fn(RoomNavigationExitW)
	}
}

func addRandExits(rnd *pcg.PCG32, maxDepth, depth uint, exits util.Bits) util.Bits {
	return exits.Set(util.Bits(rnd.Bounded(15) + 1))
}

func oppositeExit(now util.Bits) util.Bits {
	switch now {
	case RoomNavigationExitN:
		return RoomNavigationExitS
	case RoomNavigationExitE:
		return RoomNavigationExitW
	case RoomNavigationExitS:
		return RoomNavigationExitN
	case RoomNavigationExitW:
		return RoomNavigationExitE
	}
	return RoomNavigationExitN

}

func firstPass(rnd *pcg.PCG32, lastLeft bool, nextXY RoomPoint, fromRoom util.Bits) (bool, RoomPoint, util.Bits, util.Bits) {
	was := fromRoom
	now := intToExit(rnd.Bounded(4))
	now = firstPassApplyRules(lastLeft, was, now)
	nextXY = neighbor(nextXY, now)
	fromRoom = oppositeExit(now)
	exits := was.Set(now)

	return firstPassLastLeft(lastLeft, was, now), nextXY, fromRoom, exits
}

func intToExit(n uint32) util.Bits {
	switch n {
	case 0:
		return RoomNavigationExitN
	case 1:
		return RoomNavigationExitE
	case 2:
		return RoomNavigationExitS
	case 3:
		return RoomNavigationExitW
	}

	return RoomNavigationExitN
}

func firstPassLastLeft(lastLeft bool, was, now util.Bits) bool {
	type pair struct{ was, now util.Bits }
	switch (pair{was: was, now: now}) {
	case pair{was: RoomNavigationExitN, now: RoomNavigationExitE}:
		lastLeft = true
	case pair{was: RoomNavigationExitN, now: RoomNavigationExitW}:
		lastLeft = false
	case pair{was: RoomNavigationExitE, now: RoomNavigationExitS}:
		lastLeft = true
	case pair{was: RoomNavigationExitE, now: RoomNavigationExitN}:
		lastLeft = false
	case pair{was: RoomNavigationExitS, now: RoomNavigationExitW}:
		lastLeft = true
	case pair{was: RoomNavigationExitS, now: RoomNavigationExitE}:
		lastLeft = false
	case pair{was: RoomNavigationExitW, now: RoomNavigationExitN}:
		lastLeft = true
	case pair{was: RoomNavigationExitW, now: RoomNavigationExitS}:
		lastLeft = false
	}

	return lastLeft
}

func firstPassApplyRules(lastLeft bool, was, now util.Bits) util.Bits {
	if was == now {
		now = firstPassRotate(lastLeft, now)
	}

	if !firstPassValidateTurn(lastLeft, was, now) {
		// slog.Info("not valid turn")
		return firstPassApplyRules(lastLeft, was, firstPassRotate(lastLeft, now))
	}

	return now
}

func firstPassValidateTurn(lastLeft bool, was, now util.Bits) bool {
	valid := true
	switch was {
	case RoomNavigationExitN:
		valid = !(now == RoomNavigationExitE && lastLeft) && !(now == RoomNavigationExitW && !lastLeft)
	case RoomNavigationExitE:
		valid = !(now == RoomNavigationExitS && lastLeft) && !(now == RoomNavigationExitN && !lastLeft)
	case RoomNavigationExitS:
		valid = !(now == RoomNavigationExitW && lastLeft) && !(now == RoomNavigationExitE && !lastLeft)
	case RoomNavigationExitW:
		valid = !(now == RoomNavigationExitN && lastLeft) && !(now == RoomNavigationExitS && !lastLeft)
	}

	return valid
}

func firstPassRotate(lastLeft bool, now util.Bits) util.Bits {
	return firstPassRotateCClock(now)

}

func firstPassRotateCClock(now util.Bits) util.Bits {
	switch now {
	case RoomNavigationExitN:
		now = RoomNavigationExitW
	case RoomNavigationExitE:
		now = RoomNavigationExitN
	case RoomNavigationExitS:
		now = RoomNavigationExitE
	case RoomNavigationExitW:
		now = RoomNavigationExitS
	}

	return now
}
