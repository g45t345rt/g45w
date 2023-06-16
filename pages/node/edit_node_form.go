package page_node

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/walletapi"
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
	txtEndpoint      *components.TextField
	txtName          *components.TextField
	nodeConn         node_manager.NodeConnection
	submitting       bool

	confirmDelete *components.Confirm

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

	txtName := components.NewTextField(th, "Name", "Dero NFTs")
	txtEndpoint := components.NewTextField(th, "Host", "wss://node.deronfts.com/ws")

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

	w := app_instance.Window
	confirmDelete := components.NewConfirm(w, "Are you sure?", th, layout.Center)
	app_instance.Router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		confirmDelete.Layout(gtx, th)
	})

	return &PageEditNodeForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonEditNode:   buttonEditNode,
		buttonDeleteNode: buttonDeleteNode,
		txtName:          txtName,
		txtEndpoint:      txtEndpoint,

		confirmDelete: confirmDelete,

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

	p.txtEndpoint.SetValue(p.nodeConn.Endpoint)
	p.txtName.SetValue(p.nodeConn.Name)
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
		p.submitForm()
	}

	if p.buttonDeleteNode.Clickable.Clicked() {
		p.confirmDelete.SetVisible(true)
	}

	if p.confirmDelete.ClickedYes() {
		err := node_manager.Instance.DelNode(p.nodeConn.ID)
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true)
		} else {
			notification_modals.SuccessInstance.SetText("Success", "node deleted")
			notification_modals.SuccessInstance.SetVisible(true)
			page_instance.router.SetCurrent(PAGE_SELECT_NODE)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonEditNode.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			max := image.Pt(gtx.Dp(unit.Dp(gtx.Constraints.Max.X)), 5)
			paint.FillShape(gtx.Ops, color.NRGBA{A: 150}, clip.Rect{
				Min: image.Pt(0, 0),
				Max: max,
			}.Op())
			return layout.Dimensions{Size: max}
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

func (p *PageEditNodeForm) submitForm() {
	if p.submitting {
		return
	}

	p.submitting = true

	go func() {
		setError := func(err error) {
			p.submitting = false
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true)
		}

		txtName := p.txtName.EditorStyle.Editor
		txtEnpoint := p.txtEndpoint.EditorStyle.Editor

		if txtName.Text() == "" {
			setError(fmt.Errorf("enter name"))
			return
		}

		if txtEnpoint.Text() == "" {
			setError(fmt.Errorf("enter endpoint"))
			return
		}

		_, err := walletapi.TestConnect(txtEnpoint.Text())
		if err != nil {
			setError(err)
			return
		}

		err = node_manager.Instance.EditNode(node_manager.NodeConnection{
			ID:       p.nodeConn.ID,
			Name:     txtName.Text(),
			Endpoint: txtEnpoint.Text(),
		})
		if err != nil {
			setError(err)
			return
		}

		p.submitting = false
		notification_modals.SuccessInstance.SetText("Success", "data saved")
		notification_modals.SuccessInstance.SetVisible(true)
		page_instance.router.SetCurrent(PAGE_SELECT_NODE)
	}()

}
