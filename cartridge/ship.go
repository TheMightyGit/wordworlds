package cartridge

import (
	"image"
	"math"

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

	weaponProgressBar *ProgressBar
	hullProgressBar   *ProgressBar
	shieldProgressBar *ProgressBar

	overlay *Overlay

	getBaddies func() []*Baddie
}

func NewShip(
	spriteShip marvtypes.Sprite,
	areaShip marvtypes.MapBankArea,
	spritePanels marvtypes.Sprite,
	areaPanels marvtypes.MapBankArea,
	spriteUI marvtypes.Sprite,
	overlay *Overlay,
	getBaddies func() []*Baddie,
) *Ship {
	return &Ship{
		spriteShip:   spriteShip,
		areaShip:     areaShip,
		spritePanels: spritePanels,
		areaPanels:   areaPanels,
		spriteUI:     spriteUI,
		overlay:      overlay,
		getBaddies:   getBaddies,
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

	s.setupProgressBars()

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

			darkBox := darkBlueBox
			brightBox := brightRedBox
			letterType := LetterTypeWeapon

			if letterIdx >= 5 {
				// darkBox = darkOrangeBox
				brightBox = brightOrangeBox
				letterType = LetterTypeShield
			} else if letterIdx >= 3 {
				// darkBox = darkGreenBox
				brightBox = brightGreenBox
				letterType = LetterTypeHull
			}

			button := NewLetterButton(
				buttonArea,
				buttonBgArea,
				s.getRandomLetter(),
				hitBox,
				s.getRandomLetter,
				func(lb *LetterButton) {
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
						// api.ConsolePrintln(dictionary.Dictionary.GetLetterFrequency(lb.letter))
					}
					s.updateGuessWord()
				},
				darkBox,
				brightBox,
			)
			button.LetterType = letterType

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
	s.spriteGuessWord.ChangePalette(1)

	for _, updateable := range s.updateables {
		updateable.Start()
	}

	s.updateGuessWord()

	s.overlay.UpdateBaddies(s.getBaddies()...)
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

func (s *Ship) Damage(intAmount int) {
	amount := float64(intAmount) / 100.0

	// shield damage
	if amount > 0 {
		currentShield := s.shieldProgressBar.CurrentPercentage()
		shieldDamage := math.Min(currentShield, amount)
		if shieldDamage > 0 {
			s.shieldProgressBar.SetCurrentPercentage(currentShield - shieldDamage)
			s.shieldProgressBar.SetTargetPercentage(currentShield - shieldDamage)
		}
		amount -= shieldDamage
	}

	// hull damage
	if amount > 0 {
		currentHull := s.hullProgressBar.CurrentPercentage()
		hullDamage := math.Min(currentHull, amount)
		if hullDamage > 0 {
			s.hullProgressBar.SetCurrentPercentage(currentHull - hullDamage)
			s.hullProgressBar.SetTargetPercentage(currentHull - hullDamage)
		}
		amount -= hullDamage
	}

	if amount > 0 { // dead!
		api.ConsolePrintln("SHIP DEAD!")
	}
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

	s.guessWordArea.Clear(0, 0)
	if len(word) > 0 {
		pos := image.Point{0, 0}
		drawText(s.guessWordArea, pos, word, 2)
		margin := int((320.0 - (20.0 * float64(len(word)))) / 2.0)
		s.spriteGuessWord.ChangePos(image.Rectangle{image.Point{margin, 78}, image.Point{320, 20}})

		if dictionary.Dictionary.ContainsWord(word) {
			s.okButton.letter = 'o' // valid
			s.okButton.Enable()
		} else {
			s.okButton.letter = ' ' // not valid
			s.okButton.Disable()
		}

		s.delButton.letter = 'd'
	} else {
		s.okButton.letter = ' '
		s.delButton.letter = ' '
		s.okButton.Disable()
	}

	s.okButton.Start()
	s.delButton.Start()

	weaponLettersTotal := 0
	hullLettersTotal := 0
	shieldLettersTotal := 0
	for _, lb := range s.selectedLetterButtons {
		switch lb.LetterType {
		case LetterTypeWeapon:
			weaponLettersTotal += dictionary.Dictionary.GetLetterFrequency(lb.letter)
		case LetterTypeHull:
			hullLettersTotal += dictionary.Dictionary.GetLetterFrequency(lb.letter)
		case LetterTypeShield:
			shieldLettersTotal += dictionary.Dictionary.GetLetterFrequency(lb.letter)
		}
	}

	scoreMultiplier := 1.0 + (float64(len(s.selectedLetterButtons)-3) * 0.2) // 20% per letter past 3 letters
	// api.ConsolePrintln(s.weaponButtonsDown, s.hullButtonsDown, s.shieldButtonsDown, scoreMultiplier)
	weaponScore := (float64(weaponLettersTotal) * 0.01) * scoreMultiplier
	hullScore := (float64(hullLettersTotal) * 0.05) * scoreMultiplier
	shieldScore := (float64(shieldLettersTotal) * 0.01) * scoreMultiplier

	s.weaponProgressBar.SetTargetPercentage(s.weaponProgressBar.CurrentPercentage() + weaponScore)
	s.hullProgressBar.SetTargetPercentage(s.hullProgressBar.CurrentPercentage() + hullScore)
	s.shieldProgressBar.SetTargetPercentage(s.shieldProgressBar.CurrentPercentage() + shieldScore)
}

func (s *Ship) addOkButton() *LetterButton {
	return s.addActionButton(
		image.Point{29, 16},
		'o',
		func(lb *LetterButton) {
			if lb.disabled {
				return
			}

			// new letters for used buttons
			for _, lb := range s.selectedLetterButtons {
				lb.Shuffle()
			}
			s.weaponProgressBar.HitTarget()
			s.hullProgressBar.HitTarget()
			s.shieldProgressBar.HitTarget()

			if s.weaponProgressBar.CurrentPercentage() >= 0.5 {
				damage := 50

				if s.weaponProgressBar.CurrentPercentage() >= 0.80 {
					damage *= 2
				}

				baddies := s.getBaddies()
				for _, b := range baddies {
					if b.hitPoints > 0 {
						b.Damage(damage)
						break
					}
				}
				s.overlay.UpdateBaddies(baddies...)

				s.weaponProgressBar.SetCurrentPercentage(0.0)
				s.weaponProgressBar.SetTargetPercentage(0.0)
			}

			s.overlay.AddWord(s.getGuessWord())

			s.removeGuessWord()

			for _, b := range baddies {
				if b.hitPoints > 0 {
					b.TakeAction(s)
				}
			}
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

func (s *Ship) addBarText() {
	// HULL
	s.buttonIconsArea.Set(image.Point{12, 9}, 0, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{13, 9}, 1, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{14, 9}, 2, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{12, 10}, 0, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{13, 10}, 1, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{14, 10}, 2, 18, 0, 0)
	// WEAPON
	s.buttonIconsArea.Set(image.Point{3, 9}, 7, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{4, 9}, 8, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{5, 9}, 9, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{6, 9}, 10, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{7, 9}, 11, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{8, 9}, 12, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{9, 9}, 13, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{10, 9}, 14, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{11, 9}, 15, 17, 0, 0)

	s.buttonIconsArea.Set(image.Point{3, 10}, 7, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{4, 10}, 8, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{5, 10}, 9, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{6, 10}, 10, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{7, 10}, 11, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{8, 10}, 12, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{9, 10}, 13, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{10, 10}, 14, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{11, 10}, 15, 18, 0, 0)
	// SHIELD
	s.buttonIconsArea.Set(image.Point{18, 9}, 3, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{19, 9}, 4, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{20, 9}, 5, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{21, 9}, 6, 17, 0, 0)
	s.buttonIconsArea.Set(image.Point{18, 10}, 3, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{19, 10}, 4, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{20, 10}, 5, 18, 0, 0)
	s.buttonIconsArea.Set(image.Point{21, 10}, 6, 18, 0, 0)
}

func (s *Ship) addActionButton(pos image.Point, icon rune, onClick func(*LetterButton)) *LetterButton {
	hitPos := image.Point{pos.X * 10, pos.Y * 10}
	hitBox := image.Rectangle{hitPos, hitPos.Add(image.Point{30, 30})}
	buttonArea := s.buttonIconsArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{2, 2})})
	buttonBgArea := s.uiArea.GetSubArea(image.Rectangle{pos, pos.Add(image.Point{3, 3})})
	button := NewLetterButton(buttonArea, buttonBgArea, icon, hitBox, func() rune { return icon }, onClick, solidBlueBox, solidBlueBox)
	s.updateables = append(s.updateables, button)
	s.clickables = append(s.clickables, button)
	return button
}

func (s *Ship) setupProgressBars() {
	progressArea := api.MapBanksGet(MapBankGfx).GetArea(MapAreaBars)
	s.weaponProgressBar = NewProgressBar(
		image.Rectangle{image.Point{40 + 1, 100}, image.Point{90 - 2, 10}},
		api.SpritesGet(SpriteWeaponProgressBar),
		progressArea.GetSubArea(
			image.Rectangle{image.Point{0, 0}, image.Point{32, 1}},
		),
	)
	s.hullProgressBar = NewProgressBar(
		image.Rectangle{image.Point{130 + 1, 100}, image.Point{60 - 2, 10}},
		api.SpritesGet(SpriteHullProgressBar),
		progressArea.GetSubArea(
			image.Rectangle{image.Point{0, 1}, image.Point{32, 2}},
		),
	)
	s.shieldProgressBar = NewProgressBar(
		image.Rectangle{image.Point{190 + 1, 100}, image.Point{90 - 2, 10}},
		api.SpritesGet(SpriteShieldProgressBar),
		progressArea.GetSubArea(
			image.Rectangle{image.Point{0, 2}, image.Point{32, 3}},
		),
	)

	s.updateables = append(s.updateables, s.weaponProgressBar)
	s.updateables = append(s.updateables, s.hullProgressBar)
	s.updateables = append(s.updateables, s.shieldProgressBar)

	s.weaponProgressBar.SetCurrentPercentage(0.0)
	s.hullProgressBar.SetCurrentPercentage(1.0)
	s.shieldProgressBar.SetCurrentPercentage(0.5)

	s.addBarText()
}

func (s *Ship) getRandomLetter() rune {
	return dictionary.Dictionary.RandomLetterReducedByLettersInPlayHistogram(
		s.getLettersInPlayHistogram(),
	)
}

func (s *Ship) getLettersInPlayHistogram() map[rune]float64 {
	freqMap := map[rune]int{}

	for _, lb := range s.allLetterButtons {
		if !lb.disabled {
			freqMap[lb.letter]++
		}
	}

	highest := 0
	for _, v := range freqMap {
		if v > highest {
			highest = v
		}
	}

	normFreqMap := map[rune]float64{}
	for k, v := range freqMap {
		normFreqMap[k] = float64(v) / float64(highest)
	}

	// fmt.Println("getLettersInPlayHistogram: ", normFreqMap)

	return normFreqMap
}

const (
	LetterTypeNone = iota
	LetterTypeWeapon
	LetterTypeHull
	LetterTypeShield
)

type LetterButton struct {
	area       marvtypes.MapBankArea
	bgArea     marvtypes.MapBankArea
	letter     rune
	hitBox     image.Rectangle
	clicked    bool
	newLetter  func() rune
	onClick    LetterButtonOnClickFunc
	disabled   bool
	enabledBg  *[9]image.Point
	disabledBg *[9]image.Point

	LetterType int // LetterType*

	shuffleAnimCountdown int
}

type LetterButtonOnClickFunc func(*LetterButton)

func NewLetterButton(
	area marvtypes.MapBankArea,
	bgArea marvtypes.MapBankArea,
	letter rune,
	hitBox image.Rectangle,
	newLetter func() rune,
	onClick LetterButtonOnClickFunc,
	enabledBg *[9]image.Point,
	disabledBg *[9]image.Point,
) *LetterButton {
	b := &LetterButton{
		area:       area,
		bgArea:     bgArea,
		letter:     letter,
		hitBox:     hitBox,
		newLetter:  newLetter,
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
			b.letter = b.newLetter()
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
	darkBlueBox = &[9]image.Point{
		{0, 1}, {1, 1}, {2, 1},
		{0, 2}, {1, 2}, {2, 2},
		{0, 3}, {1, 3}, {2, 3},
	}
	brightRedBox = &[9]image.Point{
		{0, 11}, {1, 11}, {2, 11},
		{0, 12}, {1, 12}, {2, 12},
		{0, 13}, {1, 13}, {2, 13},
	}
	//	darkRedBox = &[9]image.Point{
	//		{0, 1}, {1, 1}, {2, 1},
	//		{0, 2}, {1, 2}, {2, 2},
	//		{0, 3}, {1, 3}, {2, 3},
	//	}
	brightGreenBox = &[9]image.Point{
		{3, 11}, {4, 11}, {5, 11},
		{3, 12}, {4, 12}, {5, 12},
		{3, 13}, {4, 13}, {5, 13},
	}
	//	darkGreenBox = &[9]image.Point{
	//		{3, 1}, {4, 1}, {5, 1},
	//		{3, 2}, {4, 2}, {5, 2},
	//		{3, 3}, {4, 3}, {5, 3},
	//	}
	brightOrangeBox = &[9]image.Point{
		{6, 11}, {7, 11}, {8, 11},
		{6, 12}, {7, 12}, {8, 12},
		{6, 13}, {7, 13}, {8, 13},
	}
	//	darkOrangeBox = &[9]image.Point{
	//		{6, 1}, {7, 1}, {8, 1},
	//		{6, 2}, {7, 2}, {8, 2},
	//		{6, 3}, {7, 3}, {8, 3},
	//	}
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
