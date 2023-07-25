package page_settings

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
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageEditIPFSGateway struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonEdit   *components.Button
	buttonDelete *components.Button
	txtEndpoint  *components.TextField
	txtName      *components.TextField
	gateway      app_data.IPFSGateway

	confirmDelete *components.Confirm

	list *widget.List
}

var _ router.Page = &PageEditIPFSGateway{}

func NewPageEditIPFSGateway() *PageEditIPFSGateway {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.OutCubic),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonEdit := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            saveIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		LoadingIcon:     loadingIcon,
	})
	buttonEdit.Label.Alignment = text.Middle
	buttonEdit.Style.Font.Weight = font.Bold

	txtName := components.NewTextField()
	txtEndpoint := components.NewTextField()

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDelete := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            deleteIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonDelete.Label.Alignment = text.Middle
	buttonDelete.Style.Font.Weight = font.Bold

	confirmDelete := components.NewConfirm(layout.Center)
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			confirmDelete.Prompt = lang.Translate("Are you sure?")
			confirmDelete.NoText = lang.Translate("NO")
			confirmDelete.YesText = lang.Translate("YES")
			confirmDelete.Layout(gtx, th)
		},
	})

	return &PageEditIPFSGateway{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonEdit:   buttonEdit,
		buttonDelete: buttonDelete,
		txtName:      txtName,
		txtEndpoint:  txtEndpoint,

		confirmDelete: confirmDelete,

		list: list,
	}
}

func (p *PageEditIPFSGateway) IsActive() bool {
	return p.isActive
}

func (p *PageEditIPFSGateway) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Edit IPFS Gateway"))
	p.animationEnter.Start()
	p.animationLeave.Reset()

	p.txtEndpoint.SetValue(p.gateway.Endpoint)
	p.txtName.SetValue(p.gateway.Name)
}

func (p *PageEditIPFSGateway) Leave() {
	if page_instance.header.IsHistory(PAGE_EDIT_IPFS_GATEWAY) {
		p.animationEnter.Reset()
		p.animationLeave.Start()
	} else {
		p.isActive = false
	}
}

func (p *PageEditIPFSGateway) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonEdit.Clicked() {
		p.submitForm(gtx)
	}

	if p.buttonDelete.Clicked() {
		p.confirmDelete.SetVisible(true)
	}

	if p.confirmDelete.ClickedYes() {
		err := p.removeGateway()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {

			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Gateway deleted"))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			page_instance.header.GoBack()
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "Dero NFTs")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th, lang.Translate("Endpoint"), "wss://node.deronfts.com/ws")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonEdit.Text = lang.Translate("SAVE")
			return p.buttonEdit.Layout(gtx, th)
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
			p.buttonDelete.Text = lang.Translate("DELETE NODE")
			return p.buttonDelete.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageEditIPFSGateway) removeGateway() error {
	endpoint := p.gateway.Endpoint
	err := app_data.DelIPFSGateway(0)
	if err != nil {
		return err
	}

	if node_manager.CurrentNode != nil {
		if node_manager.CurrentNode.Endpoint == endpoint {
			node_manager.CurrentNode = nil
			walletapi.Connected = false

			settings.App.NodeEndpoint = ""
			err := settings.Save()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PageEditIPFSGateway) submitForm(gtx layout.Context) {
	p.buttonEdit.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonEdit.SetLoading(false)
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		txtName := p.txtName.Editor()
		txtEnpoint := p.txtEndpoint.Editor()

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

		err = app_data.UpdateIPFSGatway(app_data.IPFSGateway{
			Name:     txtName.Text(),
			Endpoint: txtEnpoint.Text(),
		})
		if err != nil {
			setError(err)
			return
		}

		p.buttonEdit.SetLoading(false)
		notification_modals.SuccessInstance.SetText("Success", lang.Translate("Data saved"))
		notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		page_instance.header.GoBack()
	}()
}
