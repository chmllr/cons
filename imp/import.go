package imp

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func Import(libFolder, sourceFolder string) {
	files, err := ioutil.ReadDir(sourceFolder)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("files found:", len(files))
	folders := map[string][]string{}
	for _, f := range files {
		t, err := dateTime(filepath.Join(sourceFolder, f.Name()))
		if err != nil {
			log.Printf("skipping file %s due to errors: %v\n", f.Name(), err)
			continue
		}
		folder := t.Format("2006/01/02")
		folders[folder] = append(folders[folder], f.Name())
	}
	log.Println("new folders required:", len(folders))
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
			err := moveFile(from, to)
			if err != nil {
				log.Printf("couldn't move file from %q to %q: %v\n", from, to, err)
				continue
			}
			imported++
		}
	}
	log.Printf("files succesfully imported: %d/%d", imported, len(files))
}

func dateTime(path string) (time.Time, error) {
	jpeg := regexp.MustCompile(`(?i)\.jpe?g`)
	if jpeg.MatchString(path) {
		return imgDateTime(path)
	}
	mp4 := regexp.MustCompile(`(?i)\.mp4`)
	if mp4.MatchString(path) {
		return mp4DateTime(path)
	}
	return time.Now(), fmt.Errorf("unsupported file format: %s", path)
}

func mp4DateTime(path string) (time.Time, error) {
	_, fileName := filepath.Split(path)
	fNameParts := strings.Split(fileName, "_")
	if len(fNameParts) != 3 {
		return time.Time{}, fmt.Errorf("unexpected filename %q", fileName)
	}
	return time.Parse("20060102", fNameParts[1])
}

func imgDateTime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	return x.DateTime()
}

func moveFile(from, to string) error {
	stat, err := os.Stat(to)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	destinationNonEmpty := stat != nil

	if destinationNonEmpty {
		data, err := ioutil.ReadFile(from)
		if err != nil {
			return err
		}

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

	return os.Rename(from, to)
}
