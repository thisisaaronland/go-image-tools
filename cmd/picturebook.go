package main

import (
	"flag"
	"github.com/straup/go-image-tools/picturebook"
	"log"
	"os"
	"strings"
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
	var suffix = flag.String("with-suffix", "", "...")

	flag.Parse()

	opts := picturebook.NewPictureBookDefaultOptions()
	opts.Orientation = *orientation
	opts.Size = *size
	opts.Width = *width
	opts.Height = *height
	opts.DPI = *dpi
	opts.Debug = *debug

	if *suffix != "" {

		filter := func(path string) (bool, error) {

			ok := strings.HasSuffix(path, *suffix)
			return ok, nil
		}

		opts.Filter = filter
	}

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
