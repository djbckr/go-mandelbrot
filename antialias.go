package main

import "image/color"

func (f *frame) antialias() {

	// allocate the actual color pixels we are setting
	f.pixels = make([][]*color.RGBA64, xSize)
	for k := range f.pixels {
		f.pixels[k] = make([]*color.RGBA64, ySize)
	}

	z := make([]int, len(rads))
	xSizeMinus1 := xSize - 1
	ySizeMinus1 := ySize - 1

	for x := 0; x < xSize; x++ {
		for y := 0; y < ySize; y++ {

			// outer edge; nothing to do
			if x == 0 || x == xSizeMinus1 || y == 0 || y == ySizeMinus1 {
				f.pixels[x][y] = &colors[f.palette[x][y]]
				continue
			}

			// get surrounding pixels
			z[0] = f.palette[x-1][y-1]
			z[1] = f.palette[x+1][y-1]
			z[2] = f.palette[x-1][y+1]
			z[3] = f.palette[x+1][y+1]

			// Check if surrounding pixels differ too much.
			// If so, resample, otherwise just set the pixel to the palette color
			if needAntiAlias(z) {
				var accumulator int
				f.pixels[x][y], accumulator = performAntiAlias(0, f.pointX(x), f.pointY(y), f.xDistance/2, f.yDistance/2)
				f.aaSuper += accumulator
			} else {
				f.aaDirect++
				f.pixels[x][y] = &colors[f.palette[x][y]]
			}

		}
	}

}

func performAntiAlias(depth int, xf, yf, xDistance, yDistance float64) (rslt *color.RGBA64, numSupers int) {
	numSupers = 0

	// resample surrounding area, keeping track of where the points were
	points := make([]complex128, len(rads))
	superSamples := make([]*color.RGBA64, len(rads))
	superIndexes := make([]int, len(rads))

	for k, v := range rads {
		xg := xf + (xDistance * real(v))
		yg := yf + (yDistance * imag(v))
		points[k] = complex(xg, yg)
		superIndexes[k] = mandelbrot(points[k])
		superSamples[k] = &colors[superIndexes[k]]
		numSupers++
	}

	// Check if super-sampled pixels differ too much, or if our depth is shallow enough.
	// If so, resample the superSamples again from each point, otherwise we are done
	if needAntiAlias(superIndexes) && depth < 3 {
		var accumulator int
		for k, v := range points {
			superSamples[k], accumulator = performAntiAlias(depth+1, real(v), imag(v), xDistance/2, yDistance/2)
			numSupers += accumulator
		}
	}

	rslt = avgSamples(superSamples)

	return
}

func needAntiAlias(z []int) bool {
	var zDist int // distance between two points
	var zK int    // prior value
	for k, v := range z {

		// skip first element, but store it for comparison
		if k == 0 {
			zK = v
			continue
		}

		// get the distance between this element and the last
		zDist = v - zK

		// if either value == blackIndex and there is any distance, we need AA
		if (v == blackIndex || zK == blackIndex) && zDist != 0 {
			return true
		}

		// abs(zDist)
		if zDist < 0 {
			zDist = -zDist
		}

		// if distance is further than aaDistance, we need AA
		if zDist > aaDistance {
			return true
		}

		// store this value for next comparison
		zK = v
	}

	return false
}

func avgSamples(z []*color.RGBA64) *color.RGBA64 {

	lenZ := float64(len(z))

	// int to prevent overflow of uint16
	var r1 int
	var g1 int
	var b1 int

	for _, v := range z {
		r1 += int(v.R)
		g1 += int(v.G)
		b1 += int(v.B)
	}

	return &color.RGBA64{
		uint16(float64(r1) / lenZ),
		uint16(float64(g1) / lenZ),
		uint16(float64(b1) / lenZ),
		0xFFFF}
}
