package main

import (
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
	"time"
)

type frame struct {
	frameNumber int
	// an array of indexes into the color array (our raw color pointers, if you will)
	palette [][]int
	// our final color
	pixels [][]*color.RGBA64
	// based on zoom, the min/max values we will use to calculate the mandelbrot
	xMin float64
	xMax float64
	yMin float64
	yMax float64
	// distance between pixels (at mandelbrot measures and zoom)
	xDistance float64
	yDistance float64
	// performance counters
	aaDirect int
	aaSuper  int
	// the difference between xMin/xMax and yMin/yMax - a minor optimization to find a point
	xDiff float64
	yDiff float64
}

func renderFrame(frameNumberChannel chan int) {

	// while the channel is open...
	for frameNumber := range frameNumberChannel {

		f := &frame{}
		startTime := time.Now()
		f.init(frameNumber)
		f.fillPalette()
		f.antialias()
		f.save()

		time.Sleep(1 * time.Second)
		duration := time.Now().Sub(startTime)
		fmt.Printf("Finished %d : aaDirect %d : aaSuper %d : duration %v\n", f.frameNumber, f.aaDirect, f.aaSuper, duration)
	}

	// channel has been closed; notifiy that this goroutine is finished
	doneChannel <- 1
}

// frame initialization: mostly to setup zoom (xMin/xMax and yMin/yMax)
func (f *frame) init(fr int) {

	f.frameNumber = fr

	f.palette = make([][]int, xSize)
	for k := range f.palette {
		f.palette[k] = make([]int, ySize)
	}

	// zoom/center
	zoomFactor := zoomStart * math.Pow(magnificationFactor, float64(fr)-1)

	f.xMin = (xMin / zoomFactor) + centerX
	f.xMax = (xMax / zoomFactor) + centerX
	f.yMin = (yMin / zoomFactor) + centerY
	f.yMax = (yMax / zoomFactor) + centerY

	f.xDiff = f.xMax - f.xMin
	f.yDiff = f.yMax - f.yMin

	// distance between pixels
	f.xDistance = math.Abs(f.pointX(0) - f.pointX(1))
	f.yDistance = math.Abs(f.pointY(0) - f.pointY(1))

	fmt.Printf("Frame %d : zoomFactor %f\n", f.frameNumber, zoomFactor)

}

// Get all the raw colors for each pixel. This is the easy part.
// Once we get these colors, we then need to antialias them.
func (f *frame) fillPalette() {
	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {
			f.palette[x][y] = mandelbrot(complex(f.pointX(x), f.pointY(y)))
		}
	}
}

func (f *frame) pointX(x int) float64 {
	return f.xMin + (f.xDiff)*float64(x)/xFSizeMinus1
}

func (f *frame) pointY(y int) float64 {
	return f.yMin + (f.yDiff)*float64(y)/yFSizeMinus1
}

// The workhorse. I don't know what this does :)
func mandelbrot(C complex128) int {
	i := 0
	z := C
	for cmplx.Abs(z) < 2 && i < maxIter {
		z = z*z + C
		i++
	}

	if i < maxIter {
		z = z*z + C
		i++
		z = z*z + C
		i++
		mu := float64(i) - (math.Log(math.Log(doMod(z))))/logEscapeRadius
		colorIndex := int(mu/maxIter*lenColorsMinus1) + 1 // the +1 moves the index away from blackIndex

		for colorIndex >= lenColors {
			colorIndex = colorIndex - (lenColors - 1)
		}

		if colorIndex < 1 {
			colorIndex = 1
		}

		return colorIndex
	}

	return blackIndex
}

func doMod(z complex128) float64 {
	r := real(z)
	i := imag(z)
	return math.Sqrt(r*r + i*i)
}
