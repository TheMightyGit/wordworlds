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
}

type dict struct {
	words   []string
	wordMap map[string]struct{}
}

func newDict(rawfiledata string) Dict {
	words := strings.Split(rawfiledata, "\n")

	wordMap := map[string]struct{}{}
	for i, w := range words {
		words[i] = strings.ToUpper(w)
		wordMap[words[i]] = struct{}{}
	}

	return &dict{
		words:   words,
		wordMap: wordMap,
	}
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
