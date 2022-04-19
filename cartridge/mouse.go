package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Clickable interface {
	OnClick(image.Point) bool
}

type Pointer struct {
	sprite     marvtypes.Sprite
	area       marvtypes.MapBankArea
	clickables []Clickable
}

func NewPointer(sprite marvtypes.Sprite, area marvtypes.MapBankArea, clickables ...Clickable) *Pointer {
	return &Pointer{
		sprite:     sprite,
		area:       area,
		clickables: clickables,
	}
}

func (s *Pointer) Start() {
	s.sprite = api.SpritesGet(SpriteMousePointer)
	s.sprite.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{10, 10}})
	s.area.Set(image.Point{}, 1, 0, 0, 0)
	s.sprite.Show(GfxBankGfx, s.area)
}

func (s *Pointer) Update() {
	pos := api.InputMousePos()
	s.sprite.ChangePos(image.Rectangle{pos, image.Point{10, 10}})

	if api.InputMousePressed() {
		for _, clickable := range s.clickables {
			if clickable.OnClick(pos) {
				break
			}
		}
	}
}
