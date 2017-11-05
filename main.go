package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
)

const (
	// X is width, Y is height

	// we'll do an aspect ratio of 16:9 (HDTV)
	xRatio = float64(16)
	yRatio = float64(9)

	// the vertical mandelbrot range
	yMin = float64(-1)
	yMax = float64(1)

	// maximum "depth" to calculate if a point is in the mandelbrot
	maxIter = 100

	// If we zoom to the middle of the image, we'll just get black.
	// So we'll move the image a bit to zoom into something interesting
	centerX = -0.743643887037158704752191506114774
	centerY = 0.131825904205311970493132056385139

	superSampleLen = 4
)

var (
	numProcesses = runtime.NumCPU()
	// holds a queue of frame numbers to process
	workChannel = make(chan int, numProcesses*2)

	// used to signify when the goroutine is finished
	doneChannel = make(chan int, numProcesses)

	// we'll use a separate goroutine for file saves
	fileSaverChannel = make(chan *frame, numProcesses*10)
	fileSaverDone    = make(chan int)

	// the horizontal mandelbrot range - calculated from the vertical range and the aspect ratio
	xMin float64
	xMax float64

	// used for smoothing calculations
	logEscapeRadius = math.Log(2)

	// an array of point multipliers to move a point a little for antialiasing
	rads [superSampleLen]complex128

	// the rate at which to magnify per frame
	magnificationFactor float64

	// magnification start and end values
	zoomStart float64
	zoomEnd   float64

	// number of frames between zoomStart and zoomEnd
	numFrames int

	// where to pick up and stop rendering
	frameStart int
	frameEnd   int

	// image width (x) and height (y)
	xSize int
	ySize int

	// a minor optimization used during rendering
	xFSizeMinus1 float64
	yFSizeMinus1 float64

	colorFile string
)

func main() {

	flag.Parse()

	// calculate the vertical size
	ySize = int(float64(xSize) * (yRatio / xRatio))
	xFSizeMinus1 = float64(xSize - 1)
	yFSizeMinus1 = float64(ySize - 1)

	if colorFile != "" {
		loadColorFile(colorFile)
	}

	fmt.Printf("xSize=%d : ySize=%d : magX=%f : frameStart=%d : frameEnd=%d : numFrames=%d\n", xSize, ySize, magnificationFactor, frameStart, frameEnd, numFrames)

	// setup go-routines that will take a queue of operations
	for i := 1; i <= numProcesses; i++ {
		go renderFrame(workChannel)
	}

	// a single goprocess for file saves
	go saveImage(fileSaverChannel)

	// fill pipe with frame number (basically work-id)
	for frameNumber := frameStart; frameNumber <= frameEnd; frameNumber++ {
		workChannel <- frameNumber
	}

	// the renderFrame() method will stop
	// running when reaching the end of the channel
	close(workChannel)

	// wait for all processes to finish
	for i := 1; i <= numProcesses; i++ {
		<-doneChannel
	}

	// since all the renderers are done, close the file saver
	close(fileSaverChannel)

	// and wait for any outstanding file saves
	<-fileSaverDone

}

func init() {

	flag.IntVar(&numFrames, "numFrames", 3600, "The number of total frames that used used to calculate zoom at any particular frame.") // at 60fps, this is 1 minute
	flag.Float64Var(&zoomStart, "zoomStart", 2, "Start zoom factor; probably want to stick with this default.")                        // everybody knows what zoom x1 looks like; start a little further in
	flag.Float64Var(&zoomEnd, "zoomEnd", 100000, "End zoom factor; exceeding this default isn't recommended.")                         // this is about where float64 has its limits
	flag.IntVar(&frameStart, "frameStart", 1, "Which frame to start rendering.")                                                       // assume we're rendering all frames
	flag.IntVar(&frameEnd, "frameEnd", 3600, "Which frame to end rendering")                                                           // assume we're rendering all frames
	flag.IntVar(&xSize, "xSize", 3840, "Width of image in pixels; default is a 4k UHD image.")                                         // assume we're rendering 4k
	flag.StringVar(&colorFile, "colorFile", "colors.json", "A color description file. The contents of this file is a single array, with objects of \"R\", \"G\", and \"B\" attributes. Each attribute can be 0 to 255.")

	// calculate where xMin and xMax go, based on horizontal center between -2 and +0.5
	hCenter := lerp(-2.0, 0.5, 0.5)
	sz := (yMax - yMin) * (xRatio / yRatio)

	xMin = hCenter - sz/2
	xMax = hCenter + sz/2

	rads[0] = complex(-1, -1)
	rads[1] = complex(-1, 1)
	rads[2] = complex(1, -1)
	rads[3] = complex(1, 1)

	magnificationFactor = math.Pow(zoomEnd/zoomStart, 1.0/(float64(numFrames)-1.0))

}
