package health

// Verify matches new file list against existing records
func Verify(lib string, mapping1, mapping2 map[string]string) (corrupted []string) {
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
