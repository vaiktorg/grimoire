package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vaiktorg/grimoire/errs"
)

func main() {
	engine := gin.Default()
	engine.StaticFS("/", gin.Dir("dist", false))
	errs.Must(engine.Run(":8080"))
}
