package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/vaiktorg/grimoire/helpers"
)

var (
	sampleText = "It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layouts. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like)."
)

func main() {
	ticker := helpers.NewVariableTicker(80 * time.Millisecond)
	//ticker := helpers.NewRandomTicker(1*time.Millisecond, 500*time.Millisecond)
	//words := strings.Split(sampleText, " ")
	accum := 0
	for range ticker.C {
		if accum == len(sampleText) {
			return
		}
		fmt.Printf("%+v", string(sampleText[accum]))
		if strings.ContainsAny(string(sampleText[accum]), ".") {
			time.Sleep(time.Second)
		}
		accum++
	}
}
