package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Overlay struct {
	sprite marvtypes.Sprite
	area   marvtypes.MapBankArea
}

func NewOverlay(sprite marvtypes.Sprite, area marvtypes.MapBankArea) *Overlay {
	return &Overlay{
		sprite: sprite,
		area:   area,
	}
}

func (o *Overlay) Start() {
	o.sprite.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 100}})
	o.area.ClearWithColour(5, 5, 1, 0)
	o.area.StringToMap(image.Point{}, 1, 0, "Hello World!\nHi World!\nWotcha World!")
	o.sprite.Show(GfxBankSmallFont, o.area)
}

func (o *Overlay) Update() {
}
