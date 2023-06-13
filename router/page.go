package router

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

type Page interface {
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
	IsActive() bool
	Enter()
	Leave()
}
