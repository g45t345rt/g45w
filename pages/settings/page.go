package page_settings

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
)

type Page struct {
	isActive bool
}

var _ router.Container = &Page{}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	app_instance.Current.BottomBar.SetActive("settings")
	p.isActive = true
}

func (p *Page) Leave() {
	p.isActive = false
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return app_instance.Current.BottomBar.Layout(gtx, th)
		}),
	)
}
