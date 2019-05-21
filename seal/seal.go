package seal

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type LibRef struct {
	Path, Hash string
	Size       int64
}

var (
	jpegRegexp = regexp.MustCompile(`(?i)\.jpe?g`)
	mp4Regexp  = regexp.MustCompile(`(?i)\.mp4`)
)

// Reports walks through the folder structure and returns a mapping
// file path -> md5 hash
func Report(lib string, deep bool) (res []LibRef, err error) {
	maxLength := 0
	err = filepath.Walk(lib, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if info.IsDir() || !jpegRegexp.MatchString(path) && !mp4Regexp.MatchString(path) {
			return nil
		}

		var h string
		out := fmt.Sprintf("touching %s...", path)
		if maxLength < len(out) {
			maxLength = len(out)
		}
		fmt.Printf("\r%s", pad(out, maxLength))

		if deep {
			h, err = hash(path)
			if err != nil {
				return err
			}
		}
		res = append(res, LibRef{path, h, info.Size()})
		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("\r%s", "")
	log.Println(pad("sorting...", maxLength))
	sort.Slice(res, func(i, j int) bool { return res[i].Path < res[j].Path })
	return
}

func hash(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return fmt.Sprintf("%x", md5.Sum(data)), err
}

func pad(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}
