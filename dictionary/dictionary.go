package dictionary

import (
	_ "embed"
	"strings"
)

var (
	//go:embed "dictionary.txt"
	rawDictionary string
	Dictionary    = newDict(rawDictionary)
)

type Dict interface {
	Words() []string
}

func newDict(rawfiledata string) Dict {
	return &dict{
		rawDictionary: rawfiledata,
	}
}

type dict struct {
	rawDictionary string
}

func (d *dict) Words() []string {
	return strings.Split(d.rawDictionary, "\n")
}
