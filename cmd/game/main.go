package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
)

func main() {
	handler := cli.New(os.Stdout, &cli.Options{
		HandlerOptions: slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	})
	logger := slog.New(handler)

	slog.SetDefault(logger)

	// 640x360 * 2
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Game title goes here")

	if err := ebiten.RunGame(game.New()); err != nil {
		fmt.Print(err)
		return
	}
}
