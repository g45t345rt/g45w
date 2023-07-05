package page_wallet

import (
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
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSettings struct {
	isActive bool

	buttonDeleteWallet  *components.Button
	txtWalletName       *components.TextField
	buttonSave          *components.Button
	modalWalletPassword *prefabs.PasswordModal

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSettings{}

func NewPageSettings() *PageSettings {
	th := app_instance.Theme
	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDeleteWallet := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            deleteIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonDeleteWallet.Label.Alignment = text.Middle
	buttonDeleteWallet.Style.Font.Weight = font.Bold

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	buttonSave := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            saveIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSave.Label.Alignment = text.Middle
	buttonSave.Style.Font.Weight = font.Bold

	modalWalletPassword := prefabs.NewPasswordModal()

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			modalWalletPassword.Layout(gtx, th)
		},
	})

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	txtWalletName := components.NewTextField(th, lang.Translate("Name"), "")

	return &PageSettings{
		buttonDeleteWallet:  buttonDeleteWallet,
		animationEnter:      animationEnter,
		animationLeave:      animationLeave,
		list:                list,
		modalWalletPassword: modalWalletPassword,
		txtWalletName:       txtWalletName,
		buttonSave:          buttonSave,
	}
}

func (p *PageSettings) IsActive() bool {
	return p.isActive
}

func (p *PageSettings) Enter() {
	openedWallet := wallet_manager.OpenedWallet
	walletName := openedWallet.Info.Name
	page_instance.header.SetTitle(lang.Translate("Settings"))
	p.txtWalletName.SetValue(walletName)
	page_instance.header.Subtitle = nil

	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_SETTINGS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageSettings) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSettings) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.buttonDeleteWallet.Clickable.Clicked() {
		p.modalWalletPassword.Modal.SetVisible(gtx, true)
	}

	if p.buttonSave.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Changes applied successfully"))
			notification_modals.SuccessInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	submitted, text := p.modalWalletPassword.Submit()
	if submitted {
		openedWallet := wallet_manager.OpenedWallet
		err := wallet_manager.DeleteWallet(openedWallet.Info.Addr, text)
		if err == nil {
			p.modalWalletPassword.Modal.SetVisible(gtx, false)
			page_instance.pageRouter.SetCurrent(PAGE_BALANCE_TOKENS)
			app_instance.Router.SetCurrent(app_instance.PAGE_WALLET_SELECT)
			wallet_manager.OpenedWallet = nil
		} else {
			if err.Error() == "Invalid Password" {
				p.modalWalletPassword.StartWrongPassAnimation()
			} else {
				notification_modals.ErrorInstance.SetText("Error", err.Error())
				notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
			}
		}
	}

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

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletName.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonSave.Text = lang.Translate("SAVE CHANGES")
			return p.buttonSave.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = gtx.Dp(5)
			paint.FillShape(gtx.Ops, color.NRGBA{A: 150}, clip.Rect{
				Max: gtx.Constraints.Max,
			}.Op())

			return layout.Dimensions{Size: gtx.Constraints.Max}
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonDeleteWallet.Text = lang.Translate("DELETE WALLET")
			return p.buttonDeleteWallet.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageSettings) submitForm() error {
	walletInfo := wallet_manager.OpenedWallet.Info
	newWalletName := p.txtWalletName.Value()
	if walletInfo.Name != newWalletName {
		err := wallet_manager.RenameWallet(walletInfo.Addr, newWalletName)
		if err != nil {
			return err
		}

		wallet_manager.OpenedWallet.Info.Name = newWalletName
	}

	return nil
}
