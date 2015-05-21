package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"

	"github.com/inkyblackness/res/chunk/dos"

	"github.com/inkyblackness/construct/chunks"
)

const (
	// Version contains the current version number
	Version = "0.1.0"
	// Name is the name of the application
	Name = "InkyBlackness Construct"
	// Title contains a combined string of name and version
	Title = Name + " v." + Version
)

func main() {
	arguments, _ := docopt.Parse(usage(), nil, true, Title, false)
	outFileName := arguments["--file"].(string)
	writer, errOut := os.Create(outFileName)
	if errOut != nil {
		fmt.Printf("Error creating destination: %v\n", errOut)
	}

	chunkConsumer := dos.NewChunkConsumer(writer)

	chunks.AddArchiveName(chunkConsumer, "Starting Game")
	chunks.AddGameState(chunkConsumer)
	chunks.AddLevel(chunkConsumer, 1)

	chunkConsumer.Finish()
}

func usage() string {
	return Title + `

Usage:
  construct [--file=<file-name>]
	construct -h | --help
	construct --version

Options:
  --file=<file-name>  specifies the target file name. [default: archive.dat]
  -h --help           Show this screen.
  --version           Show version.
`
}
