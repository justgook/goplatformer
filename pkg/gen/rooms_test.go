package gen_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoomsLayout(t *testing.T) {

}

func FuzzRoomsLayout(f *testing.F) {
	//f.Add(6, 3)
	//f.Fuzz(func(t *testing.T, goalDistance, branchLength int) {
	//
	//})
	f.Add(10)
	f.Fuzz(func(t *testing.T, n int) {
		n %= 20
		expect := 10
		aaa := n
		assert.Equal(t, expect, aaa)
	})
}
