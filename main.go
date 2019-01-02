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
	path := "./testlib"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		t, err := datetime(filepath.Join(path, f.Name()))
		if err != nil {
			fmt.Printf("skipping file %s due to errors: %v", f.Name, err)
			continue
		}
		fmt.Println(f.Name(), t)
	}
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
