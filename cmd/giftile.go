package main

import (
	"flag"
	"github.com/nfnt/resize"	
	"github.com/straup/go-image-tools/util"
	"image"
	"image/draw"
	"log"
	"math"
	"os"
)

func main() {

	flag.Parse()

	cols := 10
	rows := 0

	for _, path := range flag.Args() {

		g, err := util.DecodeAnimatedGIF(path)

		if err != nil {
			log.Fatal(err)
		}

		for {
			count := len(g.Image)
			rows = int(math.Ceil(float64(count) / float64(cols)))

			if rows > cols {
				cols += 1
			} else {
				break
			}
		}

		bounds := g.Image[0].Bounds()
		dims := bounds.Max

		w := dims.X / 2
		h := dims.Y / 2

		margin := 10

		total_w := (w + (margin * 2)) * cols
		total_h := (h + (margin * 2)) * rows
		log.Println(total_w, total_h)
		log.Println(cols, rows)

		images := make([]image.Image, 0)

		for _, im := range g.Image {

			i := im.SubImage(im.Bounds())
			i = resize.Thumbnail(uint(w), uint(h), i, resize.Lanczos3)
			
			images = append(images, i)
		}

		count := len(images)

		out := image.NewRGBA(image.Rect(0, 0, total_w, total_h))
		idx := 0

		for i := 0; i < rows; i++ {

			for j := 0; j < cols; j++ {

				x1 := j * (w + (margin * 2))
				y1 := i * (h + (margin * 2))
				x2 := x1 + w
				y2 := y1 + h

				r := image.Rect(x1, y1, x2, y2)
				log.Println(count, idx, r)

				draw.Draw(out, r, images[idx], image.ZP, draw.Over)
				idx++

				if idx == count {
					break
				}
			}

			if idx == count {
				break
			}

		}

		fh, err := os.Create("test.png")

		if err != nil {
			log.Fatal(err)
		}

		err = util.EncodeImage(out, "png", fh)

		if err != nil {
			log.Fatal(err)
		}

	}
}
