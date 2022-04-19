package cartridge

import (
	"embed"
	"image"
	"math"

	"github.com/TheMightyGit/marv/marvlib"
	"github.com/TheMightyGit/marv/marvtypes"
)

//go:embed "resources/*"
var Resources embed.FS

const (
	GfxBankFont = iota
	GfxBankGfx
)
const (
	MapBankGfx = iota
)
const (
	MapAreaUI = iota
	MapAreaPanels
	MapAreaShip
	MapAreaStars
)
const (
	SpriteStars = iota
	SpritePlanetsStart
	SpritePlanetsEnd = iota + 10
	SpriteBaddieStart
	SpriteBaddieEnd = iota + 10
	SpriteShip
	SpritePanels
	SpriteUI
	SpriteButtonLetters
	SpriteGuessWord
)

var (
	spriteStars         marvtypes.Sprite
	spriteShip          marvtypes.Sprite
	spritePanels        marvtypes.Sprite
	spriteUI            marvtypes.Sprite
	spriteButtonLetters marvtypes.Sprite
	spriteGuessWord     marvtypes.Sprite

	buttonLettersArea marvtypes.MapBankArea
	guessWordArea     marvtypes.MapBankArea

	api = marvlib.API
)

type letterGfx [4]image.Point

var (
	letters = map[rune]letterGfx{
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
	}
)

func Start() {
	spriteStars = api.SpritesGet(SpriteStars)
	spriteStars.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	spriteStars.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaStars))

	spriteShip = api.SpritesGet(SpriteShip)
	spriteShip.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	spriteShip.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaShip))

	spritePanels = api.SpritesGet(SpritePanels)
	spritePanels.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	spritePanels.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaPanels))

	spriteUI = api.SpritesGet(SpriteUI)
	spriteUI.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	spriteUI.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaUI))

	spriteButtonLetters = api.SpritesGet(SpriteButtonLetters)
	spriteButtonLetters.ChangePos(image.Rectangle{image.Point{5 + (4 * 10), 5 + (11 * 10)}, image.Point{8 * 30, 30 * 3}})
	buttonLettersArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{8 * 3, 3 * 3})
	spriteButtonLetters.Show(GfxBankGfx, buttonLettersArea)

	pos := image.Point{0, 0}
	drawText(buttonLettersArea, pos, "ABCDEFGH", 3)
	pos = image.Point{0, 3}
	drawText(buttonLettersArea, pos, "IJKLMNOP", 3)
	pos = image.Point{0, 6}
	drawText(buttonLettersArea, pos, "QRSTUVWX", 3)

	spriteGuessWord = api.SpritesGet(SpriteGuessWord)
	spriteGuessWord.ChangePos(image.Rectangle{image.Point{80, 80}, image.Point{30 * 16, 30}})
	guessWordArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{16 * 2, 2})
	spriteGuessWord.Show(GfxBankGfx, guessWordArea)

	pos = image.Point{0, 0}
	drawText(guessWordArea, pos, "SOMEWORD", 2)

	api.SpritesSort()

	/*
		for _, w := range dictionary.Dictionary.Words() {
			api.ConsolePrintln(w)
		}
	*/
}

var (
	starsOffset image.Point
	cnt         float64
)

func Update() {
	spriteStars.ChangePos(image.Rectangle{starsOffset, image.Point{320, 200}})
	starsOffset.Y = -50 + int(math.Sin(cnt)*5)
	cnt += 0.05
}

func drawText(area marvtypes.MapBankArea, pos image.Point, txt string, spacing int) {
	for _, letter := range txt {
		area.Set(pos.Add(image.Point{0, 0}), uint8(letters[letter][0].X), uint8(letters[letter][0].Y), 0, 0)
		area.Set(pos.Add(image.Point{1, 0}), uint8(letters[letter][1].X), uint8(letters[letter][1].Y), 0, 0)
		area.Set(pos.Add(image.Point{0, 1}), uint8(letters[letter][2].X), uint8(letters[letter][2].Y), 0, 0)
		area.Set(pos.Add(image.Point{1, 1}), uint8(letters[letter][3].X), uint8(letters[letter][3].Y), 0, 0)
		pos = pos.Add(image.Point{spacing, 0})
	}
}
