package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

type ProgressBar struct {
	sprite marvtypes.Sprite
	area   marvtypes.MapBankArea
	rect   image.Rectangle

	currentPercent float64
	targetPercent  float64

	dirty bool
}

func NewProgressBar(rect image.Rectangle, sprite marvtypes.Sprite, area marvtypes.MapBankArea) *ProgressBar {
	return &ProgressBar{
		sprite: sprite,
		area:   area,
		rect:   rect,
	}
}

func (b *ProgressBar) Start() {
	b.sprite.ChangePos(b.rect)
	b.sprite.Show(GfxBankGfx, b.area)
	b.UpdatePercentages()
}

func (b *ProgressBar) Update() {
	if b.dirty {
		b.UpdatePercentages()
	}
}

func (b *ProgressBar) CurrentPercentage() float64 {
	return b.currentPercent
}

func (b *ProgressBar) SetCurrentPercentage(percentage float64) {
	b.currentPercent = percentage
	b.dirty = true
}

func (b *ProgressBar) SetTargetPercentage(percentage float64) {
	b.targetPercent = percentage
	b.dirty = true
}

func (b *ProgressBar) HitTarget() {
	b.currentPercent = b.targetPercent
	b.dirty = true
}

func (b *ProgressBar) UpdatePercentages() {
	//	r := b.rect // copy original rect
	//	r.Max.X = int(float64(r.Max.X) * b.currentPercent)
	//	b.sprite.ChangePos(r)

	r := b.rect // copy original rect

	curWidth := int(float64(b.rect.Max.X) * b.currentPercent)
	width := curWidth + int(float64(b.rect.Max.X)*(b.targetPercent-b.currentPercent))

	// clamp
	if width < b.rect.Max.X {
		r.Max.X = width
	} else {
		r.Max.X = b.rect.Max.X
	}

	b.sprite.ChangePos(r)

	b.sprite.ChangeViewport(image.Point{160 - curWidth, 0})
}
