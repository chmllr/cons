package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chmllr/imgtb/checksum"
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
	case "checksum":
		log.Println("computing checksums in", *lib)
		content, err := checksum.Report(*lib)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(content)
	default:
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("TBD")
}
