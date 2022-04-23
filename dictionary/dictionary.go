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
	RandomLetter() rune
}

type dict struct {
	words   []string
	wordMap map[string]struct{}

	lettersFreq      []letterFreq
	lettersFreqTotal int
}

type letterFreq struct {
	letter rune
	freq   int
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

func (d *dict) RandomLetter() rune {
	// This tends towards the most common letters over time :/
	pos := rand.Intn(d.lettersFreqTotal)
	total := 0
	for _, lf := range d.lettersFreq {
		total += lf.freq
		if total > pos {
			return lf.letter
		}
	}
	return ' '
}

func (d *dict) calcLetterFrequency() {
	countMap := map[rune]int{}
	for _, word := range d.words {
		for _, letter := range word {
			countMap[letter]++
		}
	}
	for k, v := range countMap {
		d.lettersFreq = append(d.lettersFreq, letterFreq{
			letter: k,
			freq:   v,
		})
		d.lettersFreqTotal += v
	}
}
