package util

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

func EncodeImage(im image.Image, format string, wr io.Writer) error {

	var err error

	switch format {
	case "jpeg":
		opts := jpeg.Options{Quality: 100}
		err = jpeg.Encode(wr, im, &opts)
	case "png":
		err = png.Encode(wr, im)
	case "gif":
		opts := gif.Options{}
		err = gif.Encode(wr, im, &opts)
	default:
		err = errors.New("Invalid or unsupported format")
	}

	return err
}
