package main

import (
	"flag"
	"github.com/straup/go-image-tools/picturebook"
	"log"
	"os"
)

func main() {

	var orientation = flag.String("orientation", "P", "...")
	var size = flag.String("size", "letter", "...")
	var width = flag.Float64("width", 8.5, "...")
	var height = flag.Float64("height", 11, "...")
	var dpi = flag.Float64("dpi", 150, "...")
	var filename = flag.String("filename", "picturebook.pdf", "...")
	var debug = flag.Bool("debug", false, "...")
	var mode = flag.String("mode", "files", "...")

	var include picturebook.RegexpFlag
	var exclude picturebook.RegexpFlag
	var preprocess picturebook.PreProcessFlag

	flag.Var(&include, "include", "...")
	flag.Var(&exclude, "exclude", "...")
	flag.Var(&preprocess, "pre-process", "...")

	flag.Parse()

	opts := picturebook.NewPictureBookDefaultOptions()
	opts.Orientation = *orientation
	opts.Size = *size
	opts.Width = *width
	opts.Height = *height
	opts.DPI = *dpi
	opts.Debug = *debug

	filter := func(path string) (bool, error) {

		for _, pat := range include {

			if !pat.MatchString(path) {
				return false, nil
			}
		}

		for _, pat := range exclude {

			if pat.MatchString(path) {
				return false, nil
			}
		}

		return true, nil
	}

	opts.Filter = filter

	pb, err := picturebook.NewPictureBook(opts)

	if err != nil {
		log.Fatal(err)
	}

	sources := flag.Args()

	err = pb.AddPictures(*mode, sources)

	if err != nil {
		log.Fatal(err)
	}

	err = pb.Save(*filename)

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
