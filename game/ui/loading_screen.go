package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ProgressBar displays loading or other progress visually
type ProgressBar struct {
	X, Y          float64
	Width, Height float64
	Progress      float64 // 0.0 to 1.0

	BackgroundColor color.Color
	FillColor       color.Color
	BorderColor     color.Color
	BorderWidth     float64

	// Animation properties
	AnimateProgress bool
	CurrentDisplay  float64
	AnimationSpeed  float64
}

// NewProgressBar creates a new progress bar
func NewProgressBar(x, y, width, height float64) *ProgressBar {
	return &ProgressBar{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Progress:        0,
		CurrentDisplay:  0,
		BackgroundColor: color.RGBA{40, 40, 40, 255},
		FillColor:       color.RGBA{65, 135, 230, 255},
		BorderColor:     color.RGBA{200, 200, 200, 255},
		BorderWidth:     2,
		AnimateProgress: true,
		AnimationSpeed:  0.05,
	}
}

// Update animates the progress bar if animation is enabled
func (pb *ProgressBar) Update() {
	if pb.AnimateProgress {
		if pb.CurrentDisplay < pb.Progress {
			pb.CurrentDisplay += pb.AnimationSpeed
			if pb.CurrentDisplay > pb.Progress {
				pb.CurrentDisplay = pb.Progress
			}
		} else if pb.CurrentDisplay > pb.Progress {
			pb.CurrentDisplay -= pb.AnimationSpeed
			if pb.CurrentDisplay < pb.Progress {
				pb.CurrentDisplay = pb.Progress
			}
		}
	} else {
		pb.CurrentDisplay = pb.Progress
	}
}

// Draw renders the progress bar
func (pb *ProgressBar) Draw(screen *ebiten.Image) {
	// Draw background
	vector.DrawFilledRect(screen,
		float32(pb.X), float32(pb.Y),
		float32(pb.Width), float32(pb.Height),
		pb.BackgroundColor, false)

	// Draw fill based on progress
	fillWidth := pb.Width * pb.CurrentDisplay
	vector.DrawFilledRect(screen,
		float32(pb.X), float32(pb.Y),
		float32(fillWidth), float32(pb.Height),
		pb.FillColor, false)

	// Draw border
	vector.StrokeRect(screen,
		float32(pb.X), float32(pb.Y),
		float32(pb.Width), float32(pb.Height),
		float32(pb.BorderWidth), pb.BorderColor, false)
}
