package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chmllr/imgtb/imp"
)

func main() {

	lib := flag.String("lib", "", "path to the photo library")
	source := flag.String("source", "", "source directory")
	flag.Parse()

	if len(flag.Args()) != 1 {
		printHelp()
		os.Exit(1)
	}

	cmd := flag.Args()[0]

	switch cmd {
	case "import":
		log.Println("importing to", *lib, "from", *source)
		imp.Import(*lib, *source)
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("TBD")
}
