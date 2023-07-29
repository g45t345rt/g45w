package main

import (
	"time"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type LogoSplash struct {
	animation *animation.Animation
	image     *components.Image
}

func NewLogoSplash() *LogoSplash {
	animtation := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, 1, ease.Linear),
	))
	animtation.Sequence.SetLoop(-1)

	src, _ := assets.GetImage("dero.jpg")

	image := &components.Image{
		Src: paint.NewImageOp(src),
		Fit: components.Cover,
	}

	return &LogoSplash{
		animation: animtation,
		image:     image,
	}
}

func (l *LogoSplash) Layout(gtx layout.Context) layout.Dimensions {
	r := op.Record(gtx.Ops)
	dims := l.image.Layout(gtx)
	c := r.Stop()

	gtx.Constraints.Min = dims.Size

	{
		state := l.animation.Update(gtx)
		if state.Active {
			defer animation.TransformRotate(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	c.Add(gtx.Ops)
	return dims
}

type LoadState struct {
	status     string
	err        error
	window     *app.Window
	loaded     bool
	logoSplash *LogoSplash
}

func NewLoadState(window *app.Window) *LoadState {
	logoSplash := NewLogoSplash()
	return &LoadState{
		logoSplash: logoSplash,
		window:     window,
	}
}

func (l *LoadState) Complete() {
	l.loaded = true

	l.window.Invalidate()
	time.Sleep(50 * time.Millisecond)
}

func (l *LoadState) SetStatus(status string, err error) {
	if err != nil {
		l.err = err
	} else {
		l.status = status
	}

	l.window.Invalidate()
	time.Sleep(50 * time.Millisecond)
}

func (l *LoadState) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.UniformInset(unit.Dp(30)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if l.err != nil {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(20), l.status)
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), l.err.Error())
						return lbl.Layout(gtx)
					}),
				)
			} else {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.X = gtx.Dp(100)
							gtx.Constraints.Max.Y = gtx.Dp(100)
							return l.logoSplash.Layout(gtx)
						})
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(40)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(20), l.status)
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						})
					}),
				)
			}
		})
	})
}
