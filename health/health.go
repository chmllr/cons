package health

import (
	"log"

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
	var err error
	sealed, err = seal.Registry(lib)
	if err != nil {
		log.Fatal(err)
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
