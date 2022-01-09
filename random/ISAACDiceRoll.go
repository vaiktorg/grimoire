package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"

	"github.com/aybabtme/rgbterm"
	"github.com/vaiktorg/grimoire/errs"
	"github.com/vaiktorg/grimoire/helpers"

	"github.com/gtank/isaac"
	"github.com/vaiktorg/grimoire/uid"
)

type Die struct {
	rng  isaac.ISAAC
	Side int
}

func (d *Die) Seed() {
	d.rng.Seed(uid.NewUID(d.Side).String())
}

func (d *Die) Roll() int {
	return int(d.rng.Rand())%d.Side + 1
}

func main() {
	d := Die{Side: 100}
	d.Seed()
	d.Render()

	sides := d.Side
	for x := 0; x < sides; x++ {
		d.Roll()
	}
}

// DebugFuncs ----------------------------------------
func (d *Die) clamp(val int) color.RGBA64 {
	if val <= 1 {
		return color.RGBA64{R: 1<<16 - 1, A: 0}
	}

	if val >= d.Side {
		return color.RGBA64{B: 1<<16 - 1, A: 0}
	}

	return color.RGBA64{
		R: 1<<16 - 1,
		G: 1<<16 - 1,
		B: 1<<16 - 1,
		A: 1<<16 - 1,
	}
}

func (d *Die) heatMap(val float64, min, max int) color.RGBA64 {
	minimum, maximum := float64(min), float64(max)
	ratio := 2*val - minimum/(maximum-minimum)
	b := int(math.Max(0, 255*(1-ratio)))
	r := int(math.Max(0, 255*(ratio-1)))
	g := 255 - b - r
	return color.RGBA64{
		R: uint16(r << 8),
		G: uint16(g << 8),
		B: uint16(b << 8),
	}
}

func (d *Die) Render() {
	drawCell := func(img *image.RGBA64, X, Y, cellsize int, col color.RGBA64) {
		for y := Y; y < Y+cellsize<<2; y++ {
			for x := X; x < X+cellsize<<2; x++ {
				img.Set(x, y, col)
			}
		}
	}

	cellSize := 8
	img := image.NewRGBA64(image.Rect(0, 0, d.Side*cellSize, d.Side*cellSize))

	for y := 0; y < d.Side; y++ {
		fmt.Printf("%3d ", y+1)
		for x := 0; x < d.Side; x++ {
			col := d.heatMap(float64(d.Roll()), 1, d.Side)
			//col := d.clamp(d.Roll())
			drawCell(img, x*cellSize, y*cellSize, 16, col)

			fmt.Print(rgbterm.FgString("#", uint8(col.R>>8), uint8(col.G>>8), uint8(col.B>>8)))
		}
		fmt.Println("")
	}

	file := helpers.OpenFile("heatmap.jpg")
	defer file.Close()

	err := jpeg.Encode(file, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	errs.Must(err)
}
