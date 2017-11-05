package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
)

const blackIndex = 0

var (
	colors          []color.RGBA64
	lenColors       int
	lenColorsMinus1 float64
	aaDistance      int
)

type jColor struct {
	R int
	G int
	B int
}

type jColors []jColor

func loadColorFile(colorFile string) {
	raw, err := ioutil.ReadFile(colorFile)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Sticking with pre-determined colors")
		return
	}

	var vColors jColors
	json.Unmarshal(raw, &vColors)

	fmt.Printf("Loaded %d colors\n", len(vColors))
	if len(vColors) > 1 {
		parseColors(vColors)
	}
}

func parseColors(c jColors) {
	for k := range c {
		c[k].R = c[k].R << 8
		c[k].G = c[k].G << 8
		c[k].B = c[k].B << 8
	}
	lenCInt := len(c)
	lenColors = 100000

	// make palette evenly divisible of the incoming color array
	for lenColors%lenCInt > 0 {
		lenColors++
	}
	valuesPer := lenColors / lenCInt
	lenColors++ // one for black

	colors = make([]color.RGBA64, lenColors)
	lenColorsMinus1 = float64(lenColors - 1)
	aaDistance = int(lenColorsMinus1 * 0.05)

	// black
	colors[0].R = 0
	colors[0].G = 0
	colors[0].B = 0
	colors[0].A = 0xFFFF

	// interpolate colors
	cp1 := 0
	cp2 := 1
	cpCount := 0
	for cIndex := 1; cIndex < lenColors; cIndex++ {
		if cpCount >= valuesPer {
			cp1++
			cp2++
			if cp2 == lenCInt {
				cp2 = 0
			}
			cpCount = 0
		}

		colors[cIndex].R = uint16(lerp(float64(c[cp1].R), float64(c[cp2].R), float64(cpCount)/float64(valuesPer)))
		colors[cIndex].G = uint16(lerp(float64(c[cp1].G), float64(c[cp2].G), float64(cpCount)/float64(valuesPer)))
		colors[cIndex].B = uint16(lerp(float64(c[cp1].B), float64(c[cp2].B), float64(cpCount)/float64(valuesPer)))
		colors[cIndex].A = 0xFFFF

		cpCount++
	}

	if cpCount != valuesPer {
		fmt.Printf("color interpolation error: %d", cpCount)
	}
}

// this block sets up a basic color set
func init() {

	c := make([]jColor, 7)

	c[0].R = 0
	c[0].G = 255
	c[0].B = 255

	c[1].R = 255
	c[1].G = 0
	c[1].B = 0

	c[2].R = 255
	c[2].G = 0
	c[2].B = 255

	c[3].R = 0
	c[3].G = 255
	c[3].B = 0

	c[4].R = 255
	c[4].G = 255
	c[4].B = 0

	c[5].R = 0
	c[5].G = 0
	c[5].B = 255

	c[6].R = 0
	c[6].G = 255
	c[6].B = 255

	parseColors(c)
}
