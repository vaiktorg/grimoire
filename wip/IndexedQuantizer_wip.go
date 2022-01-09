package main

import (
	"image/color"
	"image/png"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()
	//m := melody.NewTable()

	server.POST("/palette", ExtractColorPalette)
	//server.POST("/index", IndexImage)

	err := server.Run(":80")
	CheckError(err)
}

func IndexImage(c *gin.Context) {
	body, err := c.Request.GetBody()
	CheckError(err)

	var buff []byte
	_, err = body.Read(buff)
	CheckError(err)

	if multipart, err := c.FormFile("image"); err == nil {
		file, err := multipart.Open()
		CheckError(err)
		defer CheckError(file.Close())

		pngImg, err := png.Decode(file)
		CheckError(err)

		dups := make(map[color.Color]int)
		for y := 0; y < Y; y++ {
			for x := 0; x < X; x++ {
				col := pngImg.At(x, y)
				if _, ok := dups[col]; !ok {
					pltt.Palette = append(pltt.Palette, col)
					dups[col] = true
				}
			}
		}

	}

}

func ExtractColorPalette(c *gin.Context) {

	if multipart, err := c.FormFile("image"); err == nil {
		file, err := multipart.Open()
		CheckError(err)
		defer CheckError(file.Close())

		pngImg, err := png.Decode(file)
		CheckError(err)

		var pltt ColorPalette
		pltt.Name = multipart.Filename

		X, Y := pngImg.Bounds().Dx(), pngImg.Bounds().Dy()

		dups := make(map[color.Color]bool)
		for y := 0; y < Y; y++ {
			for x := 0; x < X; x++ {
				col := pngImg.At(x, y)
				if _, ok := dups[col]; !ok {
					pltt.Palette = append(pltt.Palette, col)
					dups[col] = true
				}
			}
		}
		c.JSON(200, pltt)
	} else {
		CheckError(err)
	}
}

func CheckError(err error) {
	defer func() {
		if err := recover().(error); err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		panic(err)
	}
}
