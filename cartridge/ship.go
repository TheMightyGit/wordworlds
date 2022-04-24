package cartridge

import (
	"image"
	"strings"

	"github.com/TheMightyGit/marv/marvtypes"
	"github.com/TheMightyGit/wordworlds/dictionary"
)

type Ship struct {
	spriteShip   marvtypes.Sprite
	areaShip     marvtypes.MapBankArea
	spritePanels marvtypes.Sprite
	areaPanels   marvtypes.MapBankArea
	spriteUI     marvtypes.Sprite

	spriteButtonIcons marvtypes.Sprite
	spriteGuessWord   marvtypes.Sprite

	buttonIconsArea marvtypes.MapBankArea
	guessWordArea   marvtypes.MapBankArea
	uiArea          marvtypes.MapBankArea

	clickables  []Clickable
	updateables []Updateable

	allLetterButtons      []*LetterButton
	selectedLetterButtons []*LetterButton

	okButton      *LetterButton
	delButton     *LetterButton
	shuffleButton *LetterButton
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
	s.uiArea = api.MapBanksGet(MapBankGfx).GetArea(MapAreaUI)
	s.spriteUI.Show(GfxBankGfx, s.uiArea)

	s.buttonIconsArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{32, 20}) // full screen

	s.spriteButtonIcons = api.SpritesGet(SpriteButtonLetters)
	s.spriteButtonIcons.ChangePos(image.Rectangle{image.Point{5, 5}, image.Point{320 - 5, 200 - 5}})
	s.spriteButtonIcons.Show(GfxBankGfx, s.buttonIconsArea)
	spriteButtonLettersOffset := image.Point{40, 110}

	pos := image.Point{4, 11}
	for rowIdx, letterRow := range []string{"--------", "--------", "--------"} {
		for letterIdx := range letterRow {
			hitBox := image.Rectangle{
				image.Point{30 * letterIdx, 30 * rowIdx},
				image.Point{(30 * letterIdx) + 30, (30 * rowIdx) + 30},
			}
			hitBox = hitBox.Add(spriteButtonLettersOffset)

			buttonArea := s.buttonIconsArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{2, 2})})
			buttonBgArea := s.uiArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{3, 3})})

			darkBox := darkRedBox
			brightBox := brightRedBox

			if letterIdx >= 5 {
				darkBox = darkOrangeBox
				brightBox = brightOrangeBox
			} else if letterIdx >= 3 {
				darkBox = darkGreenBox
				brightBox = brightGreenBox
			}

			button := NewLetterButton(buttonArea, buttonBgArea, dictionary.Dictionary.RandomLetter(), hitBox, func(lb *LetterButton) {
				if lb.disabled {
					// remove from letter list
					newSelectedLetterButtons := []*LetterButton{}
					for _, b := range s.selectedLetterButtons {
						if b != lb {
							newSelectedLetterButtons = append(newSelectedLetterButtons, b)
						}
					}
					s.selectedLetterButtons = newSelectedLetterButtons
					lb.Enable()
				} else {
					// add letter and disable button
					lb.Disable()
					s.selectedLetterButtons = append(s.selectedLetterButtons, lb)
				}
				s.updateGuessWord()
			}, darkBox, brightBox)
			s.updateables = append(s.updateables, button)
			s.clickables = append(s.clickables, button)
			s.allLetterButtons = append(s.allLetterButtons, button)
			pos = pos.Add(image.Point{3, 0})
		}
		pos = pos.Add(image.Point{-(3 * len(letterRow)), 3})
	}

	s.okButton = s.addOkButton()
	s.delButton = s.addDelButton()
	s.shuffleButton = s.addShuffleButton()

	s.spriteGuessWord = api.SpritesGet(SpriteGuessWord)
	s.spriteGuessWord.ChangePos(image.Rectangle{image.Point{0, 78}, image.Point{320, 30}})
	s.guessWordArea = api.MapBanksGet(MapBankGfx).AllocArea(image.Point{32, 2})
	s.spriteGuessWord.Show(GfxBankGfx, s.guessWordArea)

	for _, updateable := range s.updateables {
		updateable.Start()
	}

	s.updateGuessWord()
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

func (s *Ship) getGuessWord() string {
	word := ""
	for _, b := range s.selectedLetterButtons {
		word += string(b.letter)
	}
	return word
}

func (s *Ship) updateGuessWord() {
	word := s.getGuessWord()
	pad := (16 - len(word)) / 2
	if pad < 0 {
		pad = 0
	}
	paddedWord := strings.Repeat(" ", pad) + word

	s.guessWordArea.Clear(0, 0)
	if len(word) > 0 {
		pos := image.Point{0, 0}
		drawText(s.guessWordArea, pos, paddedWord, 2)

		if dictionary.Dictionary.ContainsWord(word) {
			// api.ConsolePrintln(word, " VALID")
			s.okButton.letter = 'o'
		} else {
			// api.ConsolePrintln(word, " INVALID")
			s.okButton.letter = ' '
		}

		s.delButton.letter = 'd'
	} else {
		s.okButton.letter = ' '
		s.delButton.letter = ' '
	}

	s.okButton.Start()
	s.delButton.Start()
}

func (s *Ship) addOkButton() *LetterButton {
	return s.addActionButton(
		image.Point{29, 16},
		'o',
		func(lb *LetterButton) {
			// new letters for used buttons
			for _, lb := range s.selectedLetterButtons {
				lb.Shuffle()
			}
			s.removeGuessWord()
		},
	)
}

func (s *Ship) addDelButton() *LetterButton {
	return s.addActionButton(
		image.Point{29, 12},
		'd',
		func(lb *LetterButton) {
			siz := len(s.selectedLetterButtons)
			if siz > 0 {
				s.selectedLetterButtons[siz-1].Enable()
				s.selectedLetterButtons = s.selectedLetterButtons[:siz-1]
				s.updateGuessWord()
			}
		},
	)
}

func (s *Ship) addShuffleButton() *LetterButton {
	return s.addActionButton(
		image.Point{0, 12},
		's',
		func(lb *LetterButton) {
			s.shuffleButtons()
		},
	)
}

func (s *Ship) removeGuessWord() {
	// re-enable associated buttons
	for _, lb := range s.selectedLetterButtons {
		lb.Enable()
	}
	// remove selected buttons
	s.selectedLetterButtons = []*LetterButton{}
	s.updateGuessWord()
}

func (s *Ship) shuffleButtons() {
	// only shuffle un-selected buttons
	for _, lb := range s.allLetterButtons {
		if !lb.disabled {
			lb.Shuffle()
		}
	}
}

func (s *Ship) addActionButton(pos image.Point, icon rune, onClick func(*LetterButton)) *LetterButton {
	hitPos := image.Point{pos.X * 10, pos.Y * 10}
	hitBox := image.Rectangle{hitPos, hitPos.Add(image.Point{30, 30})}
	buttonArea := s.buttonIconsArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{2, 2})})
	buttonBgArea := s.uiArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{3, 3})})
	button := NewLetterButton(buttonArea, buttonBgArea, icon, hitBox, onClick, solidBlueBox, solidBlueBox)
	s.updateables = append(s.updateables, button)
	s.clickables = append(s.clickables, button)
	return button
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

	shuffleAnimCountdown int
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
	if b.shuffleAnimCountdown > 0 {
		if b.shuffleAnimCountdown%5 == 0 {
			b.letter = dictionary.Dictionary.RandomLetter()
			b.Start()
		}
		b.shuffleAnimCountdown--
		if b.shuffleAnimCountdown == 0 {
			b.Enable()
		}
	} else if b.clicked {
		b.onClick(b)
	}
	b.clicked = false
}

func (b *LetterButton) OnClick(pos image.Point) bool {
	// api.ConsolePrintln(b.letter, b.hitBox, pos)
	if pos.In(b.hitBox) {
		b.clicked = true
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
	solidBlueBox = &[9]image.Point{
		{1, 8}, {1, 8}, {1, 8},
		{1, 8}, {1, 8}, {1, 8},
		{1, 8}, {1, 8}, {1, 8},
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

func (b *LetterButton) Shuffle() {
	b.shuffleAnimCountdown = 60 * 1 // 1 seconds at 60fps
	b.Disable()
}
