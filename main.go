package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/chmllr/imgtb/checksum"
	"github.com/chmllr/imgtb/health"
	"github.com/chmllr/imgtb/imp"

	"github.com/fatih/color"
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
		log.Println("importing to", *lib, "from", *source, "...")
		imp.Import(*lib, *source)
	case "seal":
		log.Println("computing checksums in", *lib, "...")
		hashes, err := checksum.Report(*lib)
		if err != nil {
			log.Fatal(err)
		}
		saveReport(*lib, hashes)
	case "health":
		log.Printf("checking health of %q...\n", *lib)
		hashes, err := checksum.Report(*lib)
		if err != nil {
			log.Fatal(err)
		}
		mapping1 := map[string]string{}
		for _, e := range hashes {
			mapping1[e.Path] = e.Hash
		}
		mapping2 := map[string]string{}
		content, err := ioutil.ReadFile(filepath.Join(*lib, "checksums.txt"))
		if err != nil {
			log.Fatal(err)
		}
		for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
			fields := strings.Split(line, "::")
			if len(fields) != 2 {
				log.Fatalf("unexpected line: %q", line)
			}
			mapping2[fields[0]] = fields[1]
		}
		corrupted := health.Verify(*lib, mapping1, mapping2)
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range corrupted {
			color.Red("File %s is corrupted!", path)
		}
		for path := range mapping2 {
			color.Red("File %s is missing!", path)
		}
		for path := range mapping1 {
			color.Cyan("File %s is new!", path)
		}
		if len(corrupted) == 0 && len(mapping1) > 0 && len(mapping2) == 0 {
			fmt.Println("Do you want me do seal new files? [n/Y]")
			var answer string
			fmt.Scanf("%s", &answer)
			if answer == "" || answer == "y" || answer == "Y" {
				log.Println("sealing...")
				for path, hash := range mapping1 {
					hashes = append(hashes, struct{ Path, Hash string }{path, hash})
					delete(mapping1, path)
				}
				saveReport(*lib, hashes)
			}
		}
		if len(corrupted) == 0 && len(mapping1) == 0 && len(mapping2) == 0 {
			log.Println(*lib, "is in perfect health! ✅")
		}
	default:
		log.Printf("Error: unknown command %q\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Usage: imgtb [OPTIONS] <COMMAND>")
	fmt.Println("Avaliable options:")
	flag.PrintDefaults()
}

func saveReport(lib string, hashes []struct{ Path, Hash string }) {
	filepath := filepath.Join(lib, "checksums.txt")
	var buf bytes.Buffer
	for _, e := range hashes {
		buf.WriteString(e.Path)
		buf.WriteString("::")
		buf.WriteString(e.Hash)
		buf.WriteString("\n")
	}
	err := ioutil.WriteFile(filepath, buf.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("checksums written to", filepath)
}
