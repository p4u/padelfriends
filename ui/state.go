package ui

import "github.com/maxence-charriere/go-app/v10/pkg/app"

const (
	StateGroupName     = "groupName"
	StateGroupPassword = "groupPassword"
	StateGroupID       = "groupID"
)

func SetGroupCredentials(ctx app.Context, groupName, password string) {
	ctx.SetState(StateGroupName, groupName)
	ctx.SetState(StateGroupPassword, password)
}

// New function to set groupID
func SetGroupID(ctx app.Context, groupID string) {
	ctx.SetState(StateGroupID, groupID)
}

func GetGroupCredentials(ctx app.Context) (string, string) {
	var name, pwd string
	ctx.GetState(StateGroupName, &name)
	ctx.GetState(StateGroupPassword, &pwd)
	return name, pwd
}

func GetGroupID(ctx app.Context) string {
	var id string
	ctx.GetState(StateGroupID, &id)
	return id
}
