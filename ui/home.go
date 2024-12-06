package ui

import (
	"encoding/json"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Group struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
}

type HomePage struct {
	app.Compo

	groupName string
	password  string
	groupID   string
	errorMsg  string
	groups    []Group
}

func (h *HomePage) OnMount(ctx app.Context) {
	h.groupName, h.password = GetGroupCredentials(ctx)
	h.fetchGroups(ctx)
}

func (h *HomePage) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("Padel Friends"),

		app.H2().Text("Existing Groups"),
		app.If(h.errorMsg != "", func() app.UI {
			return app.Div().Style("color", "red").Text(h.errorMsg)
		}).Else(func() app.UI {
			// No error, show groups or a "No groups" message
			if len(h.groups) == 0 {
				return app.P().Text("No groups available")
			}
			return app.Ul().Body(
				app.Range(h.groups).Slice(func(i int) app.UI {
					g := h.groups[i]
					return app.Li().Body(
						app.Text(g.Name),
						app.Button().Text("Use this Group").OnClick(func(ctx app.Context, e app.Event) {
							SetGroupCredentials(ctx, g.Name, "")
							h.groupName = g.Name
							h.groupID = g.ID
							ctx.Reload()
						}),
					)
				}),
			)
		}),

		app.H2().Text("Register a New Group"),
		app.Div().Body(
			app.Label().Text("Group Name:"),
			app.Input().Type("text").ID("group-name"),
		),
		app.Div().Body(
			app.Label().Text("Password:"),
			app.Input().Type("password").ID("group-password"),
		),
		app.Button().Text("Register").OnClick(h.onRegisterClick),
	)
}

func (h *HomePage) onRegisterClick(ctx app.Context, e app.Event) {
	document := app.Window().Get("document")
	if !document.Truthy() {
		// Document not available, do nothing
		return
	}

	name := document.Call("getElementById", "group-name").Get("value").String()
	pwd := document.Call("getElementById", "group-password").Get("value").String()

	SetGroupCredentials(ctx, name, pwd)
	h.fetchGroupByName(ctx, name, pwd)
}

func (h *HomePage) fetchGroups(ctx app.Context) {
	ctx.Async(func() {
		url := "/api/groups"
		respPromise := app.Window().Call("fetch", url)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				h.setErrorMsg(ctx, "Error fetching groups: "+resp.Get("statusText").String())
				return nil
			}
			return resp.Call("json").Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
				data := args[0]
				jsonStr := data.String()
				var groups []Group
				if err := json.Unmarshal([]byte(jsonStr), &groups); err != nil {
					h.setErrorMsg(ctx, "Error parsing groups: "+err.Error())
				} else {
					ctx.Dispatch(func(ctx app.Context) {
						h.groups = groups
						h.errorMsg = ""
					})
				}
				return nil
			}))
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) any {
			errMsg := args[0].String()
			h.setErrorMsg(ctx, "Error fetching groups: "+errMsg)
			return nil
		}))
	})
}

func (h *HomePage) fetchGroupByName(ctx app.Context, name, pwd string) {
	ctx.Async(func() {
		url := "/api/group/byname/" + name + "?password=" + pwd
		respPromise := app.Window().Call("fetch", url)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				status := resp.Get("status").Int()
				if status == http.StatusNotFound {
					// Group not found, create it
					h.createGroup(ctx, name, pwd)
					return nil
				}
				h.setErrorMsg(ctx, "Error retrieving group: "+resp.Get("statusText").String())
				return nil
			}
			return resp.Call("json").Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
				data := args[0]
				jsonStr := data.String()
				var g Group
				if err := json.Unmarshal([]byte(jsonStr), &g); err != nil {
					h.setErrorMsg(ctx, "Error parsing group data: "+err.Error())
				} else {
					ctx.Dispatch(func(ctx app.Context) {
						SetGroupID(ctx, g.ID)
						h.groupID = g.ID
					})
				}
				return nil
			}))
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) any {
			errMsg := args[0].String()
			h.setErrorMsg(ctx, "Error fetching group: "+errMsg)
			return nil
		}))
	})
}

func (h *HomePage) createGroup(ctx app.Context, name, pwd string) {
	ctx.Async(func() {
		url := "/api/group"
		payload := map[string]interface{}{
			"name":     name,
			"password": pwd,
		}
		body, _ := json.Marshal(payload)
		// Make sure headers is map[string]interface{}
		headers := map[string]interface{}{
			"Content-Type": "application/json",
		}
		opts := map[string]interface{}{
			"method":  "POST",
			"headers": headers,
			"body":    string(body),
		}

		respPromise := app.Window().Call("fetch", url, opts)
		respPromise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
			resp := args[0]
			if !resp.Get("ok").Bool() {
				h.setErrorMsg(ctx, "Error creating group: "+resp.Get("statusText").String())
				return nil
			}
			return resp.Call("json").Call("then", app.FuncOf(func(this app.Value, args []app.Value) any {
				data := args[0]
				jsonStr := data.String()
				var g Group
				if err := json.Unmarshal([]byte(jsonStr), &g); err != nil {
					h.setErrorMsg(ctx, "Error parsing group creation response: "+err.Error())
					return nil
				}

				ctx.Dispatch(func(ctx app.Context) {
					SetGroupID(ctx, g.ID)
					h.groupID = g.ID
				})
				return nil
			}))
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) any {
			errMsg := args[0].String()
			h.setErrorMsg(ctx, "Error creating group: "+errMsg)
			return nil
		}))
	})
}

func (h *HomePage) setErrorMsg(ctx app.Context, msg string) {
	ctx.Dispatch(func(ctx app.Context) {
		h.errorMsg = msg
	})
}
