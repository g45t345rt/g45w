package page_wallet_select

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageCreateWalletHexSeedForm struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	list *widget.List

	txtHexSeed         *prefabs.TextField
	txtWalletName      *prefabs.TextField
	txtPassword        *prefabs.TextField
	txtConfirmPassword *prefabs.TextField
	buttonCreate       *components.Button
}

var _ router.Page = &PageCreateWalletHexSeedForm{}

func NewPageCreateWalletHexSeedForm() *PageCreateWalletHexSeedForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	txtWalletName := prefabs.NewTextField()
	txtPassword := prefabs.NewPasswordTextField()
	txtConfirmPassword := prefabs.NewPasswordTextField()

	txtHexSeed := prefabs.NewTextField()

	iconCreate, _ := widget.NewIcon(icons.ContentAddBox)
	buttonCreate := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      iconCreate,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonCreate.Style.Font.Weight = font.Bold
	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_CREATE_WALLET_HEXSEED_FORM)

	return &PageCreateWalletHexSeedForm{
		list:                list,
		headerPageAnimation: headerPageAnimation,

		txtHexSeed:         txtHexSeed,
		txtWalletName:      txtWalletName,
		txtPassword:        txtPassword,
		txtConfirmPassword: txtConfirmPassword,
		buttonCreate:       buttonCreate,
	}
}

func (p *PageCreateWalletHexSeedForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Recover from Hex Seed") }
}

func (p *PageCreateWalletHexSeedForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageCreateWalletHexSeedForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletHexSeedForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonCreate.Clicked(gtx) {
		err := p.submitForm()
		if err != nil {
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
		} else {
			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.SUCCESS,
				Title:      lang.Translate("Success"),
				Text:       lang.Translate("New wallet created."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtHexSeed.Layout(gtx, th, lang.Translate("Hex Seed"), lang.Translate("Enter hex seed of 64 chars."))
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletName.Layout(gtx, th, lang.Translate("Wallet Name"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtPassword.Layout(gtx, th, lang.Translate("Password"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtConfirmPassword.Layout(gtx, th, lang.Translate("Confirm Password"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonCreate.Text = lang.Translate("RECOVER WALLET")
			p.buttonCreate.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonCreate.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtHexSeed.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtWalletName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(1)
	}

	if p.txtPassword.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(2)
	}

	if p.txtConfirmPassword.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(3)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageCreateWalletHexSeedForm) submitForm() error {
	txtName := p.txtWalletName.Editor()
	txtPassword := p.txtPassword.Editor()
	txtConfirmPassword := p.txtConfirmPassword.Editor()
	txtHexSeed := p.txtHexSeed.Editor()

	if txtHexSeed.Text() == "" {
		return fmt.Errorf("enter hex seed")
	}

	if txtName.Text() == "" {
		return fmt.Errorf("enter wallet name")
	}

	if txtPassword.Text() == "" {
		return fmt.Errorf("enter password")
	}

	if txtPassword.Text() != txtConfirmPassword.Text() {
		return fmt.Errorf("the confirm password does not match")
	}

	err := wallet_manager.CreateWalletFromHexSeed(txtName.Text(), txtPassword.Text(), txtHexSeed.Text())
	if err != nil {
		return err
	}

	txtName.SetText("")
	txtPassword.SetText("")
	txtConfirmPassword.SetText("")
	txtHexSeed.SetText("")

	page_instance.header.GoBack()
	return nil
}
