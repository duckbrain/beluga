package grifts

import (
	"github.com/duckbrain/beluga/belugad/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
