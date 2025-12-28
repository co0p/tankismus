package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/co0p/tankismus/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowTitle("tankismus")
	ebiten.SetWindowSize(800, 600)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
