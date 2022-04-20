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

	spriteButtonLetters marvtypes.Sprite
	spriteGuessWord     marvtypes.Sprite

	buttonLettersArea marvtypes.MapBankArea
	guessWordArea     marvtypes.MapBankArea

	clickables  []Clickable
	updateables []Updateable
}

func NewShip(
	spriteShip marvtypes.Sprite,
	areaShip marvtypes.MapBankArea,
	spritePanels marvtypes.Sprite,
	areaPanels marvtypes.MapBankArea,
	spriteUI marvtypes.Sprite,
) *Ship {
	return &Ship{
		spriteShip:   spriteShip,
		areaShip:     areaShip,
		spritePanels: spritePanels,
		areaPanels:   areaPanels,
		spriteUI:     spriteUI,
	}
}

func (s *Ship) Start() {
	s.spriteShip = api.SpritesGet(SpriteShip)
	s.spriteShip.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	s.spriteShip.Show(GfxBankGfx, s.areaShip)

	s.spritePanels = api.SpritesGet(SpritePanels)
	s.spritePanels.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	s.spritePanels.Show(GfxBankGfx, s.areaPanels)

	s.spriteUI = api.SpritesGet(SpriteUI)
	s.spriteUI.ChangePos(image.Rectangle{image.Point{0, 0}, image.Point{320, 200}})
	uiArea := api.MapBanksGet(MapBankGfx).GetArea(MapAreaUI)
	s.spriteUI.Show(GfxBankGfx, uiArea)

	spriteButtonLettersOffset := image.Point{5 + (4 * 10), 5 + (11 * 10)}
	s.spriteButtonLetters = api.SpritesGet(SpriteButtonLetters)
	s.spriteButtonLetters.ChangePos(image.Rectangle{spriteButtonLettersOffset, image.Point{8 * 30, 30 * 3}})
	s.buttonLettersArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{8 * 3, 3 * 3})
	s.spriteButtonLetters.Show(GfxBankGfx, s.buttonLettersArea)

	pos := image.Point{0, 0}
	for rowIdx, letterRow := range []string{"ABCDEFGH", "IJKLMNOP", "QRSTUVWX"} {
		for letterIdx, letter := range letterRow {
			hitBox := image.Rectangle{
				image.Point{30 * letterIdx, 30 * rowIdx},
				image.Point{(30 * letterIdx) + 30, (30 * rowIdx) + 30},
			}
			hitBox = hitBox.Add(spriteButtonLettersOffset)
			hitBox = hitBox.Sub(image.Point{5, 5})

			buttonArea := s.buttonLettersArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{2, 2})})
			buttonBgArea := uiArea.GetSubArea(image.Rectangle{pos.Add(image.Point{4, 11}), pos.Add(image.Point{4 + 3, 11 + 3})})

			button := NewLetterButton(buttonArea, buttonBgArea, letter, hitBox)
			s.updateables = append(s.updateables, button)
			s.clickables = append(s.clickables, button)
			pos = pos.Add(image.Point{3, 0})
		}
		pos = pos.Add(image.Point{-(3 * len(letterRow)), 3})
	}

	s.spriteGuessWord = api.SpritesGet(SpriteGuessWord)
	s.spriteGuessWord.ChangePos(image.Rectangle{image.Point{80, 80}, image.Point{30 * 16, 30}})
	s.guessWordArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{16 * 2, 2})
	s.spriteGuessWord.Show(GfxBankGfx, s.guessWordArea)

	pos = image.Point{0, 0}
	drawText(s.guessWordArea, pos, "SOMEWORD", 2)

	for _, updateable := range s.updateables {
		updateable.Start()
	}
}

func (s *Ship) Update() {
	for _, updateable := range s.updateables {
		updateable.Update()
	}
}

func (s *Ship) OnClick(pos image.Point) bool {
	for _, clickable := range s.clickables {
		if clickable.OnClick(pos) {
			return true
		}
	}
	return false
}

type LetterButton struct {
	area    marvtypes.MapBankArea
	bgArea  marvtypes.MapBankArea
	letter  rune
	hitBox  image.Rectangle
	clicked bool
}

func NewLetterButton(
	area marvtypes.MapBankArea,
	bgArea marvtypes.MapBankArea,
	letter rune,
	hitBox image.Rectangle,
) *LetterButton {
	return &LetterButton{
		area:   area,
		bgArea: bgArea,
		letter: letter,
		hitBox: hitBox,
	}
}

func (b *LetterButton) Start() {
	drawRune(b.area, image.Point{}, b.letter)
}

func (b *LetterButton) Update() {
	if b.clicked {
		// b.letter = 'X'
		// drawRune(b.area, image.Point{}, b.letter)
		b.bgArea.DrawBox(
			image.Rectangle{image.Point{}, image.Point{2, 2}},
			&[9]image.Point{
				{0, 11},
				{1, 11},
				{2, 11},
				{0, 12},
				{1, 12},
				{2, 12},
				{0, 13},
				{1, 13},
				{2, 13},
			},
			0, 0)
		b.clicked = false
	}
}

func (b *LetterButton) OnClick(pos image.Point) bool {
	// api.ConsolePrintln(b.letter, b.hitBox, pos)
	if pos.In(b.hitBox) {
		b.clicked = true
		return true
	}
	return false
}
