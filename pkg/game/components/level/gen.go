package level

import (
	"github.com/justgook/goplatformer/pkg/pcg"
	"github.com/justgook/goplatformer/pkg/util"
)

type TmpRoom struct {
	RoomMetaData
	depth       uint
	branchDepth uint
}

func generateData(rnd *pcg.PCG32, goalDistance, branchLength int) *Data {
	goalDepth := uint(goalDistance) + 2
	maxBrunchDepth := uint(branchLength)

	output := &Data{RoomsInfo: []RoomMetaData{}}
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
			RoomMetaData: RoomMetaData{RoomPoint: nextXY},
			depth:        i,
		}
		minX = min(minX, nextXY.X)
		minY = min(minY, nextXY.Y)
		lastXY = nextXY

		lastLeft, nextXY, fromRoom, exits = firstPass(rnd, lastLeft, nextXY, fromRoom)
		// add random exits
		room.Exits = addRandExits(rnd, goalDepth, room.depth, exits)
		roomMap[RoomPoint{X: room.X, Y: room.Y}] = room
		if room.Exits != exits {
			unfinished[RoomPoint{X: room.X, Y: room.Y}] = room
		}
	}
	roomMap[RoomPoint{X: 0, Y: 0}].Start = true
	roomMap[lastXY].Goal = true

	// Trim goal exits
	trimExits(lastXY, roomMap)
	delete(unfinished, lastXY)

	// Second pass - all brunches
	for len(unfinished) > 0 {
		for k := range unfinished {
			forEachExit(unfinished[k].Exits, func(e util.Bits) {
				newXY := neighbor(k, e)
				v, ok := roomMap[newXY]
				if ok {
					if !v.Exits.Has(e) {
						v.Exits = v.Exits.Set(oppositeExit(e))
					}

					return
				}
				depth := unfinished[k].depth + 1
				branchDepth := uint(1)
				if roomMap[k].branchDepth > 0 {
					branchDepth = roomMap[k].branchDepth + 1
				}
				room := &TmpRoom{
					RoomMetaData: RoomMetaData{RoomPoint: newXY},
					depth:        depth,
					branchDepth:  branchDepth,
				}
				minX = min(minX, newXY.X)
				minY = min(minY, newXY.Y)
				room.Exits = oppositeExit(e)
				room.Exits = addRandExits(rnd, goalDepth, room.depth, room.Exits)
				roomMap[newXY] = room
				if branchDepth >= maxBrunchDepth {
					return
				}
				unfinished[RoomPoint{X: room.X, Y: room.Y}] = room
			})
			delete(unfinished, k)
		}
	}

	for k := range roomMap {
		trimExits(k, roomMap)
		item := roomMap[k]
		item.RoomMetaData.X -= minX
		item.RoomMetaData.Y -= minY
		if item.RoomMetaData.Start {
			output.Start = item.RoomPoint
		}
		if item.RoomMetaData.Goal {
			output.Goal = item.RoomPoint
		}
		output.RoomsInfo = append(output.RoomsInfo, item.RoomMetaData)
	}

	return output
}

func trimExits(p RoomPoint, roomMap map[RoomPoint]*TmpRoom) {
	forEachExit(roomMap[p].Exits, func(e util.Bits) {
		newXY := neighbor(p, e)
		if v, ok := roomMap[newXY]; ok && v.Exits.Has(oppositeExit(e)) {
			return
		}
		roomMap[p].Exits = roomMap[p].Exits.Clear(e)
	})
}

func neighbor(p RoomPoint, now util.Bits) RoomPoint {
	switch now {
	case ExitN:
		p.Y -= 1
	case ExitE:
		p.X += 1
	case ExitS:
		p.Y += 1
	case ExitW:
		p.X -= 1
	}
	return p
}

func forEachExit(exits util.Bits, fn func(util.Bits)) {
	if exits.Has(ExitN) {
		fn(ExitN)
	}
	if exits.Has(ExitE) {
		fn(ExitE)
	}
	if exits.Has(ExitS) {
		fn(ExitS)
	}
	if exits.Has(ExitW) {
		fn(ExitW)
	}
}

func addRandExits(rnd *pcg.PCG32, maxDepth, depth uint, exits util.Bits) util.Bits {
	return exits.Set(util.Bits(rnd.Bounded(15) + 1))
}

func oppositeExit(now util.Bits) util.Bits {
	switch now {
	case ExitN:
		return ExitS
	case ExitE:
		return ExitW
	case ExitS:
		return ExitN
	case ExitW:
		return ExitE
	}
	return ExitN

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
		return ExitN
	case 1:
		return ExitE
	case 2:
		return ExitS
	case 3:
		return ExitW
	}

	return ExitN
}

func firstPassLastLeft(lastLeft bool, was, now util.Bits) bool {
	type pair struct{ was, now util.Bits }
	switch (pair{was: was, now: now}) {
	case pair{was: ExitN, now: ExitE}:
		lastLeft = true
	case pair{was: ExitN, now: ExitW}:
		lastLeft = false
	case pair{was: ExitE, now: ExitS}:
		lastLeft = true
	case pair{was: ExitE, now: ExitN}:
		lastLeft = false
	case pair{was: ExitS, now: ExitW}:
		lastLeft = true
	case pair{was: ExitS, now: ExitE}:
		lastLeft = false
	case pair{was: ExitW, now: ExitN}:
		lastLeft = true
	case pair{was: ExitW, now: ExitS}:
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
	case ExitN:
		valid = !(now == ExitE && lastLeft) && !(now == ExitW && !lastLeft)
	case ExitE:
		valid = !(now == ExitS && lastLeft) && !(now == ExitN && !lastLeft)
	case ExitS:
		valid = !(now == ExitW && lastLeft) && !(now == ExitE && !lastLeft)
	case ExitW:
		valid = !(now == ExitN && lastLeft) && !(now == ExitS && !lastLeft)
	}

	return valid
}

func firstPassRotate(lastLeft bool, now util.Bits) util.Bits {
	return firstPassRotateCClock(now)

}

func firstPassRotateCClock(now util.Bits) util.Bits {
	switch now {
	case ExitN:
		now = ExitW
	case ExitE:
		now = ExitN
	case ExitS:
		now = ExitE
	case ExitW:
		now = ExitS
	}

	return now
}
