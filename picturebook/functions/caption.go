package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

	var rsp gjson.Result
	var title string
	var added int64
	
	rsp = gjson.GetBytes(body, "refers_to.title")

	if !rsp.Exists() {
		return "", errors.New("Object information missing title")
	}

	title = rsp.String()

	rsp = gjson.GetBytes(body, "created")

	if rsp.Exists() {
	   added = rsp.Int()
	}

	caption := fmt.Sprintf("%s (added %d)", title, added)
	
	return caption, nil
}
