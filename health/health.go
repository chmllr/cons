package health

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/chmllr/imgtb/seal"
)

// Verify matches new file list against existing records
func Verify(lib string, hashes []seal.LibRef) (
	corrupted []string,
	found, sealed map[string]string,
	duplicates map[string][]string) {
	found = map[string]string{}
	duplicates = map[string][]string{}
	for _, e := range hashes {
		found[e.Path] = e.Hash
		duplicates[e.Hash] = append(duplicates[e.Hash], e.Path)
	}
	sealed = map[string]string{}
	content, err := ioutil.ReadFile(filepath.Join(lib, "checksums.txt"))
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		fields := strings.Split(line, "::")
		if len(fields) != 2 {
			log.Fatalf("unexpected line: %q", line)
		}
		sealed[fields[0]] = fields[1]
	}

	for path, hash1 := range found {
		hash2, ok := sealed[path]
		if !ok {
			continue
		}
		if ok && hash1 != hash2 {
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
