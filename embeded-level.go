package goplatformer

import (
	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
)

var EmbeddedLevel = func() *resources.Level {
	output := &resources.Level{}
	util.OrDie(resources.Load(embeddedLevel, output))

	return output
}()
