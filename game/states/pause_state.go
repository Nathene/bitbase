package states

import (
	"image/color"
	"os"

	"github.com/Nathene/bitbase/game"
	"github.com/Nathene/bitbase/game/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// PauseState handles the pause menu that appears when the game is paused
type PauseState struct {
	background    *ebiten.Image
	overlay       *ebiten.Image
	buttons       []*ui.Button
	selectedIndex int
	assetManager  *game.AssetManager
	stateManager  *StateManager
	previousState GameState
}

// NewPauseState creates a new pause state
func NewPauseState(assetManager *game.AssetManager, stateManager *StateManager, previousState GameState) *PauseState {
	return &PauseState{
		assetManager:  assetManager,
		stateManager:  stateManager,
		selectedIndex: 0,
		previousState: previousState,
	}
}

// Initialize sets up the pause state
func (ps *PauseState) Initialize() error {
	// Create a semi-transparent overlay
	ps.overlay = ebiten.NewImage(game.ScreenWidth, game.ScreenHeight)
	ps.overlay.Fill(color.RGBA{0, 0, 0, 180}) // Semi-transparent black

	// Create buttons
	buttonWidth := 200.0
	buttonHeight := 50.0
	buttonX := (game.ScreenWidth - buttonWidth) / 2
	buttonStartY := float64(game.ScreenHeight)/2 - 50 // Centered vertically
	buttonSpacing := 70.0

	// Resume button
	resumeButton := ui.NewButton(buttonX, buttonStartY, buttonWidth, buttonHeight, "Resume", nil)
	resumeButton.OnClick = func() {
		// Pop this state to return to the game
		ps.stateManager.PopState()
	}

	// Make buttons visually distinguishable
	resumeButton.BackgroundColor = color.RGBA{0, 120, 0, 255} // Green
	resumeButton.HoverColor = color.RGBA{0, 180, 0, 255}      // Brighter green on hover
	resumeButton.BorderWidth = 3                              // Thicker border

	// Menu button
	menuButton := ui.NewButton(buttonX, buttonStartY+buttonSpacing, buttonWidth, buttonHeight, "Main Menu", nil)
	menuButton.OnClick = func() {
		// Return to main menu
		menuState := NewMenuState(ps.assetManager, ps.stateManager)
		ps.stateManager.RequestStateChange(StateChange{
			changeType: Replace, // Replace both the game and pause states
			state:      menuState,
		})
	}
	menuButton.BackgroundColor = color.RGBA{60, 60, 180, 255} // Blue
	menuButton.HoverColor = color.RGBA{80, 80, 255, 255}      // Brighter blue on hover
	menuButton.BorderWidth = 3                                // Thicker border

	// Exit button
	exitButton := ui.NewButton(buttonX, buttonStartY+buttonSpacing*2, buttonWidth, buttonHeight, "Exit Game", nil)
	exitButton.OnClick = func() {
		// Exit the game
		os.Exit(0)
	}
	exitButton.BackgroundColor = color.RGBA{180, 0, 0, 255} // Red
	exitButton.HoverColor = color.RGBA{255, 0, 0, 255}      // Brighter red on hover
	exitButton.BorderWidth = 3                              // Thicker border

	ps.buttons = []*ui.Button{resumeButton, menuButton, exitButton}

	return nil
}

// Enter is called when this state becomes active
func (ps *PauseState) Enter() error {
	return nil
}

// Exit is called when this state is no longer active
func (ps *PauseState) Exit() error {
	return nil
}

// Update handles pause menu logic and input
func (ps *PauseState) Update() error {
	// Check for Escape key to resume
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ps.stateManager.PopState()
		return nil
	}

	// Update all buttons
	for _, button := range ps.buttons {
		button.Update()
	}

	// Keyboard navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ps.selectedIndex = (ps.selectedIndex + 1) % len(ps.buttons)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ps.selectedIndex = (ps.selectedIndex - 1 + len(ps.buttons)) % len(ps.buttons)
	}

	// Highlight selected button and remove highlight from others
	for i, button := range ps.buttons {
		if i == ps.selectedIndex {
			button.BorderColor = color.RGBA{255, 255, 100, 255}
			button.BorderWidth = 4 // Extra thick for selected
		} else {
			button.BorderColor = color.RGBA{200, 200, 200, 255}
			button.BorderWidth = 3
		}
	}

	// Handle selection with keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if ps.selectedIndex >= 0 && ps.selectedIndex < len(ps.buttons) {
			if ps.buttons[ps.selectedIndex].OnClick != nil {
				ps.buttons[ps.selectedIndex].OnClick()
			}
		}
	}

	return nil
}

// HandleInput processes all input for this state
func (ps *PauseState) HandleInput() error {
	// Input handling is done in Update
	return nil
}

// Draw renders the pause menu
func (ps *PauseState) Draw(screen *ebiten.Image) {
	// First draw the previous state (the game) if it implements Draw
	if drawer, ok := ps.previousState.(GameState); ok {
		drawer.Draw(screen)
	}

	// Draw semi-transparent overlay
	screen.DrawImage(ps.overlay, nil)

	// Pause title (no text, just a colored box)
	titleBarHeight := 60.0
	titleBarWidth := 300.0
	titleBarX := (game.ScreenWidth - titleBarWidth) / 2
	titleBarY := float64(game.ScreenHeight)/2 - 150

	// Draw a distinctive title box since we can't display text
	vector.DrawFilledRect(
		screen,
		float32(titleBarX),
		float32(titleBarY),
		float32(titleBarWidth),
		float32(titleBarHeight),
		color.RGBA{180, 60, 60, 200}, // Red tone for pause indicator
		false,
	)

	// Add border to make the title box more distinctive
	vector.StrokeRect(
		screen,
		float32(titleBarX),
		float32(titleBarY),
		float32(titleBarWidth),
		float32(titleBarHeight),
		3.0, // Thicker border
		color.RGBA{220, 220, 220, 255},
		false,
	)

	// Draw all buttons
	for _, button := range ps.buttons {
		button.Draw(screen)
	}
}

// GetStateID returns a unique identifier for this state
func (ps *PauseState) GetStateID() string {
	return "PauseMenu"
}
