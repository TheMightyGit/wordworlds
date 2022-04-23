package cartridge

import (
	"image"
	"math/rand"
	"strings"
	"time"

	"github.com/TheMightyGit/marv/marvtypes"
	"github.com/TheMightyGit/wordworlds/dictionary"
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
	for rowIdx, letterRow := range []string{"--------", "--------", "--------"} {
		for letterIdx := range letterRow {
			hitBox := image.Rectangle{
				image.Point{30 * letterIdx, 30 * rowIdx},
				image.Point{(30 * letterIdx) + 30, (30 * rowIdx) + 30},
			}
			hitBox = hitBox.Add(spriteButtonLettersOffset)
			hitBox = hitBox.Sub(image.Point{5, 5})

			buttonArea := s.buttonLettersArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{2, 2})})
			buttonBgArea := uiArea.GetSubArea(image.Rectangle{pos.Add(image.Point{4, 11}), pos.Add(image.Point{4 + 3, 11 + 3})})

			darkBox := darkRedBox
			brightBox := brightRedBox

			if letterIdx >= 5 {
				darkBox = darkOrangeBox
				brightBox = brightOrangeBox
			} else if letterIdx >= 3 {
				darkBox = darkGreenBox
				brightBox = brightGreenBox
			}

			button := NewLetterButton(buttonArea, buttonBgArea, randomLetter(), hitBox, func(lb *LetterButton) {
				api.ConsolePrintln(string(lb.letter))
			}, darkBox, brightBox)
			s.updateables = append(s.updateables, button)
			s.clickables = append(s.clickables, button)
			pos = pos.Add(image.Point{3, 0})
		}
		pos = pos.Add(image.Point{-(3 * len(letterRow)), 3})
	}

	s.spriteGuessWord = api.SpritesGet(SpriteGuessWord)
	s.spriteGuessWord.ChangePos(image.Rectangle{image.Point{0, 78}, image.Point{320, 30}})
	s.guessWordArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{32, 2})
	s.spriteGuessWord.Show(GfxBankGfx, s.guessWordArea)

	var timer *time.Timer
	timer = time.AfterFunc(
		time.Second*2,
		func() {
			s.guessWordArea.Clear(0, 0)
			pos = image.Point{0, 0}
			word := dictionary.Dictionary.RandomWord()
			pad := (16 - len(word)) / 2
			if pad < 0 {
				pad = 0
			}
			word = strings.Repeat(" ", pad) + word
			drawText(s.guessWordArea, pos, word, 2)
			timer.Reset(time.Second * 2)
		},
	)

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

var (
	validLetters = []rune("AABCDEEFGHIIJKLMNOOPQRSTUUVWXYYZ")
)

func randomLetter() rune {
	return validLetters[rand.Intn(len(validLetters))]
}

type LetterButton struct {
	area       marvtypes.MapBankArea
	bgArea     marvtypes.MapBankArea
	letter     rune
	hitBox     image.Rectangle
	clicked    bool
	onClick    LetterButtonOnClickFunc
	disabled   bool
	enabledBg  *[9]image.Point
	disabledBg *[9]image.Point
}

type LetterButtonOnClickFunc func(*LetterButton)

func NewLetterButton(
	area marvtypes.MapBankArea,
	bgArea marvtypes.MapBankArea,
	letter rune,
	hitBox image.Rectangle,
	onClick LetterButtonOnClickFunc,
	enabledBg *[9]image.Point,
	disabledBg *[9]image.Point,
) *LetterButton {
	b := &LetterButton{
		area:       area,
		bgArea:     bgArea,
		letter:     letter,
		hitBox:     hitBox,
		onClick:    onClick,
		enabledBg:  enabledBg,
		disabledBg: disabledBg,
	}
	b.Enable()
	return b
}

func (b *LetterButton) Start() {
	drawRune(b.area, image.Point{}, b.letter)
}

func (b *LetterButton) Update() {
	if b.clicked {
		// b.letter = 'X'
		// drawRune(b.area, image.Point{}, b.letter)
		b.Disable()
		b.onClick(b)

		b.clicked = false
	}
}

func (b *LetterButton) OnClick(pos image.Point) bool {
	// api.ConsolePrintln(b.letter, b.hitBox, pos)
	if pos.In(b.hitBox) {
		if !b.disabled {
			b.clicked = true
		}
		return true
	}
	return false
}

var (
	brightRedBox = &[9]image.Point{
		{0, 11}, {1, 11}, {2, 11},
		{0, 12}, {1, 12}, {2, 12},
		{0, 13}, {1, 13}, {2, 13},
	}
	darkRedBox = &[9]image.Point{
		{0, 1}, {1, 1}, {2, 1},
		{0, 2}, {1, 2}, {2, 2},
		{0, 3}, {1, 3}, {2, 3},
	}
	brightGreenBox = &[9]image.Point{
		{3, 11}, {4, 11}, {5, 11},
		{3, 12}, {4, 12}, {5, 12},
		{3, 13}, {4, 13}, {5, 13},
	}
	darkGreenBox = &[9]image.Point{
		{3, 1}, {4, 1}, {5, 1},
		{3, 2}, {4, 2}, {5, 2},
		{3, 3}, {4, 3}, {5, 3},
	}
	brightOrangeBox = &[9]image.Point{
		{6, 11}, {7, 11}, {8, 11},
		{6, 12}, {7, 12}, {8, 12},
		{6, 13}, {7, 13}, {8, 13},
	}
	darkOrangeBox = &[9]image.Point{
		{6, 1}, {7, 1}, {8, 1},
		{6, 2}, {7, 2}, {8, 2},
		{6, 3}, {7, 3}, {8, 3},
	}
)

func (b *LetterButton) Enable() {
	b.bgArea.DrawBox(
		image.Rectangle{image.Point{}, image.Point{2, 2}},
		b.enabledBg,
		0, 0)
	b.disabled = false
}

func (b *LetterButton) Disable() {
	b.bgArea.DrawBox(
		image.Rectangle{image.Point{}, image.Point{2, 2}},
		b.disabledBg,
		0, 0)
	b.disabled = true
}
