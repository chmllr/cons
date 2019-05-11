package health

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// Verify matches new file list against existing records
func Verify(lib string, hashes []struct{ Path, Hash string }) (corrupted []string, mapping1, mapping2 map[string]string) {
	mapping1 = map[string]string{}
	for _, e := range hashes {
		mapping1[e.Path] = e.Hash
	}
	mapping2 = map[string]string{}
	content, err := ioutil.ReadFile(filepath.Join(lib, "checksums.txt"))
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		fields := strings.Split(line, "::")
		if len(fields) != 2 {
			log.Fatalf("unexpected line: %q", line)
		}
		mapping2[fields[0]] = fields[1]
	}

	for path, hash1 := range mapping1 {
		hash2, found := mapping2[path]
		if !found {
			continue
		}
		if found && hash1 != hash2 {
			corrupted = append(corrupted, path)
		}
		delete(mapping1, path)
		delete(mapping2, path)
	}
	return
}
