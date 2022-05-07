package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

var outline = &[9]image.Point{
	{34, 0}, {35, 0}, {36, 0},
	{34, 1}, {35, 1}, {36, 1},
	{34, 2}, {35, 2}, {36, 2},
}

type Overlay struct {
	sprite        marvtypes.Sprite
	area          marvtypes.MapBankArea
	usedWordsArea marvtypes.MapBankArea
	cursorPos     image.Point
}

func NewOverlay(sprite marvtypes.Sprite, area marvtypes.MapBankArea) *Overlay {
	return &Overlay{
		sprite: sprite,
		area:   area,
	}
}

func (o *Overlay) AddWord(txt string) {
	o.usedWordsArea.StringToMap(o.cursorPos, 14, 16, txt)
	o.cursorPos.Y++
}

func (o *Overlay) Start() {
	o.sprite.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 100}})
	o.area.ClearWithColour(0, 0, 14, 16)
	o.sprite.Show(GfxBankSmallFont, o.area)
	o.area.DrawBox(
		image.Rectangle{
			image.Point{1, 2},
			image.Point{16, 12},
		},
		outline,
		14,
		16,
	)
	o.usedWordsArea = o.area.GetSubArea(image.Rectangle{
		image.Point{2, 3},
		image.Point{16, 12},
	})
	o.cursorPos = o.usedWordsArea.StringToMap(image.Point{}, 14, 16, "Used Words:\n")
}

func (o *Overlay) Update() {
}
