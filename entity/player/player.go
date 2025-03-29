package player

import (
	"github.com/Nathene/bitbase/common"
	"github.com/Nathene/bitbase/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	x, y          float64
	Speed         float64
	inventory     entity.Inventory
	ShowInventory bool

	AnimTimer float64
	AnimFrame int
}

// NewInventory creates a new empty inventory
func NewInventory() entity.Inventory {
	return entity.Inventory{
		Items: []string{},
	}
}

func (p *Player) GetInventory() entity.Inventory {
	return p.inventory
}

func (p *Player) SetInventory(inventory entity.Inventory) {
	p.inventory = inventory
}

func (p *Player) GetX() float64 {
	return p.x
}

func (p *Player) SetX(x float64) {
	p.x = x
}

func (p *Player) GetY() float64 {
	return p.y
}

func (p *Player) SetY(y float64) {
	p.y = y
}

// Draw draws the player on the screen
func (p *Player) Draw(screen *ebiten.Image, camera common.Camera) {
	// Drawing is handled in the game package for now
}
