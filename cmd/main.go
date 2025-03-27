package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tile struct {
	X, Y int
}

type Player struct {
	X, Y  float64
	Speed float64
}

type Camera struct {
	X, Y float64
}

type Game struct {
	Tiles  []Tile
	Player Player
	Camera Camera
}

type TrailPoint struct {
	X, Y  float64
	Alpha float64
}

var trail []TrailPoint

const (
	tileSize     = 32
	tilesX       = 100
	tilesY       = 100
	screenWidth  = 800
	screenHeight = 600

	maxTrailLength = 20
)

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Player.Y -= g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Player.Y += g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Player.X -= g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Player.X += g.Player.Speed
	}

	g.Camera.X = g.Player.X - screenWidth/2
	g.Camera.Y = g.Player.Y - screenHeight/2

	trail = append(trail, TrailPoint{X: g.Player.X, Y: g.Player.Y, Alpha: 1.0})
	if len(trail) > maxTrailLength {
		trail = trail[1:]
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	for _, t := range g.Tiles {
		sx := float64(t.X*tileSize) - g.Camera.X
		sy := float64(t.Y*tileSize) - g.Camera.Y
		vector.DrawFilledRect(screen, float32(sx), float32(sy), float32(tileSize), float32(tileSize), color.RGBA{50, 50, 50, 255}, false)
	}

	px := g.Player.X - g.Camera.X
	py := g.Player.Y - g.Camera.Y
	vector.DrawFilledRect(screen, float32(px), float32(py), float32(tileSize), float32(tileSize), color.RGBA{255, 0, 0, 255}, false)

	for _, t := range trail {
		tx := t.X - g.Camera.X
		ty := t.Y - g.Camera.Y
		vector.DrawFilledRect(screen, float32(tx), float32(ty), float32(tileSize), float32(tileSize),
			color.RGBA{200, 0, 0, uint8(t.Alpha * 255)}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	tiles := make([]Tile, 0)

	for x := 0; x < tilesX; x++ {
		for y := 0; y < tilesY; y++ {
			tiles = append(tiles, Tile{X: x, Y: y})
		}
	}

	g := &Game{
		Tiles:  tiles,
		Player: Player{X: 1000, Y: 1000, Speed: 4},
		Camera: Camera{},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Phase 1: Movement + Camera")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
