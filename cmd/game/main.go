package main

import (
	"2048/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
	boardSize    = 4
)

func main() {
	g, err := game.NewGame(screenWidth, screenHeight, boardSize)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2048")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
