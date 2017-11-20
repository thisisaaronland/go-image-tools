package main

// https://maxhalford.github.io/blog/halftoning-1/

import (
	"flag"
	"fmt"
	"github.com/MaxHalford/halfgone"
	"github.com/nfnt/resize"
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

	mode := flag.String("mode", "atkinson", "...")

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

		dims := im.Bounds()
		w := uint(dims.Max.X)
		h := uint(dims.Max.Y)

		half_w := uint(float64(w) / 2.0)
		half_h := uint(float64(h) / 2.0)

		thumb := resize.Thumbnail(half_w, half_h, im, resize.Lanczos3)

		grey := halfgone.ImageToGray(thumb)

		switch *mode {
		case "atkinson":
			grey = halfgone.AtkinsonDitherer{}.Apply(grey)
		case "threshold":
			grey = halfgone.ThresholdDitherer{Threshold: 127}.Apply(grey)
		default:
			log.Fatal("Invalid or unsupported mode")
		}

		dither := resize.Resize(w, h, grey, resize.Lanczos3)

		root := filepath.Dir(abs_path)
		fname := filepath.Base(abs_path)
		ext := filepath.Ext(abs_path)

		new_ext := fmt.Sprintf("-%s%s", *mode, ext)
		fname = strings.Replace(fname, ext, new_ext, -1)

		new_path := filepath.Join(root, fname)
		err = SaveImage(dither, new_path)

		if err != nil {
			log.Fatal(err)
		}

	}
}
