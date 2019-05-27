package index

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
	"strconv"
	"strings"
)

type LibRef struct {
	Path, Hash string
	Size       int64
}

func NewLibRef(lib, path string, size int64) (LibRef, error) {
	ref := newLibRefUnhashed(lib, path, size)
	var err error
	ref.Hash, err = hash(path)
	if err != nil {
		err = fmt.Errorf("couldn't create a libref for %s: %v", path, err)
	}
	return ref, err
}

func newLibRefUnhashed(lib, path string, size int64) LibRef {
	return LibRef{removeLibPart(lib, path), "", size}
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
			return err
		}
		if info.IsDir() || !jpegRegexp.MatchString(path) && !mp4Regexp.MatchString(path) {
			return nil
		}

		out := fmt.Sprintf("checking %s", removeLibPart(lib, path))
		if maxLength < len(out) {
			maxLength = len(out)
		}
		fmt.Printf("\r%s", pad(out, maxLength))
		var ref LibRef
		if deep {
			ref, err = NewLibRef(lib, path, info.Size())
			if err != nil {
				return err
			}
		} else {
			ref = newLibRefUnhashed(lib, path, info.Size())
		}
		res = append(res, ref)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't traverse folder structure: %v", err)
	}

	fmt.Printf("\n\r%s", "")
	return
}

func Index(lib string) (map[string]LibRef, error) {
	sealed := map[string]LibRef{}
	content, err := ioutil.ReadFile(filepath.Join(lib, "index.csv"))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]LibRef{}, nil
		}
		return nil, fmt.Errorf("couldn't open index file: %v", err)
	}
	r := csv.NewReader(bytes.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("couldn't parse index file: %v", err)
	}
	for _, record := range records {
		size, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse index entry: %v", err)
		}

		sealed[record[0]] = LibRef{record[0], record[2], size}
	}
	return sealed, nil
}

func hash(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return fmt.Sprintf("%x", md5.Sum(data)), err
}

func removeLibPart(lib, path string) string {
	for lib[len(lib)-1] != '/' {
		lib += "/"
	}
	return strings.Replace(path, lib, "", -1)
}

func pad(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}

func Save(lib string, refs []LibRef) {
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
