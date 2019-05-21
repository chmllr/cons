package seal

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type LibRef struct {
	Path, Hash string
	Size       int64
}

func (r LibRef) Record() []string {
	return []string{r.Path, strconv.FormatInt(r.Size, 10), r.Hash}
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

func Registry(lib string) (map[string]LibRef, error) {
	sealed := map[string]LibRef{}
	content, err := ioutil.ReadFile(filepath.Join(lib, "index.csv"))
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(bytes.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		size, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}

		sealed[record[0]] = LibRef{record[0], record[2], size}
	}

	return sealed, nil
}

func hash(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return fmt.Sprintf("%x", md5.Sum(data)), err
}

func pad(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}
