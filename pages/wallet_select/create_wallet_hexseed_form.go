package page_wallet_select

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageCreateWalletHexSeedForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List

	txtHexSeed         *components.TextField
	txtWalletName      *components.TextField
	txtPassword        *components.TextField
	txtConfirmPassword *components.TextField
	buttonCreate       *components.Button
}

var _ router.Page = &PageCreateWalletHexSeedForm{}

func NewPageCreateWalletHexSeedForm() *PageCreateWalletHexSeedForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	txtWalletName := components.NewTextField()
	txtPassword := components.NewPasswordTextField()
	txtConfirmPassword := components.NewPasswordTextField()

	txtHexSeed := components.NewTextField()

	iconCreate, _ := widget.NewIcon(icons.ContentAddBox)
	buttonCreate := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            iconCreate,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	return &PageCreateWalletHexSeedForm{
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		txtHexSeed:         txtHexSeed,
		txtWalletName:      txtWalletName,
		txtPassword:        txtPassword,
		txtConfirmPassword: txtConfirmPassword,
		buttonCreate:       buttonCreate,
	}
}

func (p *PageCreateWalletHexSeedForm) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Recover from Hex Seed"))

	if !page_instance.header.IsHistory(PAGE_CREATE_WALLET_HEXSEED_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageCreateWalletHexSeedForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageCreateWalletHexSeedForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletHexSeedForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	if p.buttonCreate.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText("Success", lang.Translate("New wallet created"))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtHexSeed.Layout(gtx, th, lang.Translate("Hex Seed"), lang.Translate("Enter hex seed of 64 chars"))
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
			return p.buttonCreate.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtHexSeed.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtWalletName.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	if p.txtPassword.Input.Clickable.Clicked() {
		p.list.ScrollTo(2)
	}

	if p.txtConfirmPassword.Input.Clickable.Clicked() {
		p.list.ScrollTo(3)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
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
