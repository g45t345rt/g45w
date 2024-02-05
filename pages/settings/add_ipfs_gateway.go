package page_settings

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAddIPFSGateway struct {
	isActive            bool
	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonAdd   *components.Button
	txtEndpoint *prefabs.TextField
	txtName     *prefabs.TextField

	list *widget.List
}

var _ router.Page = &PageAddIPFSGateway{}

func NewPageAddIPFSGateway() *PageAddIPFSGateway {
	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        addIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonAdd.Label.Alignment = text.Middle
	buttonAdd.Style.Font.Weight = font.Bold

	txtName := prefabs.NewTextField()
	txtEndpoint := prefabs.NewTextField()

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_ADD_IPFS_GATEWAY)
	return &PageAddIPFSGateway{
		headerPageAnimation: headerPageAnimation,

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
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Add IPFS Gateway") }
	page_instance.header.Subtitle = nil
	page_instance.header.RightLayout = nil
}

func (p *PageAddIPFSGateway) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageAddIPFSGateway) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonAdd.Clicked(gtx) {
		p.submitForm(gtx)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("Here, you can add your own IPFS Gateway. The endpoint connection must be a HTTP connection, starting with http:// or https:// for TLS connection. Use {cid} to set where the content identifier must be pasted."))
			lbl.Color = theme.Current.TextMuteColor
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
			p.buttonAdd.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonAdd.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddIPFSGateway) submitForm(gtx layout.Context) {
	p.buttonAdd.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonAdd.SetLoading(false)
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
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
		gateway := app_db.IPFSGateway{
			Name:     txtName.Text(),
			Endpoint: endpoint,
			Active:   true,
		}

		err := gateway.TestFetch()
		if err != nil {
			setError(err)
			return
		}

		err = app_db.InsertIPFSGateway(gateway)
		if err != nil {
			setError(err)
			return
		}

		p.buttonAdd.SetLoading(false)
		notification_modal.Open(notification_modal.Params{
			Type:       notification_modal.SUCCESS,
			Title:      lang.Translate("Success"),
			Text:       lang.Translate("New node added."),
			CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
		})
		page_instance.header.GoBack()
	}()
}
