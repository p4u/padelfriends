package ui

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type MatchesPage struct {
	app.Compo

	groupName   string
	password    string
	groupID     string
	players     []Player
	selectedIDs map[string]bool

	errorMsg   string
	successMsg string

	matchTeams  string
	lastMatchID string // Store last created match ID

	scoreTeam1 string
	scoreTeam2 string
}

func (m *MatchesPage) OnMount(ctx app.Context) {
	m.groupName, m.password = GetGroupCredentials(ctx)
	m.groupID = GetGroupID(ctx)
	m.selectedIDs = make(map[string]bool)
	m.loadPlayers(ctx)
}

func (m *MatchesPage) Render() app.UI {
	if m.groupName == "" || m.password == "" {
		return app.Div().Body(
			app.H2().Text("Please set group credentials on Home page."),
		)
	}
	if m.groupID == "" {
		return app.Div().Body(
			app.H2().Text("Group is not initialized yet, please wait or set it on Home page."),
		)
	}

	return app.Div().Body(
		app.H1().Text("Matches"),
		app.If(m.errorMsg != "", func() app.UI {
			return app.Div().Style("color", "red").Text(m.errorMsg)
		}),
		app.If(m.successMsg != "", func() app.UI {
			return app.Div().Style("color", "green").Text(m.successMsg)
		}),

		app.H2().Text("Select 4 Players"),
		app.Ul().Body(
			app.Range(m.players).Slice(func(i int) app.UI {
				p := m.players[i]
				return app.Li().Body(
					app.Input().Type("checkbox").
						Value(p.ID).
						Checked(m.selectedIDs[p.ID]).
						OnChange(m.onPlayerCheck),
					app.Text(p.Name),
				)
			}),
		),
		app.Button().Text("Create Match").OnClick(m.onCreateMatchClick),

		app.If(m.matchTeams != "", func() app.UI {
			// Display teams and form to submit results
			return app.Div().Body(
				app.H2().Text("Match Created"),
				app.P().Text("Teams: "+m.matchTeams),
				app.H3().Text("Submit Results"),
				app.Div().Body(
					app.Label().Text("Team 1 Score: "),
					app.Input().Type("number").Value(m.scoreTeam1).ID("score-team1").OnChange(m.onScoreChange),
				),
				app.Div().Body(
					app.Label().Text("Team 2 Score: "),
					app.Input().Type("number").Value(m.scoreTeam2).ID("score-team2").OnChange(m.onScoreChange),
				),
				app.Button().Text("Submit Results").OnClick(m.onSubmitResultsClick),
			)
		}),
	)
}

func (m *MatchesPage) loadPlayers(ctx app.Context) {
	m.errorMsg = ""
	m.successMsg = ""

	if m.groupID == "" {
		return
	}

	url := "/api/group/" + m.groupID + "/players?password=" + m.password
	m.fetchGET(ctx, url, &m.players, func(err error) {
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				m.errorMsg = "Error loading players: " + err.Error()
			}
		})
	})
}

func (m *MatchesPage) onPlayerCheck(ctx app.Context, e app.Event) {
	val := e.Value.Get("value").String()
	checked := e.Value.Get("checked").Bool()

	ctx.Dispatch(func(ctx app.Context) {
		if checked {
			m.selectedIDs[val] = true
		} else {
			delete(m.selectedIDs, val)
		}
	})
}

func (m *MatchesPage) onCreateMatchClick(ctx app.Context, e app.Event) {
	if len(m.selectedIDs) != 4 {
		ctx.Dispatch(func(ctx app.Context) {
			m.errorMsg = "You must select exactly 4 players"
		})
		return
	}

	var playerIDs []string
	for pid := range m.selectedIDs {
		playerIDs = append(playerIDs, pid)
	}

	payload := map[string]interface{}{
		"player_ids": playerIDs,
	}
	body, _ := json.Marshal(payload)

	url := "/api/group/" + m.groupID + "/matches?password=" + m.password

	// Create a struct to parse response with match_id
	type matchResponse struct {
		Match   map[string]interface{} `json:"match"`
		Details struct {
			Team1      []string `json:"team1"`
			Team2      []string `json:"team2"`
			ScoreTeam1 int      `json:"score_team1"`
			ScoreTeam2 int      `json:"score_team2"`
		} `json:"details"`
	}

	var respData matchResponse
	m.fetchPOST(ctx, url, body, &respData, func(err error) {
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				m.errorMsg = "Failed to create match: " + err.Error()
				return
			}
			// Extract match_id from respData.Match["id"] if it exists
			if idVal, ok := respData.Match["_id"]; ok {
				if idStr, ok := idVal.(string); ok {
					m.lastMatchID = idStr
				}
			}
			// Construct teams string
			teamsStr := "Team1: " + joinStrings(respData.Details.Team1) + " | Team2: " + joinStrings(respData.Details.Team2)
			m.matchTeams = teamsStr
			m.scoreTeam1 = "0"
			m.scoreTeam2 = "0"
		})
	})
}

func (m *MatchesPage) onScoreChange(ctx app.Context, e app.Event) {
	id := e.Value.Get("id").String()
	val := e.Value.Get("value").String()

	ctx.Dispatch(func(ctx app.Context) {
		if id == "score-team1" {
			m.scoreTeam1 = val
		} else if id == "score-team2" {
			m.scoreTeam2 = val
		}
	})
}

func (m *MatchesPage) onSubmitResultsClick(ctx app.Context, e app.Event) {
	team1Score, err1 := strconv.Atoi(m.scoreTeam1)
	team2Score, err2 := strconv.Atoi(m.scoreTeam2)

	if err1 != nil || err2 != nil {
		ctx.Dispatch(func(ctx app.Context) {
			m.errorMsg = "Invalid scores. Must be integers."
		})
		return
	}

	payload := map[string]int{
		"score_team1": team1Score,
		"score_team2": team2Score,
	}
	body, _ := json.Marshal(payload)

	url := "/api/group/" + m.groupID + "/matches/" + m.lastMatchID + "/results?password=" + m.password

	m.fetchPOST(ctx, url, body, nil, func(err error) {
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				m.errorMsg = "Failed to submit results: " + err.Error()
				return
			}
			m.successMsg = "Results submitted successfully!"
			// Clear selection or prompt user for next match
			m.matchTeams = ""
			m.lastMatchID = ""
			m.selectedIDs = make(map[string]bool)
			m.scoreTeam1 = ""
			m.scoreTeam2 = ""
		})
	})
}

func (m *MatchesPage) fetchPOST(ctx app.Context, url string, body []byte, out interface{}, cb func(err error)) {
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
			} else {
				cb(nil)
				return nil
			}
		})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			errMsg := args[0].String()
			cb(errors.New(errMsg))
			return nil
		}))
	})
}

func (m *MatchesPage) fetchGET(ctx app.Context, url string, out interface{}, cb func(err error)) {
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

func joinStrings(arr []string) string {
	out := ""
	for i, s := range arr {
		if i > 0 {
			out += ", "
		}
		out += s
	}
	return out
}
