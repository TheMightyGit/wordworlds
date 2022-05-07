package cartridge

import (
	"image"
	"math"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Baddie struct {
	name          string
	maxHitPoints  int
	hitPoints     int
	spriteBaddie  marvtypes.Sprite
	area          marvtypes.MapBankArea
	pos           image.Point
	baddiesOffset image.Point
	cnt           float64
}

func NewBaddie(name string, hitPoints int, pos image.Point, sprite marvtypes.Sprite, area marvtypes.MapBankArea) *Baddie {
	return &Baddie{
		name:         name,
		hitPoints:    hitPoints,
		maxHitPoints: hitPoints,
		spriteBaddie: sprite,
		area:         area,
		pos:          pos,
	}
}

func (b *Baddie) GetName() string {
	return b.name
}

func (b *Baddie) GetHealth() float64 {
	return float64(b.hitPoints) / float64(b.maxHitPoints)
}

func (b *Baddie) Damage(damage int) {
	b.hitPoints -= damage
	if b.hitPoints < 0 {
		b.hitPoints = 0
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
