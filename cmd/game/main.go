package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/justgook/goplatformer/pkg/game"
	"github.com/justgook/goplatformer/pkg/gameLogger/cli"
	"log/slog"
	"os"
)

func main() {
	handler := cli.New(os.Stderr, &cli.Options{
		HandlerOptions: slog.HandlerOptions{Level: slog.LevelDebug},
	})

	slog.SetDefault(slog.New(handler))

	if err := ebiten.RunGame(game.New()); err != nil {
		fmt.Print(err)
		return
	}
}
