package checksum

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var (
	jpegRegexp = regexp.MustCompile(`(?i)\.jpe?g`)
	mp4Regexp  = regexp.MustCompile(`(?i)\.mp4`)
)

func Report(lib string) (string, error) {
	var buf bytes.Buffer
	err := filepath.Walk(lib, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if info.IsDir() || !jpegRegexp.MatchString(path) && !mp4Regexp.MatchString(path) {
			return nil
		}

		h, err := hash(path)
		if err != nil {
			return err
		}
		buf.WriteString(path)
		buf.WriteString("::")
		buf.WriteString(h)
		buf.WriteRune('\n')
		return nil
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func hash(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return fmt.Sprintf("%x", md5.Sum(data)), err
}
