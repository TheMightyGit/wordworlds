package cartridge

import (
	"image"

	"github.com/TheMightyGit/marv/marvtypes"
)

type Ship struct {
	spriteShip   marvtypes.Sprite
	areaShip     marvtypes.MapBankArea
	spritePanels marvtypes.Sprite
	areaPanels   marvtypes.MapBankArea
	spriteUI     marvtypes.Sprite
	areaUI       marvtypes.MapBankArea
}

func NewShip(
	spriteShip marvtypes.Sprite,
	areaShip marvtypes.MapBankArea,
	spritePanels marvtypes.Sprite,
	areaPanels marvtypes.MapBankArea,
	spriteUI marvtypes.Sprite,
	areaUI marvtypes.MapBankArea,
) *Ship {
	return &Ship{
		spriteShip:   spriteShip,
		areaShip:     areaShip,
		spritePanels: spritePanels,
		areaPanels:   areaPanels,
		spriteUI:     spriteUI,
		areaUI:       areaUI,
	}
}

func (s *Ship) Start() {
	s.spriteShip = api.SpritesGet(SpriteShip)
	s.spriteShip.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	s.spriteShip.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaShip))

	s.spritePanels = api.SpritesGet(SpritePanels)
	s.spritePanels.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	s.spritePanels.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaPanels))

	s.spriteUI = api.SpritesGet(SpriteUI)
	s.spriteUI.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	s.spriteUI.Show(GfxBankGfx, api.MapBanksGet(MapBankGfx).GetArea(MapAreaUI))
}

func (s *Ship) Update() {
}
