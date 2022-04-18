package cartridge

import (
	"embed"
	"image"

	"github.com/TheMightyGit/wordworlds/dictionary"

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
)

var (
	spriteStars  marvtypes.Sprite
	spriteShip   marvtypes.Sprite
	spritePanels marvtypes.Sprite
	spriteUI     marvtypes.Sprite

	api = marvlib.API
)

func Start() {
	spriteStars = api.SpritesGet(SpriteStars)
	spriteStars.ChangePos(image.Rect(0, 0, 320, 200))
	spriteStars.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaStars))

	spriteShip = api.SpritesGet(SpriteShip)
	spriteShip.ChangePos(image.Rect(0, 0, 320, 200))
	spriteShip.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaShip))

	spritePanels = api.SpritesGet(SpritePanels)
	spritePanels.ChangePos(image.Rect(0, 0, 320, 200))
	spritePanels.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaPanels))

	spriteUI = api.SpritesGet(SpriteUI)
	spriteUI.ChangePos(image.Rect(0, 0, 320, 200))
	spriteUI.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaUI))

	api.SpritesSort()

	for _, w := range dictionary.Dictionary.Words() {
		api.ConsolePrintln(w)
	}
}

func Update() {
}
