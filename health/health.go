package health

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chmllr/imgtb/seal"
)

// Verify matches new file list against existing records
func Verify(lib string, deep bool, hashes []seal.LibRef) (
	corrupted []string,
	found, sealed map[string]seal.LibRef,
	duplicates map[string][]string) {
	found = map[string]seal.LibRef{}
	duplicates = map[string][]string{}
	for _, e := range hashes {
		found[e.Path] = e
		if deep {
			duplicates[e.Hash] = append(duplicates[e.Hash], e.Path)
		}
	}
	sealed = map[string]seal.LibRef{}
	content, err := ioutil.ReadFile(filepath.Join(lib, "checksums.txt"))
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		// TODO: move it LibRef's method
		fields := strings.Split(line, "::")
		if len(fields) != 3 {
			log.Fatalf("unexpected line: %q", line)
		}
		parsedSize, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			panic(err) // TODO: handle
		}

		sealed[fields[0]] = seal.LibRef{fields[0], fields[2], parsedSize}
	}

	for path, lr := range found {
		lr2, ok := sealed[path]
		if !ok {
			continue
		}
		if ok && (lr.Size != lr2.Size || deep && deep && lr.Hash != lr2.Hash) {
			corrupted = append(corrupted, path)
		}
		delete(found, path)
		delete(sealed, path)
	}
	for hash, paths := range duplicates {
		if len(paths) < 2 {
			delete(duplicates, hash)
		}
	}
	return
}
