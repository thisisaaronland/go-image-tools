package picturebook

import (
	"context"
	"io"
)

type PictureBookFilterFunc func(string) (bool, error)

type PictureBookPreProcessFunc func(io.Reader, context.Context) (io.Reader, error)

func DefaultFilterFunc(string) (bool, error) {
	return true, nil
}

func DefaultPreProcessFunc(fh io.Reader, ctx context.Context) (io.Reader, error) {
	return fh, nil
}
