package health

import (
	"github.com/chmllr/cons/index"
)

// Verify matches new file list against existing records
func Verify(lib string, deep bool, hashes []index.LibRef) (
	corrupted []string,
	found, sealed map[string]index.LibRef,
	duplicates map[string][]string,
	err error) {
	found = map[string]index.LibRef{}
	duplicates = map[string][]string{}
	for _, e := range hashes {
		found[e.Path] = e
		if deep {
			duplicates[e.Hash] = append(duplicates[e.Hash], e.Path)
		}
	}
	sealed, err = index.Index(lib)
	if err != nil {
		return
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
