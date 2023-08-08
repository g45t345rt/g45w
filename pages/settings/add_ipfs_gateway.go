package page_settings

import (
	"fmt"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAddIPFSGateway struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonAdd   *components.Button
	txtEndpoint *components.TextField
	txtName     *components.TextField

	list *widget.List
}

var _ router.Page = &PageAddIPFSGateway{}

func NewPageAddIPFSGateway() *PageAddIPFSGateway {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            addIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		LoadingIcon:     loadingIcon,
	})
	buttonAdd.Label.Alignment = text.Middle
	buttonAdd.Style.Font.Weight = font.Bold

	txtName := components.NewTextField()
	txtEndpoint := components.NewTextField()

	return &PageAddIPFSGateway{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonAdd:   buttonAdd,
		txtName:     txtName,
		txtEndpoint: txtEndpoint,

		list: list,
	}
}

func (p *PageAddIPFSGateway) IsActive() bool {
	return p.isActive
}

func (p *PageAddIPFSGateway) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Add IPFS Gateway") }
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageAddIPFSGateway) Leave() {
	if page_instance.header.IsHistory(PAGE_ADD_IPFS_GATEWAY) {
		p.animationEnter.Reset()
		p.animationLeave.Start()
	} else {
		p.isActive = false
	}
}

func (p *PageAddIPFSGateway) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonAdd.Clicked() {
		p.submitForm(gtx)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("Here, you can add your own IPFS Gateway. The endpoint connection must be a HTTP connection, starting with http:// or https:// for TLS connection. Use {cid} to set where the content identifier must be pasted."))
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "deronfts.com")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th, lang.Translate("Endpoint"), "https://ipfs.deronfts.com/ipfs/{cid}")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonAdd.Text = lang.Translate("ADD GATEWAY")
			return p.buttonAdd.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddIPFSGateway) submitForm(gtx layout.Context) {
	p.buttonAdd.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonAdd.SetLoading(false)
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		txtName := p.txtName.Editor()
		txtEndpoint := p.txtEndpoint.Editor()

		if txtName.Text() == "" {
			setError(fmt.Errorf("enter name"))
			return
		}

		if txtEndpoint.Text() == "" {
			setError(fmt.Errorf("enter endpoint"))
			return
		}

		endpoint := txtEndpoint.Text()
		gateway := app_data.IPFSGateway{
			Name:     txtName.Text(),
			Endpoint: endpoint,
			Active:   true,
		}

		err := gateway.TestFetch()
		if err != nil {
			setError(err)
			return
		}

		err = app_data.InsertIPFSGateway(gateway)
		if err != nil {
			setError(err)
			return
		}

		p.buttonAdd.SetLoading(false)
		notification_modals.SuccessInstance.SetText(lang.Translate("Success"), "new noded added")
		notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		page_instance.header.GoBack()
	}()
}
