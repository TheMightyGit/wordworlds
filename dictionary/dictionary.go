package dictionary

import (
	_ "embed"
	"math/rand"
	"strings"
)

var (
	//go:embed "dictionary.txt"
	rawDictionary string
	Dictionary    = newDict(rawDictionary)
)

type Dict interface {
	Words() []string
	RandomWord() string
	ContainsWord(string) bool
	RandomLetterReducedByLettersInPlayHistogram(map[rune]float64) rune
	GetLetterFrequency(rune) int
}

type dict struct {
	words   []string
	wordMap map[string]struct{}

	lettersFreq []*letterFreq
}

type letterFreq struct {
	letter   rune
	freq     int
	freqNorm float64 // 0.0 to 1.0
	rank     int
}

func newDict(rawfiledata string) Dict {
	words := strings.Split(rawfiledata, "\n")

	wordMap := map[string]struct{}{}
	for i, w := range words {
		words[i] = strings.ToUpper(w)
		wordMap[words[i]] = struct{}{}
	}

	d := &dict{
		words:   words,
		wordMap: wordMap,
	}

	d.calcLetterFrequency()

	return d
}

func (d *dict) Words() []string {
	return d.words
}

func (d *dict) RandomWord() string {
	return d.words[rand.Intn(len(d.words))]
}

func (d *dict) ContainsWord(word string) bool {
	_, found := d.wordMap[word]
	return found
}

func (d *dict) RandomLetterReducedByLettersInPlayHistogram(lettersInPlayHistogram map[rune]float64) rune {
	// subtract the latters in play histogram from the dictionary histogram
	histogram := map[rune]float64{}
	histTotal := float64(0)
	for _, lf := range d.lettersFreq {
		// val := lf.freqNorm - (lettersInPlayHistogram[lf.letter] * 0.9)
		val := lf.freqNorm * (1.0 - (lettersInPlayHistogram[lf.letter] * 0.95))
		if val < 0 {
			val = 0
		}
		histTotal += val
		histogram[lf.letter] = val
	}
	// fmt.Println(histogram)

	pos := rand.Float64() * histTotal // pick a random position across whole range
	// walk letters until we find the letter that pos is in.
	total := float64(0)
	for k, v := range histogram {
		total += v
		if total > pos {
			return k
		}
	}
	return '-'
}

func (d *dict) GetLetterFrequency(letter rune) int {
	for _, lf := range d.lettersFreq {
		if lf.letter == letter {
			return lf.rank
		}
	}
	return 0
}

func (d *dict) calcLetterCount() map[rune]int {
	countMap := map[rune]int{}
	for _, word := range d.words {
		for _, letter := range word {
			countMap[letter]++
		}
	}
	return countMap
}

func (d *dict) calcLetterFrequency() {
	highestFreq := 0
	for k, v := range d.calcLetterCount() {
		d.lettersFreq = append(d.lettersFreq, &letterFreq{
			letter: k,
			freq:   v,
		})
		if v > highestFreq {
			highestFreq = v
		}
	}
	// calculate bucketed rank and normalised frequency
	numBuckets := float64(8)
	for _, lf := range d.lettersFreq {
		lf.rank = int((numBuckets + 1) - ((float64(lf.freq) / float64(highestFreq)) * numBuckets))
		lf.freqNorm = float64(lf.freq) / float64(highestFreq)
		// fmt.Println(lf)
	}
}
