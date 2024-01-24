package page_wallet

import (
	"encoding/csv"
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
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/containers/password_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/pages"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSettings struct {
	isActive bool

	buttonDeleteWallet      *components.Button
	buttonInfo              *components.Button
	buttonServiceNames      *components.Button
	txtWalletName           *prefabs.TextField
	txtWalletChangePassword *prefabs.TextField
	buttonSave              *components.Button
	buttonCleanWallet       *components.Button
	buttonExportTxs         *components.Button
	buttonAddDEXTokens      *components.Button

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	action string

	list *widget.List
}

var _ router.Page = &PageSettings{}

func NewPageSettings() *PageSettings {
	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDeleteWallet := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      deleteIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonDeleteWallet.Label.Alignment = text.Middle
	buttonDeleteWallet.Style.Font.Weight = font.Bold

	cleanIcon, _ := widget.NewIcon(icons.ContentDeleteSweep)
	buttonCleanWallet := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      cleanIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonCleanWallet.Label.Alignment = text.Middle
	buttonCleanWallet.Style.Font.Weight = font.Bold

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	buttonSave := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      saveIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonSave.Label.Alignment = text.Middle
	buttonSave.Style.Font.Weight = font.Bold

	infoIcon, _ := widget.NewIcon(icons.ActionInfo)
	buttonInfo := components.NewButton(components.ButtonStyle{
		Icon:      infoIcon,
		TextSize:  unit.Sp(16),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonInfo.Label.Alignment = text.Middle
	buttonInfo.Style.Font.Weight = font.Bold

	nameIcon, _ := widget.NewIcon(app_icons.Badge)
	buttonServiceNames := components.NewButton(components.ButtonStyle{
		Icon:      nameIcon,
		TextSize:  unit.Sp(16),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonServiceNames.Label.Alignment = text.Middle
	buttonServiceNames.Style.Font.Weight = font.Bold

	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	exportIcon, _ := widget.NewIcon(icons.EditorPublish)
	buttonExportTxs := components.NewButton(components.ButtonStyle{
		Icon:      exportIcon,
		TextSize:  unit.Sp(16),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
		LoadingIcon: loadingIcon,
	})
	buttonExportTxs.Label.Alignment = text.Middle
	buttonExportTxs.Style.Font.Weight = font.Bold

	swapIcon, _ := widget.NewIcon(app_icons.Swap)
	buttonAddDEXTokens := components.NewButton(components.ButtonStyle{
		Icon:      swapIcon,
		TextSize:  unit.Sp(16),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
		LoadingIcon: loadingIcon,
	})
	buttonAddDEXTokens.Label.Alignment = text.Middle
	buttonAddDEXTokens.Style.Font.Weight = font.Bold

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	txtWalletName := prefabs.NewTextField()
	txtWalletChangePassword := prefabs.NewPasswordTextField()

	return &PageSettings{
		buttonDeleteWallet:      buttonDeleteWallet,
		animationEnter:          animationEnter,
		animationLeave:          animationLeave,
		list:                    list,
		txtWalletName:           txtWalletName,
		txtWalletChangePassword: txtWalletChangePassword,
		buttonSave:              buttonSave,
		buttonInfo:              buttonInfo,
		buttonCleanWallet:       buttonCleanWallet,
		buttonExportTxs:         buttonExportTxs,
		buttonServiceNames:      buttonServiceNames,
		buttonAddDEXTokens:      buttonAddDEXTokens,
	}
}

func (p *PageSettings) IsActive() bool {
	return p.isActive
}

func (p *PageSettings) Enter() {
	openedWallet := wallet_manager.OpenedWallet
	walletName := openedWallet.Info.Name
	page_instance.header.Title = func() string { return lang.Translate("Settings") }
	p.txtWalletName.SetValue(walletName)
	page_instance.header.Subtitle = nil
	page_instance.header.RightLayout = nil
	page_instance.header.LeftLayout = nil

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
	if p.buttonDeleteWallet.Clicked(gtx) {
		p.action = "delete_wallet"
		password_modal.Instance.SetVisible(true)
	}

	if p.buttonSave.Clicked(gtx) {
		p.action = "save_changes"
		password_modal.Instance.SetVisible(true)
	}

	if p.buttonServiceNames.Clicked(gtx) {
		page_instance.pageRouter.SetCurrent(PAGE_SERVICE_NAMES)
		page_instance.header.AddHistory(PAGE_SERVICE_NAMES)
	}

	if p.buttonInfo.Clicked(gtx) {
		p.action = "wallet_info"
		password_modal.Instance.SetVisible(true)
	}

	if p.buttonCleanWallet.Clicked(gtx) {
		p.action = "clean_wallet"
		password_modal.Instance.SetVisible(true)
	}

	if p.buttonExportTxs.Clicked(gtx) {
		p.action = "export_txs"
		password_modal.Instance.SetVisible(true)
	}

	if p.buttonAddDEXTokens.Clicked(gtx) {
		go func() {
			yes := <-confirm_modal.Instance.Open(confirm_modal.ConfirmText{})
			if yes {
				wallet := wallet_manager.OpenedWallet
				err := wallet.InsertDexTokensFolder()
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				} else {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.SUCCESS,
						Title: lang.Translate("Success"),
						Text:  lang.Translate("The folder was created."),
					})
				}
			}
		}()
	}

	submitted, password := password_modal.Instance.Input.Submitted()
	if submitted {
		wallet := wallet_manager.OpenedWallet
		validPassword := wallet.Memory.Check_Password(password)

		if !validPassword {
			password_modal.Instance.StartWrongPassAnimation()
		} else {
			password_modal.Instance.Modal.SetVisible(false)
			err := p.submitForm(gtx, password)

			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
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
			p.buttonServiceNames.Text = lang.Translate("Service Names")

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.buttonServiceNames.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonServiceNames.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Manage wallet service names"))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonInfo.Text = lang.Translate("Wallet Information")

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.buttonInfo.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonInfo.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("e.g seed phrase, hex seed, etc..."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, unit.Dp(5))
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletName.Layout(gtx, th, lang.Translate("Wallet Name"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtWalletChangePassword.Layout(gtx, th, lang.Translate("Change Password"), "Enter new password")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonSave.Text = lang.Translate("SAVE CHANGES")
			p.buttonSave.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonSave.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, unit.Dp(5))
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.buttonAddDEXTokens.Text = lang.Translate("Add DEX tokens")
					p.buttonAddDEXTokens.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonAddDEXTokens.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Create a folder and insert DEX tokens automatically."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.buttonExportTxs.Text = lang.Translate("Export Transactions")
					p.buttonExportTxs.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonExportTxs.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Export all your transactions into a CSV file."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.buttonCleanWallet.Text = lang.Translate("CLEAN WALLET")
					p.buttonCleanWallet.Style.Colors = theme.Current.ButtonPrimaryColors
					return p.buttonCleanWallet.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Delete cache data and rescan entire wallet."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, unit.Dp(5))
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonDeleteWallet.Text = lang.Translate("DELETE WALLET")
			p.buttonDeleteWallet.Style.Colors = theme.Current.ButtonDangerColors
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
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageSettings) submitForm(gtx layout.Context, password string) error {
	wallet := wallet_manager.OpenedWallet

	switch p.action {
	case "wallet_info":
		page_instance.pageRouter.SetCurrent(PAGE_WALLET_INFO)
		page_instance.header.AddHistory(PAGE_WALLET_INFO)
	case "clean_wallet":
		wallet.Memory.Clean()

		notification_modal.Open(notification_modal.Params{
			Type:       notification_modal.SUCCESS,
			Title:      lang.Translate("Success"),
			Text:       lang.Translate("Walled cleaned."),
			CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
		})
	case "delete_wallet":
		err := wallet.Delete()
		if err != nil {
			return err
		}

		page_instance.header.GoBack()
		app_instance.Router.SetCurrent(pages.PAGE_WALLET_SELECT)
		wallet_manager.CloseOpenedWallet()

		notification_modal.Open(notification_modal.Params{
			Type:       notification_modal.SUCCESS,
			Title:      lang.Translate("Success"),
			Text:       lang.Translate("Wallet deleted."),
			CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
		})
	case "export_txs":
		go func() error {
			setError := func(err error) error {
				p.buttonExportTxs.SetLoading(false)
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
				return err
			}

			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.INFO,
				Title:      lang.Translate("Info"),
				Text:       lang.Translate("Exporting transactions..."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})

			account := wallet.Memory.GetAccount()
			p.buttonExportTxs.SetLoading(true)

			file, err := app_instance.Explorer.CreateFile("transactions.csv")
			if err != nil {
				return setError(err)
			}

			writer := csv.NewWriter(file)
			defer writer.Flush()

			header := []string{"SCID", "TXID", "Height", "Blockhash",
				"Coinbase", "Incoming", "Destination", "Atomic Amount",
				"Atomic Burn", "Atomic Fees", "Proof", "Time", "EWData",
				"Sender", "Destination Port", "Source Port"}
			err = writer.Write(header)
			if err != nil {
				return setError(err)
			}

			for scId, entries := range account.EntriesNative {
				for _, entry := range entries {
					sender := wallet.GetTxSender(wallet_manager.Entry{Entry: entry})
					destination := wallet.GetTxDestination(wallet_manager.Entry{Entry: entry})

					row := []string{scId.String(), entry.TXID, fmt.Sprint(entry.Height), entry.BlockHash,
						fmt.Sprint(entry.Coinbase), fmt.Sprint(entry.Incoming), destination, fmt.Sprint(entry.Amount),
						fmt.Sprint(entry.Burn), fmt.Sprint(entry.Fees), entry.Proof, entry.Time.String(), entry.EWData,
						sender, fmt.Sprint(entry.DestinationPort), fmt.Sprint(entry.SourcePort)}
					err = writer.Write(row)
					if err != nil {
						return setError(err)
					}
				}
			}

			p.buttonExportTxs.SetLoading(false)
			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.SUCCESS,
				Title:      lang.Translate("Success"),
				Text:       lang.Translate("Transactions exported."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
			return nil
		}()
	case "save_changes":
		newWalletName := p.txtWalletName.Value()
		if wallet.Info.Name != newWalletName {
			err := wallet.Rename(newWalletName)
			if err != nil {
				return err
			}
		}

		newPassword := p.txtWalletChangePassword.Value()
		if newPassword != "" {
			err := wallet.ChangePassword(password, p.txtWalletChangePassword.Value())
			if err != nil {
				return err
			}

			p.txtWalletChangePassword.SetValue("")

			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.SUCCESS,
				Title:      lang.Translate("Success"),
				Text:       lang.Translate("Data saved."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
		}
	}

	return nil
}
