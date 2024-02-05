package page_wallet_select

import (
	"fmt"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
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

type PageCreateWalletForm struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	list *widget.List

	txtWalletName      *prefabs.TextField
	txtPassword        *prefabs.TextField
	txtConfirmPassword *prefabs.TextField
	buttonCreate       *components.Button

	regResultContainer *RegResultContainer
}

var _ router.Page = &PageCreateWalletForm{}

func NewPageCreateWalletForm() *PageCreateWalletForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	txtWalletName := prefabs.NewTextField()
	txtPassword := prefabs.NewPasswordTextField()
	txtConfirmPassword := prefabs.NewPasswordTextField()

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

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_CREATE_WALLET_FORM)
	return &PageCreateWalletForm{
		list:                list,
		headerPageAnimation: headerPageAnimation,

		txtWalletName:      txtWalletName,
		txtPassword:        txtPassword,
		txtConfirmPassword: txtConfirmPassword,
		buttonCreate:       buttonCreate,
	}
}

func (p *PageCreateWalletForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)
	page_instance.header.Title = func() string { return lang.Translate("Create New Wallet") }
}

func (p *PageCreateWalletForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)

}

func (p *PageCreateWalletForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	var widgets []layout.Widget

	if p.regResultContainer != nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.regResultContainer.Layout(gtx, th)
		})
	}

	widgets = append(widgets,
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
			p.buttonCreate.Text = lang.Translate("CREATE WALLET")
			p.buttonCreate.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonCreate.Layout(gtx, th)
		},
	)

	list := material.List(th, p.list)
	list.AnchorStrategy = material.Overlay

	if p.txtWalletName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtPassword.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(1)
	}

	if p.txtConfirmPassword.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(2)
	}

	return list.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageCreateWalletForm) submitForm() error {
	txtName := p.txtWalletName.Editor()
	txtPassword := p.txtPassword.Editor()
	txtConfirmPassword := p.txtConfirmPassword.Editor()

	if txtName.Text() == "" {
		return fmt.Errorf("enter wallet name")
	}

	if txtPassword.Text() == "" {
		return fmt.Errorf("enter password")
	}

	if txtPassword.Text() != txtConfirmPassword.Text() {
		return fmt.Errorf("the confirm password does not match")
	}

	if p.regResultContainer != nil {
		hexSeed := p.regResultContainer.result.HexSeed
		err := wallet_manager.CreateWalletFromHexSeed(txtName.Text(), txtPassword.Text(), hexSeed)
		if err != nil {
			return err
		}

		tx := p.regResultContainer.result.Tx
		addr := p.regResultContainer.result.Addr
		err = wallet_manager.StoreRegistrationTx(addr, tx)
		if err != nil {
			return err
		}
	} else {
		err := wallet_manager.CreateRandomWallet(txtName.Text(), txtPassword.Text())
		if err != nil {
			return err
		}
	}

	p.regResultContainer = nil
	txtName.SetText("")
	txtPassword.SetText("")
	txtConfirmPassword.SetText("")

	page_instance.header.ResetHistory()
	page_instance.pageRouter.SetCurrent(PAGE_SELECT_WALLET)
	page_instance.header.AddHistory(PAGE_SELECT_WALLET)
	return nil
}

type RegResultContainer struct {
	result         *RegResult
	addrEditor     *widget.Editor
	wordSeedEditor *widget.Editor
	hexSeedEditor  *widget.Editor
}

func NewRegResultContainer(result *RegResult) *RegResultContainer {
	addrEditor := new(widget.Editor)
	addrEditor.SetText(result.Addr)
	addrEditor.WrapPolicy = text.WrapGraphemes
	addrEditor.ReadOnly = true

	wordSeedEditor := new(widget.Editor)
	wordSeedEditor.SetText(result.WordSeed)
	wordSeedEditor.WrapPolicy = text.WrapWords
	wordSeedEditor.ReadOnly = true

	hexSeedEditor := new(widget.Editor)
	hexSeedEditor.SetText(result.HexSeed)
	hexSeedEditor.WrapPolicy = text.WrapGraphemes
	hexSeedEditor.ReadOnly = true

	return &RegResultContainer{
		result:         result,
		addrEditor:     addrEditor,
		wordSeedEditor: wordSeedEditor,
		hexSeedEditor:  hexSeedEditor,
	}
}

func (item *RegResultContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("The registration process found the POW solution. You can now create your wallet and send the registration transaction."))
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			r := op.Record(gtx.Ops)
			dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Address"))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, item.addrEditor, "")
						editor.TextSize = unit.Sp(14)
						return editor.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Word Seed"))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, item.wordSeedEditor, "")
						editor.TextSize = unit.Sp(14)
						return editor.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), lang.Translate("Hex Seed"))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, item.hexSeedEditor, "")
						editor.TextSize = unit.Sp(14)
						return editor.Layout(gtx)
					}),
				)
			})
			c := r.Stop()

			paint.FillShape(
				gtx.Ops,
				theme.Current.ListBgColor,
				clip.UniformRRect(
					image.Rectangle{Max: dims.Size},
					gtx.Dp(10),
				).Op(gtx.Ops),
			)

			c.Add(gtx.Ops)
			return dims
		}),
	)
}
