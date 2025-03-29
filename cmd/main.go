package main

import (
	"log"

	"github.com/Nathene/bitbase/game"
	"github.com/Nathene/bitbase/game/states"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameApplication struct {
	stateManager *states.StateManager
	assetManager *game.AssetManager
}

// Update updates the game state
func (g *GameApplication) Update() error {
	return g.stateManager.Update()
}

// Draw renders the game
func (g *GameApplication) Draw(screen *ebiten.Image) {
	g.stateManager.Draw(screen)
}

// Layout returns the game's logical screen dimensions
func (g *GameApplication) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.ScreenWidth, game.ScreenHeight
}

func main() {
	// Create asset manager
	assetManager := game.NewAssetManager()

	// Create state manager
	stateManager := states.NewStateManager()

	// Create application
	app := &GameApplication{
		stateManager: stateManager,
		assetManager: assetManager,
	}

	// Create menu state
	menuState := states.NewMenuState(assetManager, stateManager)

	// Create loading state as the initial state with menu as the next state
	loadingState := states.NewLoadingState(assetManager, stateManager, menuState)

	// Push the loading state as the first state
	stateManager.PushState(loadingState)

	// Set up window
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("BitBase Game")

	// Run the game
	if err := ebiten.RunGame(app); err != nil {
		log.Fatal(err)
	}
}
