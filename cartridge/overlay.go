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
	alertArea     marvtypes.MapBankArea
	statusArea    marvtypes.MapBankArea
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
	o.area.StringToMap(image.Point{3, 2}, 14, 16, "Used Words ")
	o.usedWordsArea = o.area.GetSubArea(image.Rectangle{
		image.Point{1 + 1, 2 + 1},
		image.Point{16, 12},
	})

	o.area.DrawBox(
		image.Rectangle{
			image.Point{19, 0},
			image.Point{60, 3},
		},
		outline,
		14,
		16,
	)
	o.area.StringToMap(image.Point{19 + 1 + 1, 0}, 14, 16, "Alert! ")
	o.alertArea = o.area.GetSubArea(image.Rectangle{
		image.Point{19 + 1, 1},
		image.Point{60, 3},
	})
	o.alertArea.StringToMap(image.Point{0, 0}, 14, 16, "Baddie #1: |||||||||||||||||||||||||")
	o.alertArea.StringToMap(image.Point{0, 1}, 14, 16, "Baddie #2: |||||||||||||||||||||||||")

	o.area.DrawBox(
		image.Rectangle{
			image.Point{62 + 1, 2},
			image.Point{62 + 16, 12},
		},
		outline,
		14,
		16,
	)
	o.area.StringToMap(image.Point{62 + 1 + 2, 2}, 14, 16, "Status ")
	o.statusArea = o.area.GetSubArea(image.Rectangle{
		image.Point{62 + 1 + 1, 2 + 1},
		image.Point{62 + 16, 12},
	})
	o.statusArea.StringToMap(image.Point{0, 0}, 14, 16, "Under Attack!\n20 damage\n to shield.")
}

func (o *Overlay) Update() {
}
