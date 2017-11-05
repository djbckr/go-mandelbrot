package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func saveImage(ch chan *frame) {
	// while channel is open...
	for f := range ch {

		fn := fmt.Sprintf("%04d", f.frameNumber) + ".png"

		file, err := os.Create(fn)
		if err != nil {
			log.Fatal(err)
		}

		img := image.NewRGBA64(image.Rect(0, 0, int(xSize), int(ySize)))

		for k1, v1 := range f.pixels {
			for k2, v2 := range v1 {
				img.SetRGBA64(k1, k2, v2)
			}
		}

		if err := png.Encode(file, img); err != nil {
			file.Close()
			log.Fatal(err)
		}

		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}

	// indicate we are finished
	fileSaverDone <- 1
}
