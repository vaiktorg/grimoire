package quickdraw

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"sync"
	"time"

	"github.com/vaiktorg/grimoire/errs"

	"github.com/vaiktorg/grimoire/helpers"
)

type Canvas struct {
	sync.Mutex
	img *image.RGBA
}

func main() {
	t := time.Now()
	can := NewCanvas(600, 420)
	can.DrawLine(0, 0, 300, 200, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	err := can.ExportToPNG()
	errs.Must(err)

	fmt.Println(time.Since(t))
}

func NewCanvas(w, h int) *Canvas {
	return &Canvas{
		img: image.NewRGBA(image.Rect(0, 0, w, h)),
	}
}

func (c *Canvas) ExportToPNG() error {
	file := helpers.OpenFile("CanvasTestFile.png")
	err := png.Encode(file, c.img)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (c *Canvas) DrawPixel(x, y int, colour color.Color) {
	c.Lock()
	defer c.Unlock()
	c.img.Set(x, y, colour)
}

func (c *Canvas) DrawLine(x0, y0, x1, y1 int, colour color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	d := dy*2 - dx
	incrE := dy * 2
	incrNE := (dy - dx) * 2
	x := x0
	y := y0
	c.DrawPixel(x, y, colour)
	for x < x1 {
		if d <= 0 {
			d += incrE
			x++
		} else {
			d += incrNE
			x++
			y++
		}
		c.DrawPixel(x, y, colour)
	}
}
