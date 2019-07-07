package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/chmllr/cons/health"
	"github.com/chmllr/cons/index"

	"github.com/fatih/color"
)

func main() {
	dir := flag.String("dir", "", "path to the directory")
	deep := flag.Bool("deep", false, "deep check (includes md5 comparison)")
	flag.Parse()

	if len(flag.Args()) != 1 {
		printHelp()
		return
	}

	if *dir == "" {
		var err error
		*dir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	ignoreFile := path.Join(*dir, ".consignore")
	data, err := ioutil.ReadFile(ignoreFile)
	filters := []*regexp.Regexp{
		regexp.MustCompile(`\.index.csv`),
		regexp.MustCompile(`\.consignore`),
	}
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to open %s: %v", ignoreFile, err)
	} else {
		for _, v := range strings.Split(string(data), "\n") {
			if strings.TrimSpace(v) == "" {
				continue
			}
			re, err := regexp.Compile(v)
			if err != nil {
				log.Printf("parsing %s: %v", ignoreFile, err)
			}
			filters = append(filters, re)
		}
	}

	cmd := flag.Args()[0]

	switch cmd {
	case "seal":
		fmt.Printf("sealing %q...\n", *dir)
		files, err := index.Report(*dir, filters, true)
		if err != nil {
			log.Fatalf("couldn't get report: %v", err)
		}
		index.Save(*dir, files)
	case "verify":
		fmt.Printf("verifying (deep: %t) %q...\n", *deep, *dir)
		libRefs, err := index.Report(*dir, filters, *deep)
		if err != nil {
			log.Fatalf("couldn't get report: %v", err)
		}
		corrupted, found, sealed, duplicates, err := health.Verify(*dir, *deep, libRefs)
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
				fmt.Printf("%q is sound! ✅\n", *dir)
			} else {
				fmt.Printf("%q looks sound (use --deep for a hash based check)! ✅\n", *dir)
			}
		}
	default:
		fmt.Printf("Error: unknown command %q\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`Usage: cons --dir <PATH> [OPTIONS] <COMMAND>

cons keeps track of file changes in a directory. If no directory parameter was
provided, the current directory is used. cons creates a file index inside 
".index.csv" file in CSV format. If certain files should be ignored, create
a file name ".consignore" inside the tracked directory with one regular
expression per line, matching the file names to be ignored.

Avaliable commands:

seal:
	Seals all existing files with their MD5 hashes into the directory index.
	It does not make any mutating operations on the directory!

verify (accepts option --deep):
	Checks existing file structure against the index. This command can detect 
	missing, modified, duplicated and new files. If option 'deep' is provided, 
	compares files based on the MD5 hash.`)

}
