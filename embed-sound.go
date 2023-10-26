package goplatformer

import (
	"embed"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/justgook/goplatformer/pkg/util"
)

var (
	//go:embed asset/sfx/*.ogg
	soundsEmbeded embed.FS
	AllSounds     = func() map[string]*vorbis.Stream {
		pathPrefix := "asset/sfx"
		output := map[string]*vorbis.Stream{}
		files, err := soundsEmbeded.ReadDir(pathPrefix)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			filepath := filepath.Join(pathPrefix, f.Name())
			aaa := util.GetOrDie(soundsEmbeded.Open(filepath))
			defer aaa.Close()
			key := fileNameWithoutExtSliceNotation(f.Name())
			output[key] = util.GetOrDie(vorbis.DecodeWithoutResampling(aaa))
		}

		return output
	}()
)

func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
