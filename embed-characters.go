package goplatformer

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/components"
	"github.com/justgook/goplatformer/pkg/game/components/sprite"
	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
	"image"
	_ "image/png"
)

var EmbeddedPlayerAnimation = util.GetOrDie(characterDataFromAnimatedSpriteBytes(playerAnim))

func characterDataFromAnimatedSpriteBytes(inputBytes []byte) (*components.CharacterData, error) {
	input := &resources.AnimatedSprite{}
	if err := resources.Load(inputBytes, input); err != nil {
		return nil, util.Catch(err)
	}
	img, _, err := image.Decode(bytes.NewReader(input.Image))
	if err != nil {
		return nil, util.Catch(err)
	}
	img2 := ebiten.NewImageFromImage(img)

	output := &components.CharacterData{
		Idle:         sprite.New(img2, input.Anim["Idle"].Frames, input.Anim["Run"].Durations),
		Run:          sprite.New(img2, input.Anim["Run"].Frames, input.Anim["Run"].Durations),
		JumpUp:       sprite.New(img2, input.Anim["JumpUp"].Frames, input.Anim["JumpUp"].Durations),
		JumpMax:      sprite.New(img2, input.Anim["JumpMax"].Frames, input.Anim["JumpMax"].Durations),
		Fall:         sprite.New(img2, input.Anim["Fall"].Frames, input.Anim["Fall"].Durations),
		WallSlideLow: sprite.New(img2, input.Anim["wSlideLow"].Frames, input.Anim["wSlideLow"].Durations),
	}
	output.Current = output.Run

	return output, nil
}
