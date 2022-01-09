package src

type Navbar struct {
	app.Compo
}

func (n *Navbar) Render() app.UI {
	return app.Nav().Class("models models-expand-lg models-light bg-light").Body(
		app.Img().Class("models-item mx-2").Src("https://i.imgur.com/vNxAhoY.png").Style("width", "2em"),
		app.A().Class("models-brand").Href("/").Text("VAIKTORG VIEWER"),
		app.Button().Class("models-toggler").
			Type("button").
			DataSet("bs-toggle", "collapse").
			DataSet("bs-target", "#navbarColor03").
			Aria("controls", "navbarColor03").
			Aria("expanded", "false").
			Aria("lable", "Toggle Navigation").Body(
			app.Span().Class("models-toggler-icon"),
		),
		//-------------------------------------------------------
		app.Div().Class("collapse models-collapse").ID("navbarColor03").Body(
			app.Ul().Class("nav nav-tabs models-nav me-auto").Body(
				// Logs Tab Button
				app.Li().Class("nav-item").Body(
					app.A().Class("nav-link active").
						DataSet("bs-toggle", "tab").
						Href("#logs").
						Text("Logs"),
				),
				// Statistics Tab Button
				app.Li().Class("nav-item").Body(
					app.A().Class("nav-link").
						DataSet("bs-toggle", "tab").
						Href("#statistics").
						Text("Statistics"),
				),
				// Filters Tab Button
				app.Li().Class("nav-item dropdown").Body(
					app.A().Class("nav-link dropdown-toggle").
						DataSet("bs-toggle", "dropdown").
						Aria("haspopup", "true").
						Aria("expanded", "false").
						Href("#").
						Text("Filter"),
					app.Div().Class("dropdown-menu").Body(
						app.A().Class("dropdown-item text-black-50").Href("#").Text("Trace"),
						app.A().Class("dropdown-item text-warning").Href("#").Text("Debug"),
						app.A().Class("dropdown-item text-info").Href("#").Text("Info"),
						app.A().Class("dropdown-item text-dark").Href("#").Text("Warn"),
						app.A().Class("dropdown-item text-danger").Href("#").Text("Error"),
						app.A().Class("dropdown-item text-danger fw-bold").Href("#").Text("Fatal"),
						app.Div().Class("dropdown-divider"),
						app.A().Class("dropdown-item").Href("#").Text("Show All"),
					),
				),
			),
			app.Form().Class("d-flex").Body(
				app.Input().Class("form-control me-sm-2").Type("text").Placeholder("Search"),
				app.Button().Class("btn btn-light my-2 my-sm-0").Type("submit").Text("Search"),
			),
		),
	)
}
