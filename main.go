package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
		hashes, err := checksum.Report(*lib)
		if err != nil {
			log.Fatal(err)
		}
		filepath := filepath.Join(*lib, "checksums.txt")
		var buf bytes.Buffer
		for _, e := range hashes {
			buf.WriteString(e.Path)
			buf.WriteString("::")
			buf.WriteString(e.Hash)
			buf.WriteString("\n")
		}
		err = ioutil.WriteFile(filepath, buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully written checksums to", filepath)
	case "health":
		fmt.Println("checking health of", *lib)
	default:
		fmt.Printf("Error: unknown command %q\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: imgtb [OPTIONS] <COMMAND>")
	fmt.Println("Avaliable options:")
	flag.PrintDefaults()
}
