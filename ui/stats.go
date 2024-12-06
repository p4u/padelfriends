package ui

import (
	"encoding/json"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type PlayerStats struct {
	PlayerID string  `json:"player_id"`
	Name     string  `json:"name"`
	Matches  int     `json:"matches"`
	Wins     int     `json:"wins"`
	Losses   int     `json:"losses"`
	Points   int     `json:"points"`
	WinRatio float64 `json:"win_ratio"`
}

type StatsPage struct {
	app.Compo

	groupName string
	password  string
	groupID   string

	stats    []PlayerStats
	errorMsg string
	loading  bool
}

func (s *StatsPage) OnMount(ctx app.Context) {
	s.groupName, s.password = GetGroupCredentials(ctx)
	s.groupID = GetGroupID(ctx)
	if s.groupID != "" && s.password != "" {
		s.loadStats(ctx)
	}
}

func (s *StatsPage) Render() app.UI {
	if s.groupName == "" || s.password == "" {
		return app.Div().Body(
			app.H2().Text("Please set group credentials on Home page."),
		)
	}
	if s.groupID == "" {
		return app.Div().Body(
			app.H2().Text("Group is not initialized yet."),
		)
	}

	return app.Div().Body(
		app.H1().Text("Statistics"),
		app.If(s.errorMsg != "", func() app.UI {
			return app.Div().Style("color", "red").Text(s.errorMsg)
		}),
		app.If(s.loading, func() app.UI {
			return app.Div().Text("Loading stats...")
		}),
		app.If(!s.loading && len(s.stats) > 0, func() app.UI {
			return app.Table().Body(
				app.Tr().Body(
					app.Th().Text("Name"),
					app.Th().Text("Matches"),
					app.Th().Text("Wins"),
					app.Th().Text("Losses"),
					app.Th().Text("Points"),
					app.Th().Text("Win Ratio"),
				),
				app.Range(s.stats).Slice(func(i int) app.UI {
					ps := s.stats[i]
					return app.Tr().Body(
						app.Td().Text(ps.Name),
						app.Td().Text(ps.Matches),
						app.Td().Text(ps.Wins),
						app.Td().Text(ps.Losses),
						app.Td().Text(ps.Points),
						app.Td().Textf("%.2f%%", ps.WinRatio*100),
					)
				}),
			)
		}),
		app.Button().Text("Refresh").OnClick(s.onRefreshClick),
		app.Button().Text("Export CSV").OnClick(s.onExportCSVClick), // We'll add CSV logic later
	)
}

func (s *StatsPage) loadStats(ctx app.Context) {
	s.loading = true
	s.errorMsg = ""
	url := "/api/group/" + s.groupID + "/statistics?password=" + s.password

	ctx.Async(func() {
		respPromise := app.Window().Call("fetch", url)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				s.setError(ctx, "HTTP error: "+resp.Get("statusText").String())
				return nil
			}
			return resp.Call("json").Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				data := args[0]
				jsonStr := data.String()
				var stats []PlayerStats
				if err := json.Unmarshal([]byte(jsonStr), &stats); err != nil {
					s.setError(ctx, "Error parsing stats: "+err.Error())
				} else {
					ctx.Dispatch(func(ctx app.Context) {
						s.stats = stats
						s.loading = false
					})
				}
				return nil
			}))
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			errMsg := args[0].String()
			s.setError(ctx, errMsg)
			return nil
		}))
	})
}

func (s *StatsPage) onRefreshClick(ctx app.Context, e app.Event) {
	s.loadStats(ctx)
}

func (s *StatsPage) setError(ctx app.Context, msg string) {
	ctx.Dispatch(func(ctx app.Context) {
		s.errorMsg = msg
		s.loading = false
	})
}

func (s *StatsPage) onExportCSVClick(ctx app.Context, e app.Event) {
	if s.groupID == "" || s.password == "" {
		s.setError(ctx, "Group credentials not set.")
		return
	}

	url := "/api/group/" + s.groupID + "/statistics/csv?password=" + s.password

	ctx.Async(func() {
		respPromise := app.Window().Call("fetch", url)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				s.setError(ctx, "HTTP error: "+resp.Get("statusText").String())
				return nil
			}
			return resp.Call("blob").Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				blob := args[0]
				blobURL := app.Window().Get("URL").Call("createObjectURL", blob)

				document := app.Window().Get("document") // Get the DOM document
				link := document.Call("createElement", "a")
				link.Set("href", blobURL.String())
				link.Set("download", "stats.csv")

				body := document.Get("body")
				body.Call("appendChild", link)
				link.Call("click")

				body.Call("removeChild", link)
				app.Window().Get("URL").Call("revokeObjectURL", blobURL)
				return nil
			}))
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			errMsg := args[0].String()
			s.setError(ctx, "Error exporting CSV: "+errMsg)
			return nil
		}))
	})
}
