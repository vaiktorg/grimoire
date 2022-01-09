package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/vaiktorg/grimoire/heart"
	"github.com/vaiktorg/grimoire/helpers"
)

// Client
func main() {
	hrt := heart.NewHeartbeat("ServiceX", "http://localhost:8081")
	router := http.NewServeMux()
	router.HandleFunc("/hb", hrt.ServiceHandler())

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	helpers.ConsoleCloser("Server", func() {
		err := srv.Shutdown(context.Background())
		fmt.Println(err)
	}, nil)

	// service connections
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
