package image_modal

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/router"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type ImageModal struct {
	modal       *components.Modal
	buttonClose *components.Button
	image       *components.Image
	title       string
}

var Instance *ImageModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.N,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded: components.Rounded{
			SW: unit.Dp(10), SE: unit.Dp(10),
			NW: unit.Dp(10), NE: unit.Dp(10),
		},
		Inset:     layout.UniformInset(unit.Dp(20)),
		Animation: components.NewModalAnimationDown(),
		Backdrop:  components.NewModalBackground(),
	})

	closeIcon, _ := widget.NewIcon(icons.NavigationCancel)
	buttonClose := components.NewButton(components.ButtonStyle{
		Icon:           closeIcon,
		TextColor:      color.NRGBA{R: 0, G: 0, B: 0, A: 100},
		HoverTextColor: &color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation:      components.NewButtonAnimationScale(.95),
	})

	Instance = &ImageModal{
		modal:       modal,
		buttonClose: buttonClose,
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout: func(gtx layout.Context, th *material.Theme) {
			Instance.layout(gtx, th)
		},
	})
}

func (r *ImageModal) Open(title string, imgSrc paint.ImageOp) {
	r.title = title
	r.image = &components.Image{
		Src: imgSrc,
		Fit: components.Contain,
	}
	r.modal.SetVisible(true)
}

func (r *ImageModal) layout(gtx layout.Context, th *material.Theme) {
	if r.buttonClose.Clicked() {
		r.modal.SetVisible(false)
	}

	r.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(18), r.title)
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return r.buttonClose.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return r.image.Layout(gtx)
				}),
			)
		})
	})
}
