package main

// linear interpolation; given a start and end point, and a point between 0.0 and 1.0
func lerp(v0 float64, v1 float64, t float64) float64 {
	return v0 + t*(v1-v0)
}

// find point on a slope; given point(x1, y1) and point(x2, y2), find Y for X
func findSlopePoint(x1, y1, x2, y2, x float64) float64 {
	return ((x - x1) * ((y2 - y1) / (x2 - x1))) + y1
}

// convert pixel point to the value we are calculating the mandelbrot set with
// essentially an optimized findSlopePoint() function
func getActual(p int, min float64, max float64, size int) float64 {
	return min + (max-min)*float64(p)/(float64(size)-1)
}
