package particle_test

import (
	"image/color"
	"reflect"
	"testing"
)

func assertEqual[T comparable](t testing.TB, a, b T) {
	t.Helper()
	if a == b {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
func assertColor(t testing.TB, a, b color.Color) {
	t.Helper()
	if a == b {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
