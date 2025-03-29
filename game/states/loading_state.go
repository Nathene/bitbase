package states

import (
	"image/color"
	"time"

	"github.com/Nathene/bitbase/game"
	"github.com/Nathene/bitbase/game/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

// LoadingState represents the loading state of the game
type LoadingState struct {
	assetManager *game.AssetManager
	loadingBar   *ui.ProgressBar
	progress     float64
	isComplete   bool
	nextState    GameState
	stateManager *StateManager

	// Background panning
	backgroundImage *ebiten.Image
	panOffset       float64
	panSpeed        float64

	// Logo
	logoImage *ebiten.Image
}

// NewLoadingState creates a new loading state
func NewLoadingState(assetManager *game.AssetManager, stateManager *StateManager, nextState GameState) *LoadingState {
	return &LoadingState{
		assetManager: assetManager,
		nextState:    nextState,
		stateManager: stateManager,
		panSpeed:     0.5, // Speed of panning in pixels per frame
	}
}

// Initialize sets up the loading state
func (ls *LoadingState) Initialize() error {
	// Create a progress bar
	barWidth := 400.0
	barHeight := 30.0
	ls.loadingBar = ui.NewProgressBar(
		(game.ScreenWidth-barWidth)/2,
		(game.ScreenHeight-barHeight)/2,
		barWidth,
		barHeight,
	)

	// Start the loading process with a completion callback
	ls.assetManager.StartLoading(func() {
		// This is called when all assets have finished loading
		// Wait a moment to show the completed progress bar
		go func() {
			time.Sleep(500 * time.Millisecond)

			ls.stateManager.RequestStateChange(StateChange{
				changeType: Replace,
				state:      ls.nextState,
			})
		}()
	})

	// Queue up only the assets we know exist
	// Load the logo, world background and player character
	ls.assetManager.LoadImage("logo", "assets/loading_screen/logo.png")
	ls.assetManager.LoadImage("worldBackground", "assets/world/example.png")
	ls.assetManager.LoadImage("playerSheet", "assets/character/base_idle_strip9.png")

	// Use worldBackground for the panning effect and load logo
	// Load them after a short delay to ensure the assets are loaded
	go func() {
		time.Sleep(50 * time.Millisecond) // Small delay to wait for the images to load
		ls.backgroundImage = ls.assetManager.GetImage("worldBackground")
	}()

	return nil
}

// Enter is called when this state becomes active
func (ls *LoadingState) Enter() error {
	// Reset progress tracking
	ls.progress = 0
	ls.isComplete = false
	ls.panOffset = 0
	return nil
}

// Exit is called when this state is no longer active
func (ls *LoadingState) Exit() error {
	return nil
}

// Update handles loading progress
func (ls *LoadingState) Update() error {
	// Update progress
	ls.progress = ls.assetManager.GetLoadingProgress()
	ls.loadingBar.Progress = ls.progress
	ls.loadingBar.Update()

	// Update panning animation
	ls.panOffset += ls.panSpeed

	return nil
}

// Draw renders the loading screen
func (ls *LoadingState) Draw(screen *ebiten.Image) {
	// Draw a black background first
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// Draw panning background if available
	if ls.backgroundImage != nil {
		// Get background dimensions
		bgWidth, bgHeight := ls.backgroundImage.Bounds().Dx(), ls.backgroundImage.Bounds().Dy()

		// Draw background with offset
		op := &ebiten.DrawImageOptions{}

		// Scale if necessary to make sure it's at least as tall as the screen
		scaleY := float64(game.ScreenHeight) / float64(bgHeight)
		if scaleY > 1.0 {
			op.GeoM.Scale(scaleY, scaleY)
			bgWidth = int(float64(bgWidth) * scaleY)
		}

		// Apply panning offset (wrap around if needed)
		offset := int(ls.panOffset) % bgWidth

		// Draw first part
		op.GeoM.Translate(-float64(offset), 0)
		screen.DrawImage(ls.backgroundImage, op)

		// Draw second part (wrapped) if first part doesn't fill screen
		if offset > 0 {
			op2 := &ebiten.DrawImageOptions{}
			if scaleY > 1.0 {
				op2.GeoM.Scale(scaleY, scaleY)
			}
			op2.GeoM.Translate(float64(bgWidth-offset), 0)
			screen.DrawImage(ls.backgroundImage, op2)
		}

		// Add a semi-transparent overlay to darken the background
		darkOverlay := ebiten.NewImage(game.ScreenWidth, game.ScreenHeight)
		darkOverlay.Fill(color.RGBA{0, 0, 0, 180})
		screen.DrawImage(darkOverlay, nil)
	} else {
		// Fallback to solid color if background image not available
		screen.Fill(color.RGBA{20, 20, 20, 255})
	}

	// Draw logo above the loading bar, if available
	if ls.logoImage != nil {
		logoOp := &ebiten.DrawImageOptions{}

		// Get logo dimensions
		logoWidth := float64(ls.logoImage.Bounds().Dx())
		logoHeight := float64(ls.logoImage.Bounds().Dy())

		// Position logo centered horizontally and above loading bar
		logoX := (game.ScreenWidth - logoWidth) / 2
		logoY := (game.ScreenHeight-logoHeight)/2 - 100 // 100 pixels above the center where loading bar is

		logoOp.GeoM.Translate(logoX, logoY)
		screen.DrawImage(ls.logoImage, logoOp)
	}

	// Draw progress bar
	ls.loadingBar.Draw(screen)
}

// HandleInput processes all input for this state
func (ls *LoadingState) HandleInput() error {
	// No input handling needed for loading screen
	return nil
}

// GetStateID returns a unique identifier for this state
func (ls *LoadingState) GetStateID() string {
	return "Loading"
}
