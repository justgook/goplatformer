package system

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game/state"
	"github.com/justgook/goplatformer/pkg/resolv/v2"
	"github.com/justgook/goplatformer/pkg/resources"
	"github.com/justgook/goplatformer/pkg/util"
)

var _ state.Scene = (*LevelManager)(nil)

type RoomExit = int64

const (
	ExitNorth RoomExit = 5
	ExitEast  RoomExit = 6
	ExitSouth RoomExit = 7
	ExitWest  RoomExit = 8
)

type LevelManager struct {
	Space *resolv.Space[resources.TagType]

	currentLevel *resources.Level
	tileSet      *ebiten.Image

	roomRender *ebiten.Image
}

// Draw implements System.
func (lm *LevelManager) Draw(screen *ebiten.Image) {
	screen.DrawImage(lm.roomRender, nil)
}

// Init implements System.
func (lm *LevelManager) Init() {
	lm.currentLevel = &resources.Level{}
}

func (lm *LevelManager) Load(input []byte) error {
	target := &resources.Level{}
	if err := resources.Load(input, target); err != nil {
		return util.Catch(err)
	}
	img, err := util.LoadImage(target.Image)
	if err != nil {
		return util.Catch(err)
	}
	lm.tileSet = img

	// TODO Update To level generation
	lm.currentLevel = target
	if len(target.Rooms) > 0 {
		lm.newRoom(target.Rooms[0])
	}

	return nil
}

// Terminate implements System.
func (*LevelManager) Terminate() {
}

// Update implements System.
func (lm *LevelManager) Update(aa *state.GameState) error {
	return nil
}

func (lm *LevelManager) newRoom(input *resources.Room) {
	cellSize := 16
	lm.Space = resolv.NewSpace[resources.TagType](input.W, input.H, cellSize, cellSize)
	lm.Space.Add(input.Collision...)
	lm.bakeRoomImage(input)
}

func (lm *LevelManager) bakeRoomImage(source *resources.Room) {
	lm.roomRender = ebiten.NewImage(source.W, source.H)
	for _, layer := range source.Layers {
		for _, tileData := range layer {
			id := tileData.T
			x1 := int(id%12) * 16
			y1 := int(id/12) * 16
			rect := image.Rect(x1, y1, x1+16, y1+16)
			// TODO cache tileset subimages for later use
			tileImg := lm.tileSet.SubImage(rect).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tileData.X), float64(tileData.Y))
			lm.roomRender.DrawImage(tileImg, op)
		}
	}
}

func (lm *LevelManager) ChangeRoom(input int) {
	// slog.Info("EXIT EXIT EXIT",
	// 	"EXIT", input,
	// 	"rooms", lm.currentLevel,
	// )
	lm.newRoom(lm.currentLevel.Rooms[input])
}
