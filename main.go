package main

import (
	"fmt"
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
			err := copyFile(from, to)
			if err != nil {
				log.Printf("couldn't copy file from %q to %q: %v\n", from, to, err)
				continue
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

func copyFile(from, to string) error {
	stat, err := os.Stat(to)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	destinationNonEmpty := stat != nil

	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	if destinationNonEmpty {
		destData, err := ioutil.ReadFile(to)
		if err != nil {
			return fmt.Errorf("file %q exists and couldn't be read: %v", to, err)
		}
		for i := range destData {
			if i >= len(data) || data[i] != destData[i] {
				return fmt.Errorf("file %q already exists and is different than file %q", to, from)
			}
		}
		return fmt.Errorf("file %q is already imported", from)
	}

	return ioutil.WriteFile(to, data, 0644)
}
