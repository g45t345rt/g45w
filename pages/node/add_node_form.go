package page_node

import (
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAddNodeForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonAddNode *components.Button
	txtHost       *components.TextField
	txtName       *components.TextField
	txtPort       *components.TextField

	listStyle material.ListStyle
}

var _ router.Page = &PageAddNodeForm{}

func NewPageAddNodeForm() *PageAddNodeForm {
	th := app_instance.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .5, ease.OutCubic),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAddNode := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "ADD NODE",
		Icon:            addIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonAddNode.Label.Alignment = text.Middle
	buttonAddNode.Style.Font.Weight = font.Bold

	txtName := components.NewTextField(th, "Name", "Dero")
	txtHost := components.NewTextField(th, "Host", "node.dero.io")
	txtPort := components.NewTextField(th, "Port", "10102")

	return &PageAddNodeForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonAddNode: buttonAddNode,
		txtName:       txtName,
		txtHost:       txtHost,
		txtPort:       txtPort,

		listStyle: listStyle,
	}
}

func (p *PageAddNodeForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddNodeForm) Enter() {
	p.isActive = true
	page_instance.header.SetTitle("Add Node")
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageAddNodeForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageAddNodeForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	if p.buttonAddNode.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true)
		} else {
			notification_modals.SuccessInstance.SetText("Success", "new noded added")
			notification_modals.SuccessInstance.SetVisible(true)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtHost.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtPort.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonAddNode.Layout(gtx, th)
		},
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddNodeForm) submitForm() error {
	txtName := p.txtName.EditorStyle.Editor
	txtHost := p.txtHost.EditorStyle.Editor
	txtPort := p.txtPort.EditorStyle.Editor

	if txtName.Text() == "" {
		return fmt.Errorf("enter name")
	}

	if txtHost.Text() == "" {
		return fmt.Errorf("enter host")
	}

	if txtPort.Text() == "" {
		return fmt.Errorf("enter port")
	}

	port, err := strconv.ParseUint(txtPort.Text(), 10, 64)
	if err != nil {
		return err
	}

	err = node_manager.Instance.AddNode(node_manager.NodeInfo{
		Name: txtName.Text(),
		Host: txtHost.Text(),
		Port: port,
	})
	if err != nil {
		return err
	}

	return nil
}
