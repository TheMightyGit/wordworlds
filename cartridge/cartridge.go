package cartridge

import (
	"embed"
	"image"
	"math/rand"

	"github.com/TheMightyGit/marv/marvtypes"
)

//go:embed "resources/*"
var Resources embed.FS

const (
	GfxBankFont = iota
	GfxBankGfx
	GfxBankSmallFont
)
const (
	MapBankGfx = iota
	MapBankSmallFont
)
const (
	MapAreaUI = iota
	MapAreaPanels
	MapAreaShip
	MapAreaStars
	MapAreaBars
	MapAreaBaddies
)
const (
	SpriteStars = iota
	SpritePlanetsStart
	_
	_
	_
	SpritePlanetsEnd
	SpriteBaddieStart
	_
	_
	_
	_
	_
	_
	_
	_
	SpriteBaddieEnd
	SpriteShip
	SpritePanels
	SpriteUI
	SpriteWeaponProgressBar
	SpriteHullProgressBar
	SpriteShieldProgressBar
	SpriteButtonLetters
	SpriteGuessWord
	SpriteSmallFontOverlay

	SpriteMousePointer = 127
)

var (
	api marvtypes.MarvAPI

	stars   *Stars
	baddies []*Baddie
	ship    *Ship
	pointer *Pointer
	overlay *Overlay
)

type letterGfx [4]image.Point

var (
	letters = map[rune]letterGfx{
		' ': letterGfx{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}},
		'A': letterGfx{{X: 0, Y: 19}, {X: 1, Y: 19}, {X: 0, Y: 20}, {X: 1, Y: 20}},
		'B': letterGfx{{X: 2, Y: 19}, {X: 3, Y: 19}, {X: 2, Y: 20}, {X: 3, Y: 20}},
		'C': letterGfx{{X: 4, Y: 19}, {X: 5, Y: 19}, {X: 4, Y: 20}, {X: 5, Y: 20}},
		'D': letterGfx{{X: 6, Y: 19}, {X: 7, Y: 19}, {X: 6, Y: 20}, {X: 7, Y: 20}},
		'E': letterGfx{{X: 8, Y: 19}, {X: 9, Y: 19}, {X: 8, Y: 20}, {X: 9, Y: 20}},
		'F': letterGfx{{X: 10, Y: 19}, {X: 11, Y: 19}, {X: 10, Y: 20}, {X: 11, Y: 20}},
		'G': letterGfx{{X: 12, Y: 19}, {X: 13, Y: 19}, {X: 12, Y: 20}, {X: 13, Y: 20}},
		'H': letterGfx{{X: 14, Y: 19}, {X: 15, Y: 19}, {X: 14, Y: 20}, {X: 15, Y: 20}},
		'I': letterGfx{{X: 16, Y: 19}, {X: 17, Y: 19}, {X: 16, Y: 20}, {X: 17, Y: 20}},
		'J': letterGfx{{X: 18, Y: 19}, {X: 19, Y: 19}, {X: 18, Y: 20}, {X: 19, Y: 20}},
		'K': letterGfx{{X: 20, Y: 19}, {X: 21, Y: 19}, {X: 20, Y: 20}, {X: 21, Y: 20}},
		'L': letterGfx{{X: 22, Y: 19}, {X: 23, Y: 19}, {X: 22, Y: 20}, {X: 23, Y: 20}},
		'M': letterGfx{{X: 0, Y: 21}, {X: 1, Y: 21}, {X: 0, Y: 22}, {X: 1, Y: 22}},
		'N': letterGfx{{X: 2, Y: 21}, {X: 3, Y: 21}, {X: 2, Y: 22}, {X: 3, Y: 22}},
		'O': letterGfx{{X: 4, Y: 21}, {X: 5, Y: 21}, {X: 4, Y: 22}, {X: 5, Y: 22}},
		'P': letterGfx{{X: 6, Y: 21}, {X: 7, Y: 21}, {X: 6, Y: 22}, {X: 7, Y: 22}},
		'Q': letterGfx{{X: 8, Y: 21}, {X: 9, Y: 21}, {X: 8, Y: 22}, {X: 9, Y: 22}},
		'R': letterGfx{{X: 10, Y: 21}, {X: 11, Y: 21}, {X: 10, Y: 22}, {X: 11, Y: 22}},
		'S': letterGfx{{X: 12, Y: 21}, {X: 13, Y: 21}, {X: 12, Y: 22}, {X: 13, Y: 22}},
		'T': letterGfx{{X: 14, Y: 21}, {X: 15, Y: 21}, {X: 14, Y: 22}, {X: 15, Y: 22}},
		'U': letterGfx{{X: 16, Y: 21}, {X: 17, Y: 21}, {X: 16, Y: 22}, {X: 17, Y: 22}},
		'V': letterGfx{{X: 18, Y: 21}, {X: 19, Y: 21}, {X: 18, Y: 22}, {X: 19, Y: 22}},
		'W': letterGfx{{X: 20, Y: 21}, {X: 21, Y: 21}, {X: 20, Y: 22}, {X: 21, Y: 22}},
		'X': letterGfx{{X: 22, Y: 21}, {X: 23, Y: 21}, {X: 22, Y: 22}, {X: 23, Y: 22}},
		'Y': letterGfx{{X: 0, Y: 23}, {X: 1, Y: 23}, {X: 0, Y: 24}, {X: 1, Y: 24}},
		'Z': letterGfx{{X: 2, Y: 23}, {X: 3, Y: 23}, {X: 2, Y: 24}, {X: 3, Y: 24}},
		'-': letterGfx{{X: 4, Y: 23}, {X: 5, Y: 23}, {X: 4, Y: 24}, {X: 5, Y: 24}},
		//
		's': letterGfx{{X: 18, Y: 23}, {X: 19, Y: 23}, {X: 18, Y: 24}, {X: 19, Y: 24}}, // SHUFFLE icon
		'd': letterGfx{{X: 20, Y: 23}, {X: 21, Y: 23}, {X: 20, Y: 24}, {X: 21, Y: 24}}, // DEL icon
		'o': letterGfx{{X: 22, Y: 23}, {X: 23, Y: 23}, {X: 22, Y: 24}, {X: 23, Y: 24}}, // OK icon
	}
)

type Updateable interface {
	Start()
	Update()
}

func Start(fooApi marvtypes.MarvAPI) {
	api = fooApi

	stars = NewStars(
		api.SpritesGet(SpriteStars),
		api.MapBanksGet(MapBankGfx).GetArea(MapAreaStars),
	)

	baddies = append(baddies, NewBaddie(
		"Interceptor",
		100,
		image.Point{140, 35},
		api.SpritesGet(SpriteBaddieStart+0),
		api.MapBanksGet(MapBankGfx).GetArea(MapAreaBaddies),
	))
	baddies = append(baddies, NewBaddie(
		"Interceptor",
		100,
		image.Point{180, 45},
		api.SpritesGet(SpriteBaddieStart+1),
		api.MapBanksGet(MapBankGfx).GetArea(MapAreaBaddies),
	))
	baddies = append(baddies, NewBaddie(
		"Interceptor",
		100,
		image.Point{200, 40},
		api.SpritesGet(SpriteBaddieStart+2),
		api.MapBanksGet(MapBankGfx).GetArea(MapAreaBaddies),
	))
	overlay = NewOverlay(
		api.SpritesGet(SpriteSmallFontOverlay),
		api.MapBanksGet(MapBankSmallFont).AllocArea(image.Point{80, 34}),
	)

	ship = NewShip(
		api.SpritesGet(SpriteShip),
		api.MapBanksGet(MapBankGfx).GetArea(MapAreaShip),
		api.SpritesGet(SpritePanels),
		api.MapBanksGet(MapBankGfx).GetArea(MapAreaPanels),
		api.SpritesGet(SpriteUI),
		overlay,
		func() []*Baddie {
			return baddies
		},
	)

	pointer = NewPointer(
		api.SpritesGet(SpriteMousePointer),
		api.MapBanksGet(MapBankGfx).AllocArea(image.Point{1, 1}),
		ship,
	)

	stars.Start()
	overlay.Start() // overlay needs to start before ship.
	pointer.Start()

	for _, b := range baddies {
		// b.Damage(rand.Intn(b.maxHitPoints))
		b.cnt = rand.Float64()
		b.Start()
	}

	ship.Start()

	api.SpritesSort()
}

func Update() {
	stars.Update()
	ship.Update()
	overlay.Update()
	pointer.Update()
	for _, b := range baddies {
		b.Update()
	}
}

func drawText(area marvtypes.MapBankArea, pos image.Point, txt string, spacing int) {
	for _, letter := range txt {
		drawRune(area, pos, letter)
		pos = pos.Add(image.Point{spacing, 0})
	}
}

func drawRune(area marvtypes.MapBankArea, pos image.Point, letter rune) {
	area.Set(pos.Add(image.Point{0, 0}), uint8(letters[letter][0].X), uint8(letters[letter][0].Y), 0, 0)
	area.Set(pos.Add(image.Point{1, 0}), uint8(letters[letter][1].X), uint8(letters[letter][1].Y), 0, 0)
	area.Set(pos.Add(image.Point{0, 1}), uint8(letters[letter][2].X), uint8(letters[letter][2].Y), 0, 0)
	area.Set(pos.Add(image.Point{1, 1}), uint8(letters[letter][3].X), uint8(letters[letter][3].Y), 0, 0)
}
