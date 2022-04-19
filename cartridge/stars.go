package cartridge

import (
	"image"
	"math"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Stars struct {
	spriteStars marvtypes.Sprite
	area        marvtypes.MapBankArea
	starsOffset image.Point
	cnt         float64
}

func NewStars(sprite marvtypes.Sprite, area marvtypes.MapBankArea) *Stars {
	return &Stars{
		spriteStars: sprite,
		area:        area,
	}
}

func (s *Stars) Start() {
	s.spriteStars.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	s.spriteStars.Show(GfxBankGfx, s.area)
}

func (s *Stars) Update() {
	s.spriteStars.ChangePos(image.Rectangle{s.starsOffset, image.Point{320, 200}})
	s.starsOffset.Y = -50 + int(math.Sin(s.cnt)*5)
	s.cnt += 0.05
}
