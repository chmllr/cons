package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {
	sourceFolder := "./import"
	files, err := ioutil.ReadDir(sourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("files found:", len(files))
	folders := map[string][]string{}
	for _, f := range files {
		t, err := datetime(filepath.Join(sourceFolder, f.Name()))
		if err != nil {
			log.Printf("skipping file %s due to errors: %v\n", f.Name(), err)
			continue
		}
		folder := t.Format("2006/01/02")
		folders[folder] = append(folders[folder], f.Name())
	}
	log.Println("new folders required:", len(folders))
	libFolder := "./testlib"
	imported := 0
	for folder, files := range folders {
		destinationPath := filepath.Join(libFolder, folder)
		if err := os.MkdirAll(destinationPath, os.ModePerm); err != nil && !os.IsExist(err) {
			log.Printf("couldn't create folder %q: %v\n", destinationPath, err)
			continue
		}
		for _, fileName := range files {
			from := filepath.Join(sourceFolder, fileName)
			to := filepath.Join(destinationPath, fileName)
			err := os.Rename(from, to)
			if err != nil {
				log.Printf("couldn't move file from %q to %q: %v\n", from, to, err)
			}
			imported++
		}
	}
	log.Printf("files succesfully imported: %d/%d", imported, len(files))
}

func datetime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Now(), err
	}

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return time.Now(), err
	}

	return x.DateTime()
}
