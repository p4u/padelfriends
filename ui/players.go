package ui

import (
	"encoding/json"
	"errors"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Player struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type PlayersPage struct {
	app.Compo

	players    []Player
	groupName  string
	password   string
	errorMsg   string
	successMsg string
}

func (p *PlayersPage) OnMount(ctx app.Context) {
	p.groupName, p.password = GetGroupCredentials(ctx)
	p.loadPlayers(ctx)
}

func (p *PlayersPage) Render() app.UI {
	if p.groupName == "" || p.password == "" {
		return app.Div().Body(
			app.H2().Text("Please set group credentials on Home page."),
		)
	}

	return app.Div().Body(
		app.H1().Text("Players"),
		// Wrap UI in a func() app.UI for app.If
		app.If(p.errorMsg != "", func() app.UI {
			return app.Div().Style("color", "red").Text(p.errorMsg)
		}),
		app.If(p.successMsg != "", func() app.UI {
			return app.Div().Style("color", "green").Text(p.successMsg)
		}),

		app.H2().Text("Add Player"),
		app.Div().Body(
			app.Label().Text("Name:"),
			app.Input().Type("text").ID("player-name"),
			app.Button().Text("Add").OnClick(p.onAddClick),
		),

		app.H2().Text("Current Players"),
		app.Ul().Body(
			app.Range(p.players).Slice(func(i int) app.UI {
				player := p.players[i]
				return app.Li().Text(player.Name)
			}),
		),
	)
}

func (p *PlayersPage) loadPlayers(ctx app.Context) {
	p.errorMsg = ""
	p.successMsg = ""

	groupID := GetGroupID(ctx)
	if groupID == "" {
		p.errorMsg = "Group not initialized. Please set credentials on Home page."
		return
	}

	url := "/api/group/" + groupID + "/players?password=" + p.password
	p.fetchGET(ctx, url, &p.players, func(err error) {
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				p.errorMsg = "Error loading players: " + err.Error()
			}
		})
	})
}

func (p *PlayersPage) onAddClick(ctx app.Context, e app.Event) {
	name := ctx.JSSrc().Get("document").Call("getElementById", "player-name").Get("value").String()
	if name == "" {
		ctx.Dispatch(func(ctx app.Context) {
			p.errorMsg = "Name cannot be empty"
		})
		return
	}

	groupID := GetGroupID(ctx)
	if groupID == "" {
		ctx.Dispatch(func(ctx app.Context) {
			p.errorMsg = "Group not initialized."
		})
		return
	}

	payload := map[string]string{"name": name}
	body, _ := json.Marshal(payload)

	url := "/api/group/" + groupID + "/players?password=" + p.password
	p.fetchPOST(ctx, url, body, nil, func(err error) {
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				p.errorMsg = "Failed to add player: " + err.Error()
			} else {
				p.successMsg = "Player added successfully"
			}
		})
		// Reload players after adding
		p.loadPlayers(ctx)
	})
}

func (p *PlayersPage) getGroupIDPlaceholder() string {
	return "63f1b1d7c3dbaf19d0d8475c" // fake ObjectID-like string
}

// fetchGET performs a GET request using fetch and unmarshals the JSON response into out.
func (p *PlayersPage) fetchGET(ctx app.Context, url string, out interface{}, cb func(err error)) {
	ctx.Async(func() {
		respPromise := app.Window().Call("fetch", url)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				cb(errors.New("HTTP error: " + resp.Get("statusText").String()))
				return nil
			}
			return resp.Call("json").Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				data := args[0]
				jsonStr := data.String()
				if err := json.Unmarshal([]byte(jsonStr), out); err != nil {
					cb(err)
				} else {
					cb(nil)
				}
				return nil
			}))
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			errMsg := args[0].String()
			cb(errors.New(errMsg))
			return nil
		}))
	})
}

// fetchPOST performs a POST request using fetch with a JSON payload.
func (p *PlayersPage) fetchPOST(ctx app.Context, url string, body []byte, out interface{}, cb func(err error)) {
	ctx.Async(func() {
		opts := map[string]interface{}{
			"method": "POST",
			"headers": map[string]string{
				"Content-Type": "application/json",
			},
			"body": string(body),
		}
		respPromise := app.Window().Call("fetch", url, opts)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				cb(errors.New("HTTP error: " + resp.Get("statusText").String()))
				return nil
			}
			if out != nil {
				return resp.Call("json").Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
					data := args[0]
					jsonStr := data.String()
					if err := json.Unmarshal([]byte(jsonStr), out); err != nil {
						cb(err)
					} else {
						cb(nil)
					}
					return nil
				}))
			}
			// No output expected
			cb(nil)
			return nil
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			errMsg := args[0].String()
			cb(errors.New(errMsg))
			return nil
		}))
	})
}
