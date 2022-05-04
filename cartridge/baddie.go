package cartridge

import (
	"image"
	"math"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Baddie struct {
	spriteBaddie  marvtypes.Sprite
	area          marvtypes.MapBankArea
	pos           image.Point
	baddiesOffset image.Point
	cnt           float64
}

func NewBaddie(pos image.Point, sprite marvtypes.Sprite, area marvtypes.MapBankArea) *Baddie {
	return &Baddie{
		spriteBaddie: sprite,
		area:         area,
		pos:          pos,
	}
}

func (b *Baddie) Start() {
	b.spriteBaddie.ChangePos(image.Rectangle{b.pos.Add(b.baddiesOffset), image.Point{5 * 10, 4 * 10}})
	b.spriteBaddie.Show(GfxBankGfx, b.area)
}

func (b *Baddie) Update() {
	b.spriteBaddie.ChangePos(image.Rectangle{b.pos.Add(b.baddiesOffset), image.Point{5 * 10, 4 * 10}})
	b.baddiesOffset.Y = int(math.Sin(b.cnt/2)*10) - 5
	b.baddiesOffset.X = int(math.Sin(b.cnt)*80) - 40
	b.cnt += 0.01
}
