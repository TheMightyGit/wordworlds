package cartridge

import (
	"embed"
	"image"
	"math/rand"

	"github.com/TheMightyGit/marv/marvlib"
	"github.com/TheMightyGit/marv/marvtypes"
)

//go:embed "resources/*"
var Resources embed.FS

const (
	GfxBankFont = iota
	GfxBankDinos
	GfxBankBG
)
const (
	MapBankDinos = iota
)
const (
	MapBankDinoAnimsArea = 0
	MapBankBGArea        = 1
)
const (
	SpriteBG1 = iota
	SpriteBG2
	SpriteStart
	SpriteEnd = 127
)

var (
	bg1 marvtypes.Sprite
	bg2 marvtypes.Sprite
	api = marvlib.API
)

func randomArenaStartPos(colour int) image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: Arena.Min.X + (colour * Arena.Size().X / 4) + rand.Intn(Arena.Size().X/4),
			Y: Arena.Min.Y + rand.Intn(Arena.Size().Y),
		},
		Max: DinoSize,
	}
}

func Start() {
	bg1 = api.SpritesGet(SpriteBG1)
	bg1.ChangePos(image.Rect(0, 0, 320, 200))
	bg1.Show(GfxBankBG, api.MapBanksGet(MapBankDinos).GetArea(MapBankBGArea))

	bg2 = api.SpritesGet(SpriteBG2)
	bg2.ChangePos(image.Rect(320, 0, 320, 200))
	bg2.Show(GfxBankBG, api.MapBanksGet(MapBankDinos).GetArea(MapBankBGArea))

	for i := SpriteStart; i <= SpriteEnd; i++ {
		colour := rand.Intn(4)
		Dinos = append(Dinos, NewDino(api.SpritesGet(i), randomArenaStartPos(colour), colour))
	}
	// marv.ModBanks[0].Play()
	api.SfxBanksGet(0).PlayLooped()
}

var (
	bgX           int
	bgScrollScale = 4
)

func Update() {
	bg1.ChangePos(image.Rect((bgX / bgScrollScale), 0, 320, 200))
	bg2.ChangePos(image.Rect((bgX/bgScrollScale)+320, 0, 320, 200))
	bgX--
	if bgX == (-320 * bgScrollScale) {
		bgX = 0
	}
	for _, dino := range Dinos {
		dino.Update()
	}
	api.SpritesSort()
}
