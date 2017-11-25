package picturebook

import (
	"fmt"
	"github.com/straup/go-image-tools/halftone"
	"github.com/straup/go-image-tools/util"
	"io/ioutil"
	"os"
)

type PictureBookFilterFunc func(string) (bool, error)

type PictureBookPreProcessFunc func(string) (string, error)

func DefaultFilterFunc(string) (bool, error) {
	return true, nil
}

func DefaultPreProcessFunc(path string) (string, error) {
	return path, nil
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

	fh, err := ioutil.TempFile("", "halftone")

	if err != nil {
		return "", err
	}

	defer fh.Close()

	err = util.EncodeImage(dithered, format, fh)

	if err != nil {
		return "", err
	}

	// see what's going on here - this (appending the format
	// extension) is necessary because without it fpdf.GetImageInfo
	// gets confused and FREAKS out triggering fatal errors
	// along the way... oh well (20171125/thisisaaronland)

	fname := fh.Name()
	fh.Close()

	fq_fname := fmt.Sprintf("%s.%s", fname, format)

	err = os.Rename(fname, fq_fname)

	if err != nil {
		return "", err
	}

	return fq_fname, nil
}
