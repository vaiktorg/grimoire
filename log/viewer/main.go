package main

import (
	"fmt"
	"github.com/vaiktorg/grimoire/log/logger"
	"github.com/vaiktorg/grimoire/log/viewer/src"
	"net/http"
)

type LogViewer struct {
	app.Compo

	//Components
	navbar *src.Navbar
	table  *src.Table

	//Connection Status
	connection logger.Connection
}

// Render ----------------------------------------------
func (l *LogViewer) Render() app.UI {
	return app.Body().Body(

		app.Div().ID("myTabContent").Class("tab-content").Body(
			app.Div().Class("mx-3 tab-pane fade show active").ID("logs").Body(
				l.table,
			),
			app.Div().Class("mx-3 tab-pane fade show").ID("statistics").Body(
				app.H1().Text("STATISTICS"),
			),
		),
	)
}

func main() {
	app.Route("/", &LogViewer{})

	app.RunWhenOnBrowser()

	// Router
	r := http.NewServeMux()

	// FrontEnd Handler
	r.Handle("/", &app.Handler{
		Name:        "Log Viewer",
		Description: "HTTP Log Viewer",
		Styles: []string{
			"https://cdn.jsdelivr.net/npm/bootstrap@5.0.0/dist/css/bootstrap.min.css",
			//"https://cdn.jsdelivr.net/npm/bootswatch@4.5.2/dist/simplex/bootstrap.min.css",
		},
		Scripts: []string{
			"https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js",
		},
		LoadingLabel: "Reading Grimoire...",
		Icon: struct {
			Default    string
			Large      string
			AppleTouch string
		}{
			Default:    "https://i.imgur.com/vNxAhoY.png",
			Large:      "https://i.imgur.com/vNxAhoY.png",
			AppleTouch: "https://i.imgur.com/vNxAhoY.png",
		},
	})

	fmt.Println("Listening on :8080")
	_ = http.ListenAndServe(":8080", r)
}
