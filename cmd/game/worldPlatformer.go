package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/justgook/goplatformer"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/justgook/goplatformer/pkg/aseprite"
	"github.com/justgook/goplatformer/pkg/ldtk"
	"github.com/justgook/goplatformer/pkg/resolv"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type Player struct {
	Object         *resolv.Object
	SpeedX         float64
	SpeedY         float64
	OnGround       *resolv.Object
	WallSliding    *resolv.Object
	FacingRight    bool
	IgnorePlatform *resolv.Object
	// New Stuff
	Animation *AnimSprite
}

func NewPlayer(space *resolv.Space) *Player {
	p := &Player{
		Object:      resolv.NewObject(32, 128, 16, 24),
		FacingRight: true,
	}
	p.Object.SetShape(resolv.NewRectangle(0, 0, p.Object.W, p.Object.H))

	space.Add(p.Object)
	p.Animation = initPlayerAnimation()

	return p
}

type WorldPlatformer struct {
	Game   *Game
	Space  *resolv.Space
	Player *Player

	FloatingPlatform      *resolv.Object
	FloatingPlatformTween *gween.Sequence

	TileMapRenderer *EbitenRenderer
}

func NewWorldPlatformer(game *Game) *WorldPlatformer {
	w := &WorldPlatformer{Game: game}
	w.Init()

	return w
}

type AnimSprite struct {
	Animation    string
	Source       *ebiten.Image
	Sprite       *ebiten.Image
	currentFrame int
	Data         AnimDataMap
}

type AnimDataMap = map[string][]FFrame

type FFrame struct {
	Duration int
	Layers   []FrameDrawData
}
type FrameDrawData struct {
	W  int
	H  int
	tX int
	tY int
	X0 int
	Y0 int
	X1 int
	Y1 int
}

// TODO Move this to build step and us gob file for that
// setup aseprite to:
// Item File name `{layer} {frame}`
// Item Tag Name `{tag}`
func convertAsprite2My(input *aseprite.Animation) map[string][]FFrame {
	output := make(AnimDataMap, len(input.Meta.FrameTags))
	layerCount := len(input.Meta.Layers)
	frameMaps := input.Frames
	for _, anim := range input.Meta.FrameTags {
		output[anim.Name] = make([]FFrame, 0, anim.To-anim.From+1)
		for i := anim.From; i <= anim.To; i++ {
			item := FFrame{Layers: make([]FrameDrawData, 0, layerCount)}
			for _, layer := range input.Meta.Layers {
				key := fmt.Sprintf("%s %d", layer.Name, i)
				data := FrameDrawData{
					W:  frameMaps[key].Frame.W,
					H:  frameMaps[key].Frame.H,
					X0: frameMaps[key].Frame.X,
					Y0: frameMaps[key].Frame.Y,
					X1: frameMaps[key].Frame.X + frameMaps[key].Frame.W,
					Y1: frameMaps[key].Frame.Y + frameMaps[key].Frame.H,
					tX: frameMaps[key].SpriteSourceSize.X,
					tY: frameMaps[key].SpriteSourceSize.Y,
				}
				item.Duration = int(frameMaps[key].Duration)
				item.Layers = append(item.Layers, data)
			}
			output[anim.Name] = append(output[anim.Name], item)
		}
	}

	return output
}

func initPlayerAnimation() *AnimSprite {
	playerAnim, err := aseprite.UnmarshalAnimation(goplatformer.AnimationJson)
	if err != nil {
		panic(err)
	}
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(goplatformer.AnimationPNG))
	if err != nil {
		log.Fatal(err)
	}

	return &AnimSprite{
		Animation: "Run",
		Source:    ebiten.NewImageFromImage(img),
		//TODO Update to rel size
		Sprite: ebiten.NewImage(48, 48),
		Data:   convertAsprite2My(&playerAnim),
	}
}

func (a *AnimSprite) Update() {
	animationFrame := (a.currentFrame / 5) % (len(a.Data[a.Animation]))
	currentFrameData := a.Data[a.Animation][animationFrame]
	a.Sprite.Fill(color.RGBA{})
	for _, info := range currentFrameData.Layers {
		rect := image.Rect(info.X0, info.Y0, info.X1, info.Y1)
		result := a.Source.SubImage(rect).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(info.tX), float64(info.tY))
		a.Sprite.DrawImage(result, op)
	}
	a.currentFrame++
}

type embededImageLoader struct {
	tilesets map[string]*ebiten.Image
}

func NewTilesetLoader() *embededImageLoader {
	img1, _, err := image.Decode(bytes.NewReader(goplatformer.TilesetPng1))
	if err != nil {
		log.Fatal(err)
	}
	img2, _, err := image.Decode(bytes.NewReader(goplatformer.TilesetPng2))
	if err != nil {
		log.Fatal(err)
	}
	return &embededImageLoader{
		tilesets: map[string]*ebiten.Image{
			"tileset.png":  ebiten.NewImageFromImage(img1),
			"tileset2.png": ebiten.NewImageFromImage(img2),
		},
	}

}

func (d *embededImageLoader) LoadTileset(tilesetPath string) *ebiten.Image {
	return d.tilesets[tilesetPath]
}

func (world *WorldPlatformer) initLevel() {
	ldtkProject, err := ldtk.Read(goplatformer.ExampleLevel)
	if err != nil {
		panic(err)
	}

	// Choose a level...
	level := ldtkProject.Levels[0]

	// Create a new renderer...
	// EbitenRenderer.DiskLoader loads images from disk using ebitenutil.NewImageFromFile() and takes an argument of the base path to use when loading.
	// We pass a blank string to NewDiskLoader() because for the example, the assets are in the same directory.
	world.TileMapRenderer = NewEbitenRenderer(NewTilesetLoader())

	// ... And render the tiles for the level out to layers, which will be *ebiten.Images. We'll retrieve them to draw in a Draw() loop later.
	world.TileMapRenderer.Render(level)
}

func (world *WorldPlatformer) Init() {
	world.initLevel()
	// Initialize the world.
	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	// Define the world's Space. Here, a Space is essentially a grid (the game's width and height, or 640x360), made up of 16x16 cells. Each cell can have 0 or more Objects within it,
	// and collisions can be found by checking the Space to see if the Cells at specific positions contain (or would contain) Objects. This is a broad, simplified approach to collision
	// detection.
	world.Space = resolv.NewSpace(int(gw), int(gh), 16, 16)

	// Construct the solid level geometry.
	// Note that the simple approach of checking cells in a Space for collision works simply when the geometry is
	// aligned with the cells, as it all is in this platformer example.
	world.Space.Add(
		resolv.NewObject(0, 0, 16, gh, "solid"),
		resolv.NewObject(gw-16, 0, 16, gh, "solid"),
		resolv.NewObject(0, 0, gw, 16, "solid"),
		resolv.NewObject(0, gh-24, gw, 32, "solid"),
		resolv.NewObject(160, gh-56, 160, 32, "solid"),
		resolv.NewObject(320, 64, 32, 160, "solid"),
		resolv.NewObject(64, 128, 16, 160, "solid"),
		resolv.NewObject(gw-128, 64, 128, 16, "solid"),
		resolv.NewObject(gw-128, gh-88, 128, 16, "solid"),
	)

	// Create the Player. NewPlayer adds it to the world's Space.
	world.Player = NewPlayer(world.Space)

	// Create the floating platform.
	world.FloatingPlatform = resolv.NewObject(128, gh-32, 128, 8)
	world.FloatingPlatform.AddTags("platform")

	// The floating platform moves using a *gween.Sequence sequence of tweens, moving it back and forth.
	world.FloatingPlatformTween = gween.NewSequence()
	world.FloatingPlatformTween.Add(
		gween.New(float32(world.FloatingPlatform.Y)-48, float32(world.FloatingPlatform.Y-128), 2, ease.Linear),
		gween.New(float32(world.FloatingPlatform.Y-128), float32(world.FloatingPlatform.Y)-48, 2, ease.Linear),
	)
	world.Space.Add(world.FloatingPlatform)

	// Non-moving floating Platforms.
	world.Space.Add(
		resolv.NewObject(352, 64, 48, 8, "platform"),
		resolv.NewObject(352, 64+64, 48, 8, "platform"),
		resolv.NewObject(352, 64+128, 48, 8, "platform"),
		resolv.NewObject(352, 64+192, 48, 8, "platform"),
	)

	// A ramp, which is unique as it has a non-rectangular shape. For this, we will specify a different shape for collision testing.
	ramp := resolv.NewObject(320, gh-56, 64, 32, "ramp")

	// We will construct the shape using a ConvexPolygon. It's essentially an elogated triangle, but with a "floor" afterwards,
	// ensuring the Player is always able to stand regardless of which ramp they're standing on.

	rampShape := resolv.NewConvexPolygon(
		0, 0,
		// Vertices:
		0, 0,
		2, 0, // The extra 2 pixels here make it so the Player doesn't get stuck for a frame or two when running up the ramp.
		ramp.W-2, ramp.H, // Same here; an extra 2 pixels makes it so that dismounting the ramp is nice and easy
		ramp.W, ramp.H,
		0, ramp.H,
	)

	world.Space.Add(ramp)
	ramp.SetShape(rampShape)
}

func (world *WorldPlatformer) Update() {
	// Floating platform movement needs to be done before the player's movement update to make sure there's no space between its top and the player's bottom;
	// otherwise, an alternative might be to have the platform detect to see if the Player's resting on it, and if so, move the player up manually.
	y, _, seqDone := world.FloatingPlatformTween.Update(1.0 / 60.0)
	world.FloatingPlatform.Y = float64(y)
	if seqDone {
		world.FloatingPlatformTween.Reset()
	}
	world.FloatingPlatform.Update()

	// Now we update the Player's movement. This is the real bread-an-butter of this example, naturally.
	player := world.Player

	friction := 0.5
	accel := 0.5 + friction
	maxSpeed := 4.0
	jumpSpd := 10.0
	gravity := 0.75

	player.SpeedY += gravity

	if player.WallSliding != nil && player.SpeedY > 1 {
		player.SpeedY = 1
	}

	// Horizontal movement is only possible when not wall sliding.
	if player.WallSliding == nil {
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxisValue(0, 0) > 0.1 {
			player.SpeedX += accel
			player.FacingRight = true
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxisValue(0, 0) < -0.1 {
			player.SpeedX -= accel
			player.FacingRight = false
		}
	}

	// Apply friction and horizontal speed limiting.
	if player.SpeedX > friction {
		player.SpeedX -= friction
	} else if player.SpeedX < -friction {
		player.SpeedX += friction
	} else {
		player.SpeedX = 0
	}

	if player.SpeedX > maxSpeed {
		player.SpeedX = maxSpeed
	} else if player.SpeedX < -maxSpeed {
		player.SpeedX = -maxSpeed
	}

	// Check for jumping.
	jumpKeyJustPressed := inpututil.IsKeyJustPressed(ebiten.KeyX) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
	if jumpKeyJustPressed || ebiten.IsGamepadButtonPressed(0, 0) || ebiten.IsGamepadButtonPressed(1, 0) {
		if (ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.GamepadAxisValue(0, 1) > 0.1 || ebiten.GamepadAxisValue(1, 1) > 0.1) && player.OnGround != nil && player.OnGround.HasTags("platform") {
			player.IgnorePlatform = player.OnGround
		} else {
			if player.OnGround != nil {
				player.SpeedY = -jumpSpd
			} else if player.WallSliding != nil {
				// WALL JUMPING
				player.SpeedY = -jumpSpd

				if player.WallSliding.X > player.Object.X {
					player.SpeedX = -4
				} else {
					player.SpeedX = 4
				}

				player.WallSliding = nil
			}
		}
	}

	// We handle horizontal movement separately from vertical movement. This is, conceptually, decomposing movement into two phases / axes.
	// By decomposing movement in this manner, we can handle each case properly (i.e. stop movement horizontally separately from vertical movement, as
	// necessary). More can be seen on this topic over on this blog post on:
	// http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

	// dx is the horizontal delta movement variable (which is the Player's horizontal speed). If we come into contact with something, then it will
	// be that movement instead.
	dx := player.SpeedX

	// Moving horizontally is done fairly simply; we just check to see if something solid is in front of us. If so, we move into contact with it
	// and stop horizontal movement speed. If not, then we can just move forward.

	if check := player.Object.Check(player.SpeedX, 0, "solid"); check != nil {

		dx = check.ContactWithCell(check.Cells[0]).X()
		player.SpeedX = 0

		// If you're in the air, then colliding with a wall object makes you start wall sliding.
		if player.OnGround == nil {
			player.WallSliding = check.Objects[0]
		}

	}

	// Then we just apply the horizontal movement to the Player's Object. Easy-peasy.
	player.Object.X += dx

	// Now for the vertical movement; it's the most complicated because we can land on different types of objects and need
	// to treat them all differently, but overall, it's not bad.

	// First, we set OnGround to be nil, in case we don't end up standing on anything.
	player.OnGround = nil

	// dy is the delta movement downward, and is the vertical movement by default; similarly to dx, if we come into contact with
	// something, this will be changed to move to contact instead.

	dy := player.SpeedY

	// We want to be sure to lock vertical movement to a maximum of the size of the Cells within the Space
	// so we don't miss any collisions by tunneling through.

	dy = math.Max(math.Min(dy, 16), -16)

	// We're going to check for collision using dy (which is vertical movement speed), but add one when moving downwards to look a bit deeper down
	// into the ground for solid objects to land on, specifically.
	checkDistance := dy
	if dy >= 0 {
		checkDistance++
	}

	// We check for any solid / stand-able objects. In actuality, there aren't any other Objects
	// with other tags in this Space, so we don't -have- to specify any tags, but it's good to be specific for clarity in this example.
	if check := player.Object.Check(0, checkDistance, "solid", "platform", "ramp"); check != nil {
		// So! Firstly, we want to see if we jumped up into something that we can slide around horizontally to avoid bumping the Player's head.

		// Sliding around a misspaced jump is a small thing that makes jumping a bit more forgiving, and is something different polished platformers
		// (like the 2D Mario games) do to make it a smidge more comfortable to play. For a visual example of this, see this excellent devlog post
		// from the extremely impressive indie game, Leilani's Island: https://forums.tigsource.com/index.php?topic=46289.msg1387138#msg1387138

		// To accomplish this sliding, we simply call Collision.SlideAgainstCell() to see if we can slide.
		// We pass the first cell, and tags that we want to avoid when sliding (i.e. we don't want to slide into cells that contain other solid objects).

		slide := check.SlideAgainstCell(check.Cells[0], "solid")

		// We further ensure that we only slide if:
		// 1) We're jumping up into something (dy < 0),
		// 2) If the cell we're bumping up against contains a solid object,
		// 3) If there was, indeed, a valid slide left or right, and
		// 4) If the proposed slide is less than 8 pixels in horizontal distance. (This is a relatively arbitrary number that just so happens to be half the
		// width of a cell. This is to ensure the player doesn't slide too far horizontally.)

		if dy < 0 && check.Cells[0].ContainsTags("solid") && slide != nil && math.Abs(slide.X()) <= 8 {

			// If we are able to slide here, we do so. No contact was made, and vertical speed (dy) is maintained upwards.
			player.Object.X += slide.X()

		} else {

			// If sliding -fails-, that means the Player is jumping directly onto or into something, and we need to do more to see if we need to come into
			// contact with it. Let's press on!

			// First, we check for ramps. For ramps, we can't simply check for collision with Check(), as that's not precise enough. We need to get a bit
			// more information, and so will do so by checking its Shape (a triangular ConvexPolygon, as defined in WorldPlatformer.Init()) against the
			// Player's Shape (which is also a rectangular ConvexPolygon).

			// We get the ramp by simply filtering out Objects with the "ramp" tag out of the objects returned in our broad Check(), and grabbing the first one
			// if there's any at all.
			if ramps := check.ObjectsByTags("ramp"); len(ramps) > 0 {

				ramp := ramps[0]

				// For simplicity, this code assumes we can only stand on one ramp at a time as there is only one ramp in this example.
				// In actuality, if there was a possibility to have a potential collision with multiple ramps (i.e. a ramp that sits on another ramp, and the player running down
				// one onto the other), the collision testing code should probably go with the ramp with the highest confirmed intersection point out of the two.

				// Next, we see if there's been an intersection between the two Shapes using Shape.Intersection. We pass the ramp's shape, and also the movement
				// we're trying to make horizontally, as this makes Intersection return the next y-position while moving, not the one directly
				// underneath the Player. This would keep the player from getting "stuck" when walking up a ramp into the top of a solid block, if there weren't
				// a landing at the top and bottom of the ramp.

				// We use 8 here for the Y-delta so that we can easily see if you're running down the ramp (in which case you're probably in the air as you
				// move faster than you can fall in this example). This way we can maintain contact so you can always jump while running down a ramp. We only
				// continue with coming into contact with the ramp as long as you're not moving upwards (i.e. jumping).

				if contactSet := player.Object.Shape.Intersection(dx, 8, ramp.Shape); dy >= 0 && contactSet != nil {
					// If Intersection() is successful, a ContactSet is returned. A ContactSet contains information regarding where
					// two Shapes intersect, like the individual points of contact, the center of the contacts, and the MTV, or
					// Minimum Translation Vector, to move out of contact.

					// Here, we use ContactSet.TopmostPoint() to get the top-most contact point as an indicator of where
					// we want the player's feet to be. Then we just set that position, and we're done.

					dy = contactSet.TopmostPoint()[1] - player.Object.Bottom() + 0.1
					player.OnGround = ramp
					player.SpeedY = 0

				}
			}

			// Platforms are next; here, we just see if the platform is not being ignored by attempting to drop down,
			// if the player is falling on the platform (as otherwise he would be jumping through platforms), and if the platform is low enough
			// to land on. If so, we stand on it.

			// Because there's a moving floating platform, we use Collision.ContactWithObject() to ensure the player comes into contact
			// with the top of the platform object. An alternative would be to use Collision.ContactWithCell(), but that would be only if the
			// platform didn't move and were aligned with the Spatial cellular grid.

			if platforms := check.ObjectsByTags("platform"); len(platforms) > 0 {
				platform := platforms[0]

				if platform != player.IgnorePlatform && player.SpeedY >= 0 && player.Object.Bottom() < platform.Y+4 {
					dy = check.ContactWithObject(platform).Y()
					player.OnGround = platform
					player.SpeedY = 0
				}

			}

			// Finally, we check for simple solid ground. If we haven't had any success in landing previously, or the solid ground
			// is higher than the existing ground (like if the platform passes underneath the ground, or we're walking off of solid ground
			// onto a ramp), we stand on it instead. We don't check for solid collision first because we want any ramps to override solid
			// ground (so that you can walk onto the ramp, rather than sticking to solid ground).

			// We use ContactWithObject() here because otherwise, we might come into contact with the moving platform's cells (which, naturally,
			// would be selected by a Collision.ContactWithCell() call because the cell is closest to the Player).

			if solids := check.ObjectsByTags("solid"); len(solids) > 0 && (player.OnGround == nil || player.OnGround.Y >= solids[0].Y) {
				dy = check.ContactWithObject(solids[0]).Y()
				player.SpeedY = 0

				// We're only on the ground if we land on it (if the object's Y is greater than the player's).
				if solids[0].Y > player.Object.Y {
					player.OnGround = solids[0]
				}

			}

			if player.OnGround != nil {
				player.WallSliding = nil    // Player's on the ground, so no wallsliding anymore.
				player.IgnorePlatform = nil // Player's on the ground, so reset which platform is being ignored.
			}

		}

	}

	// Move the object on dy.
	player.Object.Y += dy

	wallNext := 1.0
	if !player.FacingRight {
		wallNext = -1
	}

	// If the wall next to the Player runs out, stop wall sliding.
	if c := player.Object.Check(wallNext, 0, "solid"); player.WallSliding != nil && c == nil {
		player.WallSliding = nil
	}

	player.Object.Update() // Update the player's position in the space.
	/****************************************************************************************/
	/****************************player animation************************************************************/
	if player.WallSliding != nil {
		player.Animation.Animation = "wSlideLow"
	} else if player.OnGround == nil {
		if player.SpeedY > 2 {
			player.Animation.Animation = "Fall"
		} else if player.SpeedY < -2 {
			player.Animation.Animation = "JumpUp"
		} else {
			player.Animation.Animation = "JumpMax"
		}
	} else if player.SpeedX != 0 {
		player.Animation.Animation = "Run"
	} else {
		player.Animation.Animation = "Idle"

	}
	player.Animation.Update() // Update player animation frame

	/****************************************************************************************/
	// And that's it!
}

func (world *WorldPlatformer) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{G: 145, B: 255, A: 255})

	/*LDtk Starts */
	/* ==================================================================== */
	for _, layer := range world.TileMapRenderer.RenderedLayers {
		if layer.Layer.Identifier == "Indoor" {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}
	/* ==================================================================== */
	/*LDtk Starts */

	drawColor := color.RGBA{R: 255, G: 50, B: 100, A: 255}
	for _, o := range world.Space.Objects() {
		if o.HasTags("ramp") {
			tri := o.Shape.(*resolv.ConvexPolygon)
			world.DrawPolygon(screen, tri, drawColor)
		}
		//if o.HasTags("platform") && o != world.FloatingPlatform {
		//	drawColor := color.RGBA{R: 180, G: 100, A: 255}
		//	vector.DrawFilledRect(screen, float32(o.X), float32(o.Y), float32(o.W), float32(o.H), drawColor, false)
		//} else if o.HasTags("ramp") {
		//	drawColor := color.RGBA{R: 255, G: 50, B: 100, A: 255}
		//	tri := o.Shape.(*resolv.ConvexPolygon)
		//	world.DrawPolygon(screen, tri, drawColor)
		//} else {
		//	drawColor := color.RGBA{R: 60, G: 60, B: 60, A: 255}
		//	vector.DrawFilledRect(screen, float32(o.X), float32(o.Y), float32(o.W), float32(o.H), drawColor, false)
		//}
	}

	// We draw the floating platform separately because Space.Objects() returns Objects in order of which cells they are in, which means
	// that the platform would draw under the solid blocks if it's below it. This way, it always draws on top.
	o := world.FloatingPlatform
	drawColor = color.RGBA{R: 180, G: 100, A: 255}
	vector.DrawFilledRect(screen, float32(o.X), float32(o.Y), float32(o.W), float32(o.H), drawColor, false)

	player := world.Player.Object
	playerColor := color.RGBA{G: 255, B: 60, A: 255}
	if world.Player.OnGround == nil {
		// We draw the player as a different color when jumping so we can visually see when he's in the air.
		playerColor = color.RGBA{R: 200, B: 200, A: 255}
	}
	_ = playerColor
	// vector.DrawFilledRect(screen, float32(player.X), float32(player.Y), float32(player.W), float32(player.H), playerColor, false)

	/* ===================================================== */
	/* Player Sprite*/
	op := &ebiten.DrawImageOptions{}
	if !world.Player.FacingRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(48, 0)
	}
	op.GeoM.Translate(float64(player.X)-16, float64(player.Y)-16)
	screen.DrawImage(world.Player.Animation.Sprite, op)
	/* ===================================================== */

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	/*LDtk Starts */
	/* ==================================================================== */
	for _, layer := range world.TileMapRenderer.RenderedLayers {
		if layer.Layer.Identifier != "Indoor" {
			screen.DrawImage(layer.Image, &ebiten.DrawImageOptions{})
		}
	}
	/* ==================================================================== */
	/*LDtk Starts */

	if world.Game.ShowHelpText {
		world.Game.DrawText(screen, 16, 16,
			"~ Platformer Demo ~",
			"Move Player: Left, Right Arrow",
			"Jump: X Key",
			"Wall-slide: Move into wall in air",
			"Wall-jump: Jump while wall sliding",
			"Fall through platforms: Down + X",
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"F4: Toggle fullscreen",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.ActualFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.ActualTPS())),
		)
	}
}

func (world *WorldPlatformer) DrawPolygon(screen *ebiten.Image, polygon *resolv.ConvexPolygon, color color.Color) {
	for _, line := range polygon.Lines() {
		vector.StrokeLine(screen, float32(line.Start.X()), float32(line.Start.Y()), float32(line.End.X()), float32(line.End.Y()), 1, color, false)
	}
}
