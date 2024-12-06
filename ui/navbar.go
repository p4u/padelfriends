package ui

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// NavBar provides a simple top navigation menu.
type NavBar struct {
	app.Compo
}

func (n *NavBar) Render() app.UI {
	return app.Nav().Body(
		app.A().Href("/").Text("Home").Style("margin-right", "20px"),
		app.A().Href("/players").Text("Players").Style("margin-right", "20px"),
		app.A().Href("/matches").Text("Matches").Style("margin-right", "20px"),
		app.A().Href("/stats").Text("Stats"),
	)
}
