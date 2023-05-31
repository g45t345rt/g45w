package page_settings

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
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
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), "App Dir")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			appDir := settings.Instance.AppDir
			label := material.Label(th, unit.Sp(16), appDir)
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), "Version")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), settings.Version)
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), "Git Version")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), settings.GitVersion)
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), "Build Time")
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(th, unit.Sp(16), settings.BuildTime)
			return label.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return app_instance.Current.BottomBar.Layout(gtx, th)
		}),
	)
}
