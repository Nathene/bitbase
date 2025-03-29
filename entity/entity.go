package entity

import (
	"github.com/Nathene/bitbase/common"
	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	GetInventory() Inventory
	SetInventory(Inventory)
	GetX() float64
	GetY() float64
	Draw(screen *ebiten.Image, camera common.Camera)
}

type Inventory struct {
	Items []string
}
