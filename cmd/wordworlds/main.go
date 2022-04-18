package main

import (
	"github.com/TheMightyGit/marv/marvlib"
	"github.com/TheMightyGit/wordworlds/cartridge"
)

func main() {
	marvlib.API.ConsoleBoot(
		"wordworlds",
		cartridge.Resources,
		cartridge.Start,
		cartridge.Update,
	)
}
