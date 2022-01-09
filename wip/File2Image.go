package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"math"
	"os"

	"github.com/vaiktorg/grimoire/quickdraw"

	"github.com/vaiktorg/grimoire/errs"
)

func main() {
	file, err := os.Open("grimoire/wip/reggae.mp3")
	errs.Must(err)
	defer file.Close()

	st, err := file.Stat()
	errs.Must(err)

	scn := bufio.NewScanner(file)
	total := float64(st.Size())

	height := math.Floor(math.Sqrt(total))
	// to determine the width, we floor and ceil the length divided by the height
	width := int(math.Floor(math.Ceil(total / height)))

	img := quickdraw.NewCanvas(width/4, int(height/4))
	for scn.Scan() {
		err = scn.Err()
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		mapByteToColor(scn.Bytes(), width, img)
	}

	err = img.ExportToPNG()
	fmt.Println(err)
}

func mapByteToColor(byteArr []byte, width int, img *quickdraw.Canvas) {
	for i := 0; i < len(byteArr); i += 4 {
		if i+4 >= len(byteArr) {
			break
		}

		x := i % width
		y := i / width

		r, g, b, a := byteArr[i], byteArr[i+1], byteArr[i+2], byteArr[i+3]

		img.DrawPixel(x, y, color.NRGBA{
			R: r,
			G: g,
			B: b,
			A: a,
		})
	}
}
