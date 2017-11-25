package picturebook

import (
	"context"
	"github.com/jung-kurt/gofpdf"
	"github.com/whosonfirst/go-whosonfirst-index"
	"io"
	"log"
	_ "os"
	"sync"
)

type PictureBookOptions struct {
	Orientation string
	Size        string
	Width       float64
	Height      float64
	DPI         float64
	Filter      PictureBookFilterFunc
	PreProcess  PictureBookPreProcessFunc
	Debug       bool
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
	PDF     *gofpdf.Fpdf
	Mutex   *sync.Mutex
	Border  PictureBookBorder
	Canvas  PictureBookCanvas
	Options PictureBookOptions
	pages   int
}

func NewPictureBookDefaultOptions() PictureBookOptions {

	filter := DefaultFilterFunc
	prep := DefaultPreProcessFunc

	opts := PictureBookOptions{
		Orientation: "P",
		Size:        "letter",
		Width:       0.0,
		Height:      0.0,
		DPI:         150.0,
		Filter:      filter,
		PreProcess:  prep,
		Debug:       false,
	}

	return opts
}

func NewPictureBook(opts PictureBookOptions) (*PictureBook, error) {

	var pdf *gofpdf.Fpdf

	if opts.Size == "custom" {

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
	border_left := border_top * 1.0
	border_right := border_top * 1.0

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
		Height: canvas_h,
	}

	mu := new(sync.Mutex)

	pb := PictureBook{
		PDF:     pdf,
		Mutex:   mu,
		Border:  b,
		Canvas:  c,
		Options: opts,
		pages:   0,
	}

	return &pb, nil
}

func (pb *PictureBook) AddPictures(mode string, paths []string) error {

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		abs_path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		ok, err := pb.Options.Filter(abs_path)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		processed_path, err := pb.Options.PreProcess(abs_path)

		if err != nil {
			return nil
		}

		pb.Mutex.Lock()
		pb.pages += 1
		pagenum := pb.pages
		pb.Mutex.Unlock()

		err = pb.AddPicture(pagenum, processed_path)

		if err != nil {
			return err
		}

		return nil
	}

	idx, err := index.NewIndexer(mode, cb)

	if err != nil {
		return err
	}

	return idx.IndexPaths(paths)
}

func (pb *PictureBook) AddPicture(pagenum int, abs_path string) error {

	pb.Mutex.Lock()
	defer pb.Mutex.Unlock()

	info := pb.PDF.GetImageInfo(abs_path)

	if info == nil {
		info = pb.PDF.RegisterImage(abs_path, "")
	}

	info.SetDpi(pb.Options.DPI)

	w := info.Width() * pb.Options.DPI
	h := info.Height() * pb.Options.DPI

	if pb.Options.Debug {
		log.Printf("[%d] %s %02.f x %02.f\n", pagenum, abs_path, w, h)
	}

	x := pb.Border.Left
	y := pb.Border.Top

	if pb.Options.Debug {
		log.Printf("[%d] canvas: %0.2f x %0.2f image: %0.2f x %0.2f\n", pagenum, pb.Canvas.Width, pb.Canvas.Height, w, h)
	}

	for {

		if w >= pb.Canvas.Width || h >= pb.Canvas.Height {

			if w > h {

				ratio := pb.Canvas.Width / w
				w = pb.Canvas.Width
				h = h * ratio

			} else if w > pb.Canvas.Width {

				ratio := pb.Canvas.Width / w
				w = pb.Canvas.Width
				h = h * ratio

			} else if h > w {

				ratio := pb.Canvas.Height / h
				w = w * ratio
				h = pb.Canvas.Height

			} else if h > pb.Canvas.Height {

				ratio := pb.Canvas.Height / h
				w = w * ratio
				h = pb.Canvas.Height

			} else {
			}

		}

		if pb.Options.Debug {
			log.Printf("[%d] w: %0.2f (max w: %0.2f)  h: %0.2f (max h: %0.2f)\n", pagenum, w, pb.Canvas.Width, h, pb.Canvas.Height)
		}

		if w <= pb.Canvas.Width && h <= pb.Canvas.Height {
			break
		}

	}

	if w < pb.Canvas.Width {

		padding := pb.Canvas.Width - w
		x = x + (padding / 2.0)
	}

	// if pb.Canvas.Height > pb.Canvas.Width && h < (pb.Canvas.Height - pb.Border.Top) {
	if h < (pb.Canvas.Height - pb.Border.Top) {

		y = y + pb.Border.Top
	}

	if pb.Options.Debug {
		log.Printf("[%d] final %0.2f x %0.2f (%0.2f x %0.2f)\n", pagenum, w, h, x, y)
	}

	pb.PDF.AddPage()

	// https://godoc.org/github.com/jung-kurt/gofpdf#ImageOptions

	opts := gofpdf.ImageOptions{
		ReadDpi:   false,
		ImageType: "",
	}

	x = x / pb.Options.DPI
	y = y / pb.Options.DPI
	w = w / pb.Options.DPI
	h = h / pb.Options.DPI

	r_border := 0.01

	/*
		if pb.Options.Debug {
			log.Println((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)))
		}
	*/

	pb.PDF.Rect((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)), "FD")

	pb.PDF.ImageOptions(abs_path, x, y, w, h, false, opts, 0, "")

	return nil
}

func (pb *PictureBook) Save(path string) error {

	return pb.PDF.OutputFileAndClose(path)
}
