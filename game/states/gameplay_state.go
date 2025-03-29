package states

import (
	"time"

	"github.com/Nathene/bitbase/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// GameplayState handles the actual gameplay
type GameplayState struct {
	game          *game.Game // Your existing game implementation
	isPaused      bool
	gameStartTime time.Time
	assetManager  *game.AssetManager
	stateManager  *StateManager
}

// NewGameplayState creates a new gameplay state
func NewGameplayState(assetManager *game.AssetManager, stateManager *StateManager) *GameplayState {
	return &GameplayState{
		assetManager:  assetManager,
		stateManager:  stateManager,
		gameStartTime: time.Now(),
	}
}

// Initialize sets up the gameplay state
func (gs *GameplayState) Initialize() error {
	// Load the player sprite sheet and background
	playerSheet := gs.assetManager.GetImage("playerSheet")
	backgroundImage := gs.assetManager.GetImage("worldBackground")

	// Create the game instance
	gs.game = game.NewGame(playerSheet, backgroundImage)

	return nil
}

// Enter is called when this state becomes active
func (gs *GameplayState) Enter() error {
	gs.isPaused = false
	return nil
}

// Exit is called when this state is no longer active
func (gs *GameplayState) Exit() error {
	return nil
}

// Update handles gameplay logic
func (gs *GameplayState) Update() error {
	if gs.isPaused {
		return nil
	}

	// Check for pause action
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		gs.isPaused = true

		// Create and push the pause state
		pauseState := NewPauseState(gs.assetManager, gs.stateManager, gs)
		gs.stateManager.PushState(pauseState)

		return nil
	}

	// Update the game
	return gs.game.Update()
}

// HandleInput processes all input for this state
func (gs *GameplayState) HandleInput() error {
	// Input handling is done in Update and by the game instance
	return nil
}

// Draw renders the gameplay
func (gs *GameplayState) Draw(screen *ebiten.Image) {
	// Let the game draw itself
	gs.game.Draw(screen)
}

// GetStateID returns a unique identifier for this state
func (gs *GameplayState) GetStateID() string {
	return "Gameplay"
}
