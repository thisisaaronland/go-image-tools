package util

import (
       "bufio"
       "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"       
       "os"
       "path/filepath"
)

// https://golang.org/src/image/decode_test.go

func DecodeImage(path string) (image.Image, string, error) {

		abs_path, err := filepath.Abs(path)

		if err != nil {
		return nil, "", err		
		}

	fh, err := os.Open(abs_path)

	if err != nil {
		return nil, "", err
	}

	defer fh.Close()

	return image.Decode(bufio.NewReader(fh))
}
