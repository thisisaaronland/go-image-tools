package main

import (
	"flag"
	"github.com/jung-kurt/gofpdf"
	"github.com/straup/go-image-tools/util"
	"log"
	"path/filepath"
)

func main() {

	var orientation = flag.String("orientation", "P", "...")
	var size = flag.String("size", "letter", "...")
	var width = flag.Float64("width", 8.5, "...")
	var height = flag.Float64("height", 11, "...")
	var dpi = flag.Int("dpi", 150, "...")
	var filename = flag.String("filename", "picturebook.pdf", "...")
	var debug = flag.Bool("debug", false, "...")
	// var mode = flag.String("mode", "files", "...")

	flag.Parse()

	abs_filename, err := filepath.Abs(*filename)

	if err != nil {
		log.Fatal(err)
	}

	// https://godoc.org/github.com/jung-kurt/gofpdf#New

	var pdf *gofpdf.Fpdf

	if *size == "custom" {

		sz := gofpdf.SizeType{
			Wd: *width,
			Ht: *height,
		}

		init := gofpdf.InitType{
			OrientationStr: *orientation,
			UnitStr:        "in",
			SizeStr:        "",
			Size:           sz,
			FontDirStr:     "",
		}

		pdf = gofpdf.NewCustom(&init)

	} else {

		pdf = gofpdf.New(*orientation, "in", *size, "")
	}

	pdf.SetFillColor(0, 0, 0)

	// https://godoc.org/github.com/jung-kurt/gofpdf#ex-Fpdf-PageSize

	w, h, _ := pdf.PageSize(1)

	page_w := w * float64(*dpi)
	page_h := h * float64(*dpi)

	border_top := 1.0 * float64(*dpi)
	border_bottom := border_top * 1.5
	border_left := border_top * 0.8
	border_right := border_top * 0.8

	canvas_w := page_w - (border_left + border_right)
	canvas_h := page_h - (border_top + border_bottom)

	for i, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		im, format, err := util.DecodeImage(abs_path)

		if err != nil {
			log.Println(err)
			continue
		}

		info := pdf.GetImageInfo(abs_path)

		if info == nil {
			info = pdf.RegisterImage(abs_path, "")
		}

		info.SetDpi(float64(*dpi))

		dims := im.Bounds()

		x := border_left
		y := border_top

		w := float64(dims.Max.X)
		h := float64(dims.Max.Y)

		if *debug {
			log.Printf("[%d] %0.2f x %0.2f %0.2f x %0.2f\n", i, canvas_w, canvas_h, w, h)
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
				log.Printf("[%d] %0.2f (%0.2f) x %0.2f (%0.2f)\n", i, w, canvas_w, h, canvas_h)
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
			log.Printf("[%d] final %0.2f x %0.2f (%0.2f x %0.2f)\n", i, w, h, x, y)
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

		r_border := 0.01
		pdf.Rect((x - r_border), (y - r_border), (w + (r_border * 2)), (h + (r_border * 2)), "FD")

		pdf.ImageOptions(abs_path, x, y, w, h, false, opts, 0, "")
	}

	err = pdf.OutputFileAndClose(abs_filename)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(abs_filename)
}
