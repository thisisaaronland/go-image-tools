package main

import (
	"errors"
	"flag"
	"github.com/straup/go-image-tools/picturebook"
	"github.com/straup/go-image-tools/picturebook/functions"
	"log"
	"os"
)

func main() {

	err := Picturebook()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func Picturebook() error {

	var orientation = flag.String("orientation", "P", "...")
	var size = flag.String("size", "letter", "...")
	var width = flag.Float64("width", 8.5, "...")
	var height = flag.Float64("height", 11, "...")
	var dpi = flag.Float64("dpi", 150, "...")
	var border = flag.Float64("border", 0.01, "...")
	var caption = flag.String("caption", "default", "...")
	var filename = flag.String("filename", "picturebook.pdf", "...")
	var debug = flag.Bool("debug", false, "...")

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
	opts.Border = *border
	opts.Debug = *debug

	processed := make([]string, 0)

	defer func() {
		for _, p := range processed {

			go func(p string) {

				_, err := os.Stat(p)

				if os.IsNotExist(err) {
					return
				}

				os.Remove(p)
			}(p)
		}
	}()

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

	prep := func(path string) (string, error) {

		final := path

		for _, proc := range preprocess {

			switch proc {

			case "rotate":

				processed_path, err := functions.RotatePreProcessFunc(final)

				if err != nil {
					return "", err
				}

				if processed_path == "" {
					continue
				}

				processed = append(processed, processed_path)
				final = processed_path

			case "halftone":

				processed_path, err := functions.HalftonePreProcessFunc(final)

				if err != nil {
					return "", err
				}

				if processed_path == "" {
					continue
				}

				processed = append(processed, processed_path)
				final = processed_path

			default:
				return "", errors.New("Invalid or unsupported process")
			}
		}

		return final, nil
	}

	capt, err := functions.PictureBookCaptionFuncFromString(*caption)

	if err != nil {
		log.Fatal(err)
	}

	opts.Filter = filter
	opts.PreProcess = prep
	opts.Caption = capt

	pb, err := picturebook.NewPictureBook(opts)

	if err != nil {
		return err
	}

	sources := flag.Args()

	err = pb.AddPictures(sources)

	if err != nil {
		return err
	}

	err = pb.Save(*filename)

	if err != nil {
		return err
	}

	return nil
}
