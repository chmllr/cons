package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chmllr/imgtb/health"
	"github.com/chmllr/imgtb/imp"
	"github.com/chmllr/imgtb/seal"

	"github.com/fatih/color"
)

func main() {
	lib := flag.String("lib", "", "path to the photo library")
	source := flag.String("source", "", "source directory")
	deep := flag.Bool("deep", false, "deep check (includes md5 comparison)")
	flag.Parse()

	if len(flag.Args()) != 1 || *lib == "" {
		printHelp()
		os.Exit(1)
	}

	cmd := flag.Args()[0]

	switch cmd {
	case "import":
		if *source == "" {
			log.Fatal("no source folder specified")
		}
		log.Println("importing to", *lib, "from", *source, "...")
		refs, err := imp.Import(*lib, *source)
		if err != nil {
			log.Fatalf("couldn't import: %v", err)
		}
		sealed, err := seal.Registry(*lib)
		if err != nil {
			log.Fatalf("couldn't get registry: %v", err)
		}
		for _, ref := range sealed {
			refs = append(refs, ref)
		}
		saveReport(*lib, refs)
	case "repair":
		log.Println("sealing", *lib, "...")
		files, err := seal.Report(*lib, true)
		if err != nil {
			log.Fatalf("couldn't get report: %v", err)
		}
		saveReport(*lib, files)
	case "health":
		log.Printf("checking (deep: %t) health of %q...\n", *deep, *lib)
		libRefs, err := seal.Report(*lib, *deep)
		if err != nil {
			log.Fatalf("couldn't get report: %v", err)
		}
		corrupted, found, sealed, duplicates, err := health.Verify(*lib, *deep, libRefs)
		if err != nil {
			log.Fatalf("couldn't verify: %v", err)
		}
		for _, path := range corrupted {
			color.Red("File %s is corrupted!", path)
		}
		for path := range sealed {
			color.Red("File %s is missing!", path)
		}
		for _, paths := range duplicates {
			color.Yellow("These files are duplicates:")
			for _, path := range paths {
				color.Yellow(" - %s", path)
			}
		}
		for path := range found {
			color.Cyan("File %s is new!", path)
		}
		if len(corrupted) == 0 && len(sealed) == 0 && len(duplicates) == 0 && len(found) == 0 {
			if *deep {
				log.Printf("%q is in perfect health! ✅\n", *lib)
			} else {
				log.Printf("%q is in a good health (use --deep for a complete check)! ✅\n", *lib)
			}
		}
	default:
		log.Printf("Error: unknown command %q\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: imgtb --lib <PATH> [OPTIONS] <COMMAND>

Avaliable commands:

import (requires option --source <PATH>):
	Imports all media files from the specified source path into the lib folder.
	It creates the corresponding folder structure (<lib>/YYYY/MM/DD) if necessary.

repair:
	Seals all existing files with their md5 hashes into the lib index. It does not
	make any mutating operations on the library!

health (accepts option --deep):
	Checks existing file structure against the index. This command can detect 
	missing, modified, duplicated and new files. If option 'deep' is proveded, 
	checks the file hash as well.`)

}

func saveReport(lib string, refs []seal.LibRef) {
	filepath := filepath.Join(lib, "index.csv")
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	for _, e := range refs {
		if err := w.Write(e.Record()); err != nil {
			log.Fatalf("couldn't write csv record: %v", err)
		}
	}
	w.Flush()
	err := ioutil.WriteFile(filepath, buf.Bytes(), 0666)
	if err != nil {
		log.Fatalf("couldn't write index: %v", err)
	}
	log.Println("checksums written to", filepath)
}
