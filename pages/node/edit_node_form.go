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

type PageEditNodeForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonEditNode   *components.Button
	buttonDeleteNode *components.Button
	txtHost          *components.TextField
	txtName          *components.TextField
	txtPort          *components.TextField
	nodeInfo         node_manager.NodeInfo

	listStyle material.ListStyle
}

var _ router.Page = &PageEditNodeForm{}

func NewPageEditNodeForm() *PageEditNodeForm {
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

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	buttonEditNode := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "SAVE",
		Icon:            saveIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonEditNode.Label.Alignment = text.Middle
	buttonEditNode.Style.Font.Weight = font.Bold

	txtName := components.NewTextField(th, "Name", "Dero")
	txtHost := components.NewTextField(th, "Host", "node.dero.io")
	txtPort := components.NewTextField(th, "Port", "10102")

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDeleteNode := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "DELETE WALLET",
		Icon:            deleteIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonDeleteNode.Label.Alignment = text.Middle
	buttonDeleteNode.Style.Font.Weight = font.Bold

	return &PageEditNodeForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonEditNode:   buttonEditNode,
		buttonDeleteNode: buttonDeleteNode,
		txtName:          txtName,
		txtHost:          txtHost,
		txtPort:          txtPort,

		listStyle: listStyle,
	}
}

func (p *PageEditNodeForm) IsActive() bool {
	return p.isActive
}

func (p *PageEditNodeForm) Enter() {
	p.isActive = true
	page_instance.header.SetTitle("Edit Node")
	p.animationEnter.Start()
	p.animationLeave.Reset()

	p.txtHost.SetValue(p.nodeInfo.Host)
	p.txtName.SetValue(p.nodeInfo.Name)
	p.txtPort.SetValue(fmt.Sprint(p.nodeInfo.Port))
}

func (p *PageEditNodeForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageEditNodeForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonEditNode.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true)
		} else {
			notification_modals.SuccessInstance.SetText("Success", "new noded added")
			notification_modals.SuccessInstance.SetVisible(true)
		}
	}

	if p.buttonDeleteNode.Clickable.Clicked() {
		err := node_manager.Instance.DelNode(p.nodeInfo.ID)
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true)
		} else {
			notification_modals.SuccessInstance.SetText("Success", "node deleted")
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
			return p.buttonEditNode.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonDeleteNode.Layout(gtx, th)
		},
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageEditNodeForm) submitForm() error {
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

	err = node_manager.Instance.EditNode(node_manager.NodeInfo{
		ID:   p.nodeInfo.ID,
		Name: txtName.Text(),
		Host: txtHost.Text(),
		Port: port,
	})
	if err != nil {
		return err
	}

	return nil
}
