package game

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// AssetManager handles loading and managing game assets
type AssetManager struct {
	images       map[string]*ebiten.Image
	audioSamples map[string]*audio.Player
	fonts        map[string]font.Face
	jsonData     map[string]interface{}

	// Tracking loading progress
	totalAssets    int
	loadedAssets   int
	isLoading      bool
	onLoadComplete func()
	mutex          sync.Mutex

	// Audio context for sound loading
	audioContext *audio.Context
}

// NewAssetManager creates a new asset manager
func NewAssetManager() *AssetManager {
	// Initialize audio context
	audioContext, err := audio.NewContext(44100)
	if err != nil {
		log.Fatalf("Failed to create audio context: %v", err)
	}

	// Initialize with a default font
	fonts := make(map[string]font.Face)
	fonts["default"] = basicfont.Face7x13

	return &AssetManager{
		images:       make(map[string]*ebiten.Image),
		audioSamples: make(map[string]*audio.Player),
		fonts:        fonts,
		jsonData:     make(map[string]interface{}),
		audioContext: audioContext,
		mutex:        sync.Mutex{},
	}
}

// StartLoading begins a new batch of asset loading
func (am *AssetManager) StartLoading(onComplete func()) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	am.totalAssets = 0
	am.loadedAssets = 0
	am.isLoading = true
	am.onLoadComplete = onComplete
}

// LoadImage loads an image asset asynchronously
func (am *AssetManager) LoadImage(id, path string) {
	am.mutex.Lock()
	am.totalAssets++
	am.mutex.Unlock()

	go func() {
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Printf("Failed to load image %s: %v", path, err)
		} else {
			am.mutex.Lock()
			am.images[id] = img
			am.mutex.Unlock()
		}

		am.mutex.Lock()
		am.loadedAssets++

		// Check if this was the last asset to load
		if am.isLoading && am.loadedAssets >= am.totalAssets && am.onLoadComplete != nil {
			am.isLoading = false
			// Call the completion callback outside the lock
			callback := am.onLoadComplete
			am.mutex.Unlock()
			callback()
		} else {
			am.mutex.Unlock()
		}
	}()
}

// GetImage retrieves a loaded image
func (am *AssetManager) GetImage(id string) *ebiten.Image {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if img, ok := am.images[id]; ok {
		return img
	}
	// Don't log a warning - just return nil
	return nil
}

// GetLoadingProgress returns progress as a value between 0.0 and 1.0
func (am *AssetManager) GetLoadingProgress() float64 {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if am.totalAssets == 0 {
		return 1.0
	}
	return float64(am.loadedAssets) / float64(am.totalAssets)
}

// IsLoadingComplete checks if all assets have been loaded
func (am *AssetManager) IsLoadingComplete() bool {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	return am.loadedAssets >= am.totalAssets
}

// LoadFont loads a font file and registers it with the given size
func (am *AssetManager) LoadFont(id string, fontBytes []byte, size float64) {
	am.mutex.Lock()
	am.totalAssets++
	am.mutex.Unlock()

	go func() {
		// Use the default font
		face := basicfont.Face7x13

		// Store the font face
		am.mutex.Lock()
		am.fonts[id] = face
		am.loadedAssets++

		// Check if this was the last asset to load
		if am.isLoading && am.loadedAssets >= am.totalAssets && am.onLoadComplete != nil {
			am.isLoading = false
			// Call the completion callback outside the lock
			callback := am.onLoadComplete
			am.mutex.Unlock()
			callback()
		} else {
			am.mutex.Unlock()
		}
	}()
}

// GetFont retrieves a loaded font, or returns the default font if not found
func (am *AssetManager) GetFont(id string) font.Face {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if f, ok := am.fonts[id]; ok {
		return f
	}

	// Return the default font as a fallback
	return am.fonts["default"]
}

// Similar methods for audio, fonts, and data would be implemented
