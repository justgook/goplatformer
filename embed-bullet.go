package goplatformer

import (
	"bytes"
	_ "embed"

	"github.com/justgook/goplatformer/pkg/bulletml"
)

var (
	//go:embed asset/bullet/test.xml
	testBullet []byte
	TestBullet = func() *bulletml.BulletML {
		f := bytes.NewReader(testBullet)

		bml, err := bulletml.Load(f)
		if err != nil {
			panic(err)
		}
		return bml
	}()
)
