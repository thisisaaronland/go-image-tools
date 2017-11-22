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

type PictureBookFilterFunc func(string) (bool, error)

type PictureBookOptions struct {
	Orientation string
	Size        string
	Width       float64
	Height      float64
	DPI         float64
	Filter      PictureBookFilterFunc
}

type PictureBookBorder struct {
	Top    float64
	Bottom float64
	Left   float64
	Right  float64
}

type PictureBookCanvas struct {
	Width  float64
	Height float64
}

type PictureBook struct {
	PDF     *gofdpf.FPDF
	Mutex   *sync.Mutex
	Border  PictureBookBorder
	Canvas  PictureBookCanvas
	Options PictureBookOptions
	Debug   bool
}

func NewPictureBookDefaultOptions() PictureBookOptions {

	filter := func(string) (bool, error) {

		return ok, nil
	}

	opts := PictureBookOptions{
		Orientation: "P",
		Size:        "letter",
		Width:       0.0,
		Height:      0.0,
		DPI:         150.0,
		Filter:      filter,
	}

	return opts
}

func NewPictureBook(opts PictureBookOptions) (*PictureBook, error) {

	var pdf *gofpdf.Fpdf

	if *size == "custom" {

		sz := gofpdf.SizeType{
			Wd: opts.Width,
			Ht: opts.Height,
		}

		init := gofpdf.InitType{
			OrientationStr: opts.Orientation,
			UnitStr:        "in",
			SizeStr:        "",
			Size:           sz,
			FontDirStr:     "",
		}

		pdf = gofpdf.NewCustom(&init)

	} else {

		pdf = gofpdf.New(opts.Orientation, "in", opts.Size, "")
	}

	w, h, _ := pdf.PageSize(1)

	page_w := w * opts.DPI
	page_h := h * opts.DPI

	border_top := 1.0 * opts.DPI
	border_bottom := border_top * 1.5
	border_left := border_top * 0.8
	border_right := border_top * 0.8

	canvas_w := page_w - (border_left + border_right)
	canvas_h := page_h - (border_top + border_bottom)

	b := PictureBookBorder{
		Top:    border_top,
		Bottom: border_bottom,
		Left:   border_left,
		Right:  border_right,
	}

	c := PictureBookCanvas{
		Width:  canvas_w,
		Height: pb.Canvas.Height,
	}

	mu := new(sync.Mutex)

	pb := PictureBook{
		Debug:   false,
		PDF:     pdf,
		Mutex:   mu,
		Border:  b,
		Canvas:  c,
		Options: opts,
	}

	return &pb, nil
}

func (pb *PictureBook) AddPictures(mode string, paths []string) error {

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		abs_path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		ok, err = pb.Options.Filter(abs_path)

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

	pb.Mutex.Lock()

	info := pb.PDF.GetImageInfo(abs_path)

	if info == nil {
		info = pb.PDF.RegisterImage(abs_path, "")
	}

	pb.Mutex.Unlock()

	info.SetDpi(pb.Options.DPI)

	dims := im.Bounds()

	x := pb.Border.Left
	y := pb.Border.Top

	w := float64(dims.Max.X)
	h := float64(dims.Max.Y)

	if *debug {
		log.Printf("%0.2f x %0.2f %0.2f x %0.2f\n", pb.Canvas.Width, pb.Canvas.Height, w, h)
	}

	for {

		if w >= pb.Canvas.Width || h >= pb.Canvas.Height {

			if w > h || w > pb.Canvas.Width {
				ratio := pb.Canvas.Width / w

				w = pb.Canvas.Width
				h = h * ratio

			} else {

				ratio := pb.Canvas.Height / h
				w = w * ratio
				h = pb.Canvas.Height
			}
		}

		if *debug {
			log.Printf("%0.2f (%0.2f) x %0.2f (%0.2f)\n", w, pb.Canvas.Width, h, pb.Canvas.Height)
		}

		if w <= pb.Canvas.Width && h <= pb.Canvas.Height {
			break
		}

	}

	if w < pb.Canvas.Width {

		padding := pb.Canvas.Width - w
		x = x + (padding / 2.0)
	}

	if h < (pb.Canvas.Height - pb.Border.Top) {

		y = y + pb.Border.Top
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

	x = x / pb.Options.DPI
	y = y / pb.Options.DPI
	w = w / pb.Options.DPI
	h = h / pb.Options.DPI

	pb.Mutex.Lock()

	r_border := 0.01

	if *debug {
		log.Println((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)))
	}

	pb.PDF.Rect((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)), "FD")

	pb.PDF.ImageOptions(abs_path, x, y, w, h, false, opts, 0, "")
	pb.Mutex.Unlock()

	return nil
}

func (pb *PictureBook) Save(path string) error {

	return pb.PDF.OutputFileAndClose(path)
}
