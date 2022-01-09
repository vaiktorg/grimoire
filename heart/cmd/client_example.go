package main

import (
	"fmt"

	"github.com/vaiktorg/grimoire/heart"
)

// Server
func main() {

	monitor := heart.NewMonitor()
	service := monitor.AddService("http://localhost:8081/hb", func(heartbeat *heart.Heartbeat, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(heartbeat)
	})

	service.Ping()
}
