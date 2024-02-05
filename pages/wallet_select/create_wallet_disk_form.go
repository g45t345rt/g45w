package page_wallet_select

import (
	"fmt"
	"os"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageCreateWalletDiskForm struct {
	isActive            bool
	headerPageAnimation *prefabs.PageHeaderAnimation

	list *widget.List

	buttonCreate  *components.Button
	txtWalletName *prefabs.TextField
	txtPassword   *prefabs.TextField
	buttonLoad    *components.Button

	walletPath string
	walletData []byte
}

var _ router.Page = &PageCreateWalletDiskForm{}

func NewPageCreateWalletDiskForm() *PageCreateWalletDiskForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	txtWalletName := prefabs.NewTextField()
	txtPassword := prefabs.NewPasswordTextField()

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

	iconOpen, _ := widget.NewIcon(icons.FileFolderOpen)
	buttonLoad := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      iconOpen,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonLoad.Style.Font.Weight = font.Bold

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_CREATE_WALLET_DISK_FORM)
	return &PageCreateWalletDiskForm{
		list:                list,
		headerPageAnimation: headerPageAnimation,

		txtWalletName: txtWalletName,
		txtPassword:   txtPassword,
		buttonCreate:  buttonCreate,
		buttonLoad:    buttonLoad,
	}
}

func (p *PageCreateWalletDiskForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Load from Disk") }
}

func (p *PageCreateWalletDiskForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageCreateWalletDiskForm) IsActive() bool {
	return p.isActive
}

func (p *PageCreateWalletDiskForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonLoad.Clicked(gtx) {
		go func() {
			loadFile := func() (filePath string, data []byte, err error) {
				file, err := app_instance.Explorer.ChooseFile()
				if err != nil {
					return
				}

				switch f := file.(type) {
				case *os.File:
					filePath = f.Name()
				default:
				}

				reader := utils.ReadCloser{ReadCloser: file}
				data, err = reader.ReadAll()
				return
			}

			filePath, data, err := loadFile()
			if err != nil {
				p.walletPath = ""
				p.walletData = []byte{}

				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
				return
			} else {
				p.walletPath = filePath
				p.walletData = data
			}
		}()
	}

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
				Text:       lang.Translate("Wallet loaded successfully."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			p.buttonLoad.Text = lang.Translate("LOAD FILE")
			p.buttonLoad.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonLoad.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			path := lang.Translate("No file selected.")
			if len(p.walletData) > 0 {
				path = lang.Translate("Data was loaded.")
			}

			if p.walletPath != "" {
				path = p.walletPath
			}

			lbl := material.Label(th, unit.Sp(16), path)
			lbl.Color = theme.Current.TextMuteColor
			lbl.WrapPolicy = text.WrapGraphemes
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtPassword.Layout(gtx, th, lang.Translate("Password"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletName.Layout(gtx, th, lang.Translate("Wallet Name"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonCreate.Text = lang.Translate("CREATE WALLET")
			p.buttonCreate.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonCreate.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtPassword.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(2)
	}

	if p.txtWalletName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(3)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageCreateWalletDiskForm) submitForm() error {
	txtName := p.txtWalletName.Editor()
	txtPassword := p.txtPassword.Editor()

	if txtName.Text() == "" {
		return fmt.Errorf("enter wallet name")
	}

	err := wallet_manager.CreateWalletFromData(txtName.Text(), txtPassword.Text(), p.walletData)
	if err != nil {
		return err
	}

	page_instance.header.GoBack()
	return nil
}
