package game

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/Nathene/bitbase/common"
	"github.com/Nathene/bitbase/entity/player"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Define screen dimensions as constants
const (
	ScreenWidth  = 1920
	ScreenHeight = 1080

	tileSize = 32
	tilesX   = 100
	tilesY   = 100

	playerWidth    = tileSize
	playerHeight   = tileSize
	maxTrailLength = 20

	playerSheetStartX = 43  // Note 2: Assumed 0, adjust if first frame has left padding
	playerSheetStartY = 23  // Note 2: Assumed 0, adjust if first frame has top padding
	playerFrameWidth  = 11  // Your measurement: Width of a single frame
	playerFrameHeight = 16  // Your measurement: Height of a single frame
	playerFrameStepX  = 96  // Your measurement: Horizontal distance between frame starts
	playerFrameCount  = 9   // Note 3: Make sure this matches the actual number of frames!
	playerAnimSpeed   = 0.1 // Adjust for desired animation speed (lower = faster)
	playerDrawScale   = 3.0
)

type Tile struct {
	X, Y  int
	Color color.RGBA
}

type TileProperty struct {
	IsWall bool
}

type Game struct {
	Tiles    []Tile
	Player   player.Player
	Camera   common.Camera
	WorldMap [][]TileProperty

	PlayerSheet     *ebiten.Image
	BackgroundImage *ebiten.Image
}

// NewGame creates a new game instance with initialized components
func NewGame(playerSheet, backgroundImage *ebiten.Image) *Game {
	tiles := make([]Tile, 0)

	for x := 0; x < tilesX; x++ {
		for y := 0; y < tilesY; y++ {
			tileColor := color.RGBA{50, 50, 50, 255}
			if (x+y)%2 == 0 {
				tileColor = color.RGBA{70, 70, 70, 255}
			}
			tiles = append(tiles, Tile{X: x, Y: y, Color: tileColor})
		}
	}

	worldMap := make([][]TileProperty, tilesY)

	for row := range worldMap {
		worldMap[row] = make([]TileProperty, tilesX)
		for col := range worldMap[row] {
			if row == 0 || row == tilesY-1 || col == 0 || col == tilesX-1 {
				worldMap[row][col] = TileProperty{IsWall: true}
			}
		}
	}

	p := player.Player{}
	p.SetX(1000)
	p.SetY(1000)
	p.Speed = 4
	p.SetInventory(player.NewInventory())

	return &Game{
		Tiles:           tiles,
		Player:          p,
		Camera:          common.Camera{},
		WorldMap:        worldMap,
		PlayerSheet:     playerSheet,
		BackgroundImage: backgroundImage,
	}
}

func (g *Game) Update() error {
	var dx, dy float64
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= g.Player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += g.Player.Speed
	}

	playerMoved := dx != 0 || dy != 0

	playerX, playerY := g.Player.GetX(), g.Player.GetY()

	// --- COLLISION DETECTION ---
	nextX := playerX + dx
	nextY := playerY + dy

	if dx != 0 {
		if !g.collidesWithWorld(nextX, playerY) {
			g.Player.SetX(nextX)
		}
	}

	if dy != 0 {
		if !g.collidesWithWorld(playerX, nextY) {
			g.Player.SetY(nextY)
		}
	}
	if playerMoved {
		g.Player.AnimFrame = 0 // Reset frame when moving (adjust if needed)
		g.Player.AnimTimer = 0
	} else { // When IDLE (not moving)
		// --- Animate the idle state using stable ActualTPS ---

		// Calculate time elapsed since last Update tick using ActualTPS
		var deltaT float64
		actualTps := ebiten.ActualTPS()
		if actualTps > 0 { // Avoid division by zero if TPS is somehow zero
			deltaT = 1.0 / actualTps
		} else {
			// Fallback or assume target TPS if ActualTPS is zero (shouldn't happen often)
			deltaT = 1.0 / 60.0 // Assuming target is 60 if unavailable
		}

		g.Player.AnimTimer += deltaT // Increment timer by calculated delta time

		if g.Player.AnimTimer >= playerAnimSpeed {
			g.Player.AnimTimer -= playerAnimSpeed // Reset timer partially
			g.Player.AnimFrame++
			if g.Player.AnimFrame >= playerFrameCount { // Frame count is 9
				g.Player.AnimFrame = 0 // Loop the idle animation
			}
		}
	}
	g.Camera.X = g.Player.GetX() - ScreenWidth/2
	g.Camera.Y = g.Player.GetY() - ScreenHeight/2

	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		g.Player.ShowInventory = !g.Player.ShowInventory
	}

	return nil
}

func (g *Game) collidesWithWorld(checkX, checkY float64) bool {
	minX := checkX
	maxX := checkX + playerWidth
	minY := checkY
	maxY := checkY + playerHeight

	minTileX := int(minX / tileSize)
	maxTileX := int((maxX - 1) / tileSize)
	minTileY := int(minY / tileSize)
	maxTileY := int((maxY - 1) / tileSize)

	if g.WorldMap == nil {
		return true
	}

	for ty := minTileY; ty <= maxTileY; ty++ {
		for tx := minTileX; tx <= maxTileX; tx++ {
			if ty < 0 || ty >= len(g.WorldMap) || tx < 0 || tx >= len(g.WorldMap[0]) {
				return true
			}
			if g.WorldMap[ty][tx].IsWall {
				return true
			}
		}
	}
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	// --- Draw the World Background Image ---
	if g.BackgroundImage != nil {
		opts := &ebiten.DrawImageOptions{}

		// Translate the background based on the camera's position
		// Move the background opposite to the camera's view
		opts.GeoM.Translate(-g.Camera.X, -g.Camera.Y)

		screen.DrawImage(g.BackgroundImage, opts)
	} else {
		// Fallback if background failed to load
		screen.Fill(color.RGBA{50, 50, 50, 255})
	}

	if g.PlayerSheet != nil {
		frameStartX := playerSheetStartX + (g.Player.AnimFrame * playerFrameStepX)
		frameStartY := playerSheetStartY

		sourceRect := image.Rect(
			frameStartX,
			frameStartY,
			frameStartX+playerFrameWidth,
			frameStartY+playerFrameHeight,
		)

		frameToDraw := g.PlayerSheet.SubImage(sourceRect).(*ebiten.Image)

		opts := &ebiten.DrawImageOptions{}

		opts.GeoM.Scale(playerDrawScale, playerDrawScale)

		playerScreenX := g.Player.GetX() - g.Camera.X
		playerScreenY := g.Player.GetY() - g.Camera.Y
		opts.GeoM.Translate(playerScreenX, playerScreenY)

		screen.DrawImage(frameToDraw, opts)
	} else {
		log.Println("Player sheet is nil")
		px := g.Player.GetX() - g.Camera.X
		py := g.Player.GetY() - g.Camera.Y
		vector.DrawFilledRect(screen, float32(px), float32(py), float32(tileSize), float32(tileSize), color.RGBA{255, 0, 0, 255}, false)
	}

	// In Game.Draw() near the end
	debugText := fmt.Sprintf("X: %.1f, Y: %.1f | Frame: %d", g.Player.GetX(), g.Player.GetY(), g.Player.AnimFrame)
	ebitenutil.DebugPrint(screen, debugText)

	if g.Player.ShowInventory {
		inventoryText := "Inventory:\n"
		for i, item := range g.Player.GetInventory().Items {
			inventoryText += fmt.Sprintf("%d: %s\n", i+1, item)
		}
		ebitenutil.DebugPrint(screen, inventoryText)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Run() {
	g.Player.AnimFrame = 0
	g.Player.AnimTimer = 0

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Phase 1: Movement + Camera")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
