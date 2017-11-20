package main

import (
	"flag"
	"fmt"
	"github.com/iand/salience"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadImage(filepath string) (image.Image, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()
	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func SaveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
	return nil
}

func main() {

	width := flag.Int("width", 200, "...")
	height := flag.Int("height", 200, "...")

	flag.Parse()

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		// only does PNG files...
		// im, err := halfgone.LoadImage(abs_path)
		im, err := LoadImage(abs_path)

		if err != nil {
			log.Fatal(err)
		}

		cropped := salience.Crop(im, *width, *height)

		root := filepath.Dir(abs_path)
		fname := filepath.Base(abs_path)
		ext := filepath.Ext(abs_path)

		new_ext := fmt.Sprintf("-%s%s", "crop", ext)
		fname = strings.Replace(fname, ext, new_ext, -1)

		new_path := filepath.Join(root, fname)
		err = SaveImage(cropped, new_path)

		if err != nil {
			log.Fatal(err)
		}
	}

}
