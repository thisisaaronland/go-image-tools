package functions

import (
	"github.com/microcosm-cc/exifutil"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/straup/go-image-tools/halftone"
	"github.com/straup/go-image-tools/util"
	_ "log"
	"os"
	"path/filepath"
)

func DefaultPreProcessFunc(path string) (string, error) {
	return "", nil
}

// https://www.daveperrett.com/articles/2012/07/28/exif-orientation-handling-is-a-ghetto/

func RotatePreProcessFunc(path string) (string, error) {

	ext := filepath.Ext(path)

	if ext != ".jpg" && ext != ".jpeg" {
		return "", nil
	}

	fh, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	x, err := exif.Decode(fh)

	if err != nil {
		return "", err
	}

	tag, err := x.Get(exif.Orientation)

	if err != nil {
		return "", nil
	}

	// log.Println(path, tag)

	orientation, err := tag.Int64(0)

	if err != nil {
		return "", nil
	}

	if orientation == 1 {
		return "", nil
	}

	im, format, err := util.DecodeImage(path)

	if err != nil {
		return "", err
	}

	angle, _, _ := exifutil.ProcessOrientation(orientation)
	rotated := exifutil.Rotate(im, angle)

	return util.EncodeTempImage(rotated, format)
}

func HalftonePreProcessFunc(path string) (string, error) {

	im, format, err := util.DecodeImage(path)

	if err != nil {
		return "", err
	}

	opts := halftone.NewDefaultHalftoneOptions()
	dithered, err := halftone.Halftone(im, opts)

	if err != nil {
		return "", err
	}

	return util.EncodeTempImage(dithered, format)
}
