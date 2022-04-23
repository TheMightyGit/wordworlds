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
}

func newDict(rawfiledata string) Dict {
	words := strings.Split(rawfiledata, "\n")
	for i, w := range words {
		words[i] = strings.ToUpper(w)
	}
	return &dict{
		words: words,
	}
}

type dict struct {
	words []string
}

func (d *dict) Words() []string {
	return d.words
}

func (d *dict) RandomWord() string {
	return d.words[rand.Intn(len(d.words))]
}
