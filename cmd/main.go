package main

import (
	"log"

	"github.com/Nathene/bitbase/game"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	playerSheet, _, err := ebitenutil.NewImageFromFile("assets/character/base_idle_strip9.png")
	if err != nil {
		log.Fatal("Failed to load player sprite sheet")
	}

	backgroundImage, _, err := ebitenutil.NewImageFromFile("assets/world/Sunnyside_World_ExampleScene.png")
	if err != nil {
		log.Fatal("Failed to load world image")
	}

	g := game.NewGame(playerSheet, backgroundImage)
	g.Run()
}
