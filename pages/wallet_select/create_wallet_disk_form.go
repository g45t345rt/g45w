package page_wallet_select

import (
	"fmt"
	"image/color"
	"os"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
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

type PageCreateWalletDiskForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List

	buttonCreate  *components.Button
	txtWalletName *components.TextField
	txtPassword   *components.TextField
	buttonLoad    *components.Button

	walletPath string
}

var _ router.Page = &PageCreateWalletDiskForm{}

func NewPageCreateWalletDiskForm() *PageCreateWalletDiskForm {
	th := app_instance.Theme
	list := new(widget.List)
	list.Axis = layout.Vertical

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	txtWalletName := components.NewTextField(th, lang.Translate("Wallet Name"), "")
	txtPassword := components.NewPasswordTextField(th, lang.Translate("Password"), "")

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

	iconOpen, _ := widget.NewIcon(icons.FileFolderOpen)
	buttonLoad := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            iconOpen,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	return &PageCreateWalletDiskForm{
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		txtWalletName: txtWalletName,
		txtPassword:   txtPassword,
		buttonCreate:  buttonCreate,
		buttonLoad:    buttonLoad,
	}
}

func (p *PageCreateWalletDiskForm) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_CREATE_WALLET_DISK_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	page_instance.header.SetTitle(lang.Translate("Load from Disk"))
}

func (p *PageCreateWalletDiskForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageCreateWalletDiskForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletDiskForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonLoad.Clickable.Clicked() {
		go func() {
			read, err := app_instance.Explorer.ChooseFile()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
				return
			}

			switch f := read.(type) {
			case *os.File:
				p.walletPath = f.Name()
			default:
			}
		}()
	}

	if p.buttonCreate.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Wallet loaded successfully"))
			notification_modals.SuccessInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			p.buttonLoad.Text = lang.Translate("LOAD FILE")
			return p.buttonLoad.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			path := lang.Translate("Not file selected.")
			if p.walletPath != "" {
				path = p.walletPath
			}

			lbl := material.Label(th, unit.Sp(16), path)
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtPassword.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletName.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonCreate.Text = lang.Translate("CREATE WALLET")
			return p.buttonCreate.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtPassword.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	if p.txtWalletName.Input.Clickable.Clicked() {
		p.list.ScrollTo(2)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageCreateWalletDiskForm) submitForm() error {
	txtName := p.txtWalletName.Editor()
	txtPassword := p.txtPassword.Editor()

	if txtName.Text() == "" {
		return fmt.Errorf("enter wallet name")
	}

	err := wallet_manager.CreateWalletFromPath(txtName.Text(), txtPassword.Text(), p.walletPath)
	if err != nil {
		return err
	}

	return nil
}
