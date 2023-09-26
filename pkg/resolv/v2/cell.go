package resolv

// Cell is used to contain and organize Object information.
type Cell[T comparable] struct {
	X, Y    int          // The X and Y position of the cell in the Space - note that this is in Grid position, not World position.
	Objects []*Object[T] // The Objects that a Cell contains.
}

// newCell creates a new cell at the specified X and Y position. Should not be used directly.
func newCell[T comparable](x, y int) *Cell[T] {
	return &Cell[T]{
		X:       x,
		Y:       y,
		Objects: []*Object[T]{},
	}
}

// register registers an object with a Cell. Should not be used directly.
func (cell *Cell[T]) register(obj *Object[T]) {
	if !cell.Contains(obj) {
		cell.Objects = append(cell.Objects, obj)
	}
}

// unregister unregisters an object from a Cell. Should not be used directly.
func (cell *Cell[T]) unregister(obj *Object[T]) {
	for i, o := range cell.Objects {
		if o == obj {
			cell.Objects[i] = cell.Objects[len(cell.Objects)-1]
			cell.Objects = cell.Objects[:len(cell.Objects)-1]
			break
		}
	}
}

// Contains returns whether a Cell contains the specified Object at its position.
func (cell *Cell[T]) Contains(obj *Object[T]) bool {
	for _, o := range cell.Objects {
		if o == obj {
			return true
		}
	}

	return false
}

// ContainsTags returns whether a Cell contains an Object that has the specified tag at its position.
func (cell *Cell[T]) ContainsTags(tags ...T) bool {
	for _, o := range cell.Objects {
		if o.HaveTags(tags...) {
			return true
		}
	}

	return false
}

// Occupied returns whether a Cell contains any Objects at all.
func (cell *Cell[T]) Occupied() bool {
	return len(cell.Objects) > 0
}
