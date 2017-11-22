package picturebook

import (
	"context"
	"github.com/jung-kurt/gofpdf"
	"github.com/straup/go-image-tools/util"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"log"
	"path/filepath"
	"sync"
)

type PictureBookFilter func(string) (bool, error)

type PictureBook struct {
	pdf   *gofdpf.FPDF
	mu    *sync.Mutex
	Debug bool
}

func NewPictureBook() (*PictureBook, error) {

	mu := new(sync.Mutex)

	pb := PictureBook{
		Debug: false,
		mu:    mu,
	}

	return &pb, nil
}

func (pb *PictureBook) AddPictures(mode string, filter PictureBookFilter, paths []string) error {

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		abs_path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		ok, err = filter(abs_path)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		return pb.AddPicture(abs_path)
	}

	idx, err := index.NewIndexer(mode, cb)

	if err != nil {
		return err
	}

	return idx.IndexPaths(sources)
}

func (pb *PictureBook) AddPicture(abs_path string) error {

	im, format, err := util.DecodeImage(abs_path)

	if err != nil {
		return nil
	}

	pb.mu.Lock()

	info := pdf.GetImageInfo(abs_path)

	if info == nil {
		info = pdf.RegisterImage(abs_path, "")
	}

	pb.mu.Unlock()

	info.SetDpi(float64(*dpi))

	dims := im.Bounds()

	x := border_left
	y := border_top

	w := float64(dims.Max.X)
	h := float64(dims.Max.Y)

	if *debug {
		log.Printf("%0.2f x %0.2f %0.2f x %0.2f\n", canvas_w, canvas_h, w, h)
	}

	for {

		if w >= canvas_w || h >= canvas_h {

			if w > h || w > canvas_w {
				ratio := canvas_w / w

				w = canvas_w
				h = h * ratio

			} else {

				ratio := canvas_h / h
				w = w * ratio
				h = canvas_h
			}
		}

		if *debug {
			log.Printf("%0.2f (%0.2f) x %0.2f (%0.2f)\n", w, canvas_w, h, canvas_h)
		}

		if w <= canvas_w && h <= canvas_h {
			break
		}

	}

	if w < canvas_w {

		padding := canvas_w - w
		x = x + (padding / 2.0)
	}

	if h < (canvas_h - border_top) {

		y = y + border_top
	}

	if *debug {
		log.Printf("final %0.2f x %0.2f (%0.2f x %0.2f)\n", w, h, x, y)
	}

	pdf.AddPage()

	// https://godoc.org/github.com/jung-kurt/gofpdf#ImageOptions

	opts := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: format,
	}

	x = x / float64(*dpi)
	y = y / float64(*dpi)
	w = w / float64(*dpi)
	h = h / float64(*dpi)

	mu.Lock()

	r_border := 0.01

	if *debug {
		log.Println((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)))
	}

	pdf.Rect((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)), "FD")

	pdf.ImageOptions(abs_path, x, y, w, h, false, opts, 0, "")
	mu.Unlock()

	return nil
}

func (pb *PictureBook) Save(path string) error {

	return pb.pdf.OutputFileAndClose(path)
}
