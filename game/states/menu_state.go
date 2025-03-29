package states

import (
	"image/color"
	"math"
	"os"
	"time"

	"github.com/Nathene/bitbase/game"
	"github.com/Nathene/bitbase/game/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MenuState represents the menu state of the game
type MenuState struct {
	background    *ebiten.Image
	buttons       []*ui.Button
	selectedIndex int
	assetManager  *game.AssetManager
	stateManager  *StateManager

	// Camera movement variables
	cameraX         float64
	cameraY         float64
	cameraStartTime time.Time
}

// NewMenuState creates a new menu state
func NewMenuState(assetManager *game.AssetManager, stateManager *StateManager) *MenuState {
	return &MenuState{
		assetManager:    assetManager,
		stateManager:    stateManager,
		selectedIndex:   0,
		cameraStartTime: time.Now(),
		// Initialize camera position with an offset
		cameraX: -1000, // Negative because we move the background in the opposite direction
		cameraY: -1000, // Negative because we move the background in the opposite direction
	}
}

// Initialize sets up the menu state
func (ms *MenuState) Initialize() error {
	// Only try to load the world background since we know it exists
	ms.background = ms.assetManager.GetImage("worldBackground")

	// Create buttons
	buttonWidth := 200.0
	buttonHeight := 50.0
	buttonX := (game.ScreenWidth - buttonWidth) / 2
	buttonStartY := 400.0 // Position buttons lower on the screen
	buttonSpacing := 70.0

	// Create buttons without passing a font
	playButton := ui.NewButton(buttonX, buttonStartY, buttonWidth, buttonHeight, "Play", nil)
	playButton.OnClick = func() {
		// Create and push a new gameplay state when clicked
		gameplayState := NewGameplayState(ms.assetManager, ms.stateManager)
		ms.stateManager.RequestStateChange(StateChange{
			changeType: Replace,
			state:      gameplayState,
		})
	}

	// Make buttons visually distinguishable with clear colors and thicker borders
	playButton.BackgroundColor = color.RGBA{0, 120, 0, 255} // Green
	playButton.HoverColor = color.RGBA{0, 180, 0, 255}      // Brighter green on hover
	playButton.BorderWidth = 3                              // Thicker border for better visibility

	// Options button
	optionsButton := ui.NewButton(buttonX, buttonStartY+buttonSpacing, buttonWidth, buttonHeight, "Options", nil)
	optionsButton.OnClick = func() {
		// For now, we'll just do nothing
	}
	optionsButton.BackgroundColor = color.RGBA{0, 0, 120, 255} // Blue
	optionsButton.HoverColor = color.RGBA{0, 0, 180, 255}      // Brighter blue on hover
	optionsButton.BorderWidth = 3                              // Thicker border

	// Exit button
	exitButton := ui.NewButton(buttonX, buttonStartY+buttonSpacing*2, buttonWidth, buttonHeight, "Exit", nil)
	exitButton.OnClick = func() {
		// Exit the game
		os.Exit(0)
	}
	exitButton.BackgroundColor = color.RGBA{120, 0, 0, 255} // Red
	exitButton.HoverColor = color.RGBA{180, 0, 0, 255}      // Brighter red on hover
	exitButton.BorderWidth = 3                              // Thicker border

	ms.buttons = []*ui.Button{playButton, optionsButton, exitButton}

	return nil
}

// Enter is called when this state becomes active
func (ms *MenuState) Enter() error {
	// Start menu music would go here if implemented
	return nil
}

// Exit is called when this state is no longer active
func (ms *MenuState) Exit() error {
	// Stop menu music would go here if implemented
	return nil
}

// Update handles menu logic and input
func (ms *MenuState) Update() error {
	// Update camera position for panning movement
	elapsedTime := time.Since(ms.cameraStartTime).Seconds()

	// Create a slow, panning movement pattern
	// Total cycle takes about 30 seconds (0.2 radians per second)
	cyclePosition := math.Mod(elapsedTime*1, 2*math.Pi)

	// Use sine and cosine to create a circular panning pattern
	// Scale factor controls how far the camera moves from center
	scaleFactor := 20.0

	// Base camera position (starting at -1000, -1000 to position view at 1000,1000)
	baseX := -1000.0
	baseY := -1000.0

	// Calculate camera position based on where we are in the cycle
	ms.cameraX = baseX + scaleFactor*math.Cos(cyclePosition)
	ms.cameraY = baseY + scaleFactor*math.Sin(cyclePosition)

	// Update all buttons
	for _, button := range ms.buttons {
		button.Update()
	}

	// Keyboard navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		ms.selectedIndex = (ms.selectedIndex + 1) % len(ms.buttons)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		ms.selectedIndex = (ms.selectedIndex - 1 + len(ms.buttons)) % len(ms.buttons)
	}

	// Highlight selected button and remove highlight from others
	for i, button := range ms.buttons {
		if i == ms.selectedIndex {
			button.BorderColor = color.RGBA{255, 255, 100, 255}
			button.BorderWidth = 3
		} else {
			button.BorderColor = color.RGBA{200, 200, 200, 255}
			button.BorderWidth = 2
		}
	}

	// Handle selection with keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if ms.selectedIndex >= 0 && ms.selectedIndex < len(ms.buttons) {
			if ms.buttons[ms.selectedIndex].OnClick != nil {
				ms.buttons[ms.selectedIndex].OnClick()
			}
		}
	}

	return nil
}

// HandleInput processes all input for this state
func (ms *MenuState) HandleInput() error {
	// Input handling is done in Update for now
	return nil
}

// Draw renders the menu
func (ms *MenuState) Draw(screen *ebiten.Image) {
	// Draw background with camera movement
	if ms.background != nil {
		op := &ebiten.DrawImageOptions{}

		// Apply camera transformation - this moves the background with a parallax effect
		op.GeoM.Translate(ms.cameraX, ms.cameraY)

		screen.DrawImage(ms.background, op)
	} else {
		// Fallback background
		screen.Fill(color.RGBA{20, 30, 50, 255})
	}

	// Title box (no text, just a colored box)
	titleBarHeight := 80.0
	titleBarWidth := 400.0
	titleBarX := (game.ScreenWidth - titleBarWidth) / 2
	titleBarY := 100.0

	// Draw title box with a more distinctive color since we can't display text
	vector.DrawFilledRect(
		screen,
		float32(titleBarX),
		float32(titleBarY),
		float32(titleBarWidth),
		float32(titleBarHeight),
		color.RGBA{180, 180, 60, 255}, // Bright yellow/gold for the title box
		false,
	)

	// Draw a border for the title box to make it more distinctive
	vector.StrokeRect(
		screen,
		float32(titleBarX),
		float32(titleBarY),
		float32(titleBarWidth),
		float32(titleBarHeight),
		3.0, // Thicker border
		color.RGBA{200, 200, 200, 255},
		false,
	)

	// Draw all buttons
	for _, button := range ms.buttons {
		button.Draw(screen)
	}
}

// GetStateID returns a unique identifier for this state
func (ms *MenuState) GetStateID() string {
	return "MainMenu"
}
