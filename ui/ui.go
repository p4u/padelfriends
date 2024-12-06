package ui

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Root layout component if needed
type Root struct {
	app.Compo
}

func (r *Root) Render() app.UI {
	return app.Div().Body(
		&NavBar{},
		app.Div().Body(
			app.Text("Welcome to Padel Friends! Use the navigation above."),
		),
	)
}

func init() {
	app.Route("/", func() app.Composer { return &HomePage{} })
	app.Route("/players", func() app.Composer { return &PlayersPage{} })
	app.Route("/matches", func() app.Composer { return &MatchesPage{} })
	app.Route("/stats", func() app.Composer { return &StatsPage{} })

	app.RunWhenOnBrowser()
}
