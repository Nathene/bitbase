package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

// Button states
const (
	ButtonNormal = iota
	ButtonHovered
	ButtonPressed
	ButtonDisabled
)

// Button is an interactive UI element
type Button struct {
	X, Y          float64
	Width, Height float64
	Text          string // Kept for identification purposes only, not rendered

	// Visual properties
	Font            font.Face // Kept for future implementation
	TextColor       color.Color
	BackgroundColor color.Color
	HoverColor      color.Color
	PressedColor    color.Color
	DisabledColor   color.Color
	BorderColor     color.Color
	BorderWidth     float64

	// State
	State    int
	Disabled bool

	// Event callbacks
	OnClick func()
}

// NewButton creates a new button with default styling
func NewButton(x, y, width, height float64, text string, font font.Face) *Button {
	return &Button{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Text:            text,
		Font:            font,
		TextColor:       color.White,
		BackgroundColor: color.RGBA{60, 60, 60, 255},
		HoverColor:      color.RGBA{80, 80, 80, 255},
		PressedColor:    color.RGBA{40, 40, 40, 255},
		DisabledColor:   color.RGBA{100, 100, 100, 128},
		BorderColor:     color.RGBA{200, 200, 200, 255},
		BorderWidth:     2,
		State:           ButtonNormal,
	}
}

// Update handles button state based on mouse input
func (btn *Button) Update() {
	if btn.Disabled {
		btn.State = ButtonDisabled
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()
	isHovered := float64(mouseX) >= btn.X &&
		float64(mouseX) <= btn.X+btn.Width &&
		float64(mouseY) >= btn.Y &&
		float64(mouseY) <= btn.Y+btn.Height

	if isHovered {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			btn.State = ButtonPressed
		} else {
			if btn.State == ButtonPressed && btn.OnClick != nil {
				btn.OnClick()
			}
			btn.State = ButtonHovered
		}
	} else {
		btn.State = ButtonNormal
	}
}

// Draw renders the button
func (btn *Button) Draw(screen *ebiten.Image) {
	// Determine the fill color based on button state
	var fillColor color.Color
	switch btn.State {
	case ButtonHovered:
		fillColor = btn.HoverColor
	case ButtonPressed:
		fillColor = btn.PressedColor
	case ButtonDisabled:
		fillColor = btn.DisabledColor
	default:
		fillColor = btn.BackgroundColor
	}

	// Draw button background
	vector.DrawFilledRect(screen,
		float32(btn.X), float32(btn.Y),
		float32(btn.Width), float32(btn.Height),
		fillColor, false)

	// Draw border
	vector.StrokeRect(screen,
		float32(btn.X), float32(btn.Y),
		float32(btn.Width), float32(btn.Height),
		float32(btn.BorderWidth), btn.BorderColor, false)

	// Note: Text rendering has been removed
}
