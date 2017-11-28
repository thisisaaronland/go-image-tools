package picturebook

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/straup/go-image-tools/halftone"
	"github.com/straup/go-image-tools/util"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

type PictureBookFilterFunc func(string) (bool, error)

type PictureBookPreProcessFunc func(string) (string, error)

type PictureBookCaptionFunc func(string) (string, error)

func DefaultFilterFunc(string) (bool, error) {
	return true, nil
}

func DefaultPreProcessFunc(path string) (string, error) {
	return path, nil
}

func PictureBookCaptionFuncFromString(caption string) (PictureBookCaptionFunc, error) {

	var capt PictureBookCaptionFunc

	switch caption {

	case "cooperhewitt":
		capt = CooperHewittShoeboxCaptionFunc
	case "default":
		capt = FilenameCaptionFunc
	case "filename":
		capt = FilenameCaptionFunc
	case "parent":
		capt = FilenameAndParentCaptionFunc
	case "none":
		capt = NoneCaptionFunc
	default:
		return nil, errors.New("Invalid caption type")
	}

	return capt, nil
}

func DefaultCaptionFunc(path string) (string, error) {
	return FilenameCaptionFunc(path)
}

func FilenameCaptionFunc(path string) (string, error) {

	fname := filepath.Base(path)
	return fname, nil
}

func FilenameAndParentCaptionFunc(path string) (string, error) {

	root := filepath.Dir(path)
	parent := filepath.Base(root)
	fname := filepath.Base(path)

	return filepath.Join(parent, fname), nil
}

func NoneCaptionFunc(path string) (string, error) {
	return "", nil
}

func CooperHewittShoeboxCaptionFunc(path string) (string, error) {

	root := filepath.Dir(path)
	info := filepath.Join(root, "index.json")

	_, err := os.Stat(info)

	if err != nil {
		return "", err
	}

	fh, err := os.Open(info)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	var item interface{}
	err = json.Unmarshal(body, &item)

	if err != nil {
		return "", err
	}

	rsp := gjson.GetBytes(body, "refers_to.title")

	if !rsp.Exists() {
		return "", errors.New("Object information missing title")
	}

	return rsp.String(), nil
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
