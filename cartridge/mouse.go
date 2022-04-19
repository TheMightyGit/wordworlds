package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Pointer struct {
	sprite marvtypes.Sprite
	area   marvtypes.MapBankArea
}

func NewPointer(sprite marvtypes.Sprite, area marvtypes.MapBankArea) *Pointer {
	return &Pointer{
		sprite: sprite,
		area:   area,
	}
}

func (s *Pointer) Start() {
	s.sprite = api.SpritesGet(SpriteMousePointer)
	s.sprite.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{10, 10}})
	s.area.Set(image.Point{}, 1, 0, 0, 0)
	s.sprite.Show(GfxBankGfx, s.area)
}

func (s *Pointer) Update() {
	s.sprite.ChangePos(image.Rectangle{api.InputMousePos(), image.Point{10, 10}})
}
