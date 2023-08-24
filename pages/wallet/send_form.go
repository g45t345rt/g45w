package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/transaction"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSendForm struct {
	isActive bool

	txtAmount        *prefabs.TextField
	txtWalletAddr    *components.Input
	buttonBuildTx    *components.Button
	buttonAddr       *components.Button
	buttonOptions    *components.Button
	buttonSetMax     *components.Button
	balanceContainer *BalanceContainer
	tokenContainer   *TokenContainer
	qrScanCamModal   *prefabs.CameraQRScanModal
	addrMenuSelect   *AddrMenuSelect

	token *wallet_manager.Token

	ringSizeSelector *prefabs.RingSizeSelector

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSendForm{}

func NewPageSendForm() *PageSendForm {
	buildIcon, _ := widget.NewIcon(icons.HardwareMemory)
	buttonBuildTx := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      buildIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonBuildTx.Label.Alignment = text.Middle
	buttonBuildTx.Style.Font.Weight = font.Bold

	txtAmount := prefabs.NewNumberTextField()
	txtAmount.Input.TextSize = unit.Sp(26)
	txtAmount.Input.FontWeight = font.Bold

	txtWalletAddr := components.NewInput()

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	defaultRingSize := settings.App.SendRingSize
	ringSizeSelector := prefabs.NewRingSizeSelector(defaultRingSize)

	optionIcon, _ := widget.NewIcon(icons.ActionSettingsEthernet)
	buttonOptions := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Icon:      optionIcon,
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonOptions.Label.Alignment = text.Middle
	buttonOptions.Style.Font.Weight = font.Bold

	addrIcon, _ := widget.NewIcon(icons.SocialPeople)
	buttonAddr := components.NewButton(components.ButtonStyle{
		Rounded: components.UniformRounded(unit.Dp(5)),
		Icon:    addrIcon,
		Inset: layout.Inset{
			Top: unit.Dp(14), Bottom: unit.Dp(14),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
	})

	buttonSetMax := components.NewButton(components.ButtonStyle{
		TextSize: unit.Sp(16),
	})
	buttonSetMax.Style.Font.Weight = font.Bold

	balanceContainer := NewBalanceContainer()
	tokenContainer := NewTokenContainer()
	addrMenuSelect := NewAddrMenuSelect()

	qrScanCamModal := prefabs.NewCameraQRScanModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			qrScanCamModal.Layout(gtx, th)
		},
	})

	return &PageSendForm{
		txtAmount:        txtAmount,
		txtWalletAddr:    txtWalletAddr,
		buttonBuildTx:    buttonBuildTx,
		ringSizeSelector: ringSizeSelector,
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		list:             list,
		buttonAddr:       buttonAddr,
		buttonOptions:    buttonOptions,
		buttonSetMax:     buttonSetMax,
		balanceContainer: balanceContainer,
		tokenContainer:   tokenContainer,
		addrMenuSelect:   addrMenuSelect,
		qrScanCamModal:   qrScanCamModal,
	}
}

func (p *PageSendForm) IsActive() bool {
	return p.isActive
}

func (p *PageSendForm) Enter() {
	p.isActive = true
	if !page_instance.header.IsHistory(PAGE_SEND_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
	page_instance.pageBalanceTokens.ResetWalletHeader()
}

func (p *PageSendForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSendForm) SetToken(token *wallet_manager.Token) {
	p.token = token
	p.balanceContainer.SetToken(p.token)
	p.tokenContainer.SetToken(p.token)
}

func (p *PageSendForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonBuildTx.Clicked() {
		go func() {
			err := p.prepareTx()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			}
		}()
	}

	if p.buttonOptions.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_SEND_OPTIONS_FORM)
		page_instance.header.AddHistory(PAGE_SEND_OPTIONS_FORM)
	}

	if p.buttonAddr.Clicked() {
		p.addrMenuSelect.SelectModal.Modal.SetVisible(true)

	}

	{
		selected, key := p.addrMenuSelect.SelectModal.Selected()
		if selected {
			switch key {
			case "contact_list":
				page_instance.pageRouter.SetCurrent(PAGE_CONTACTS)
				page_instance.header.AddHistory(PAGE_CONTACTS)
			case "scan_qrcode":
				p.qrScanCamModal.Show()
			}
			p.addrMenuSelect.SelectModal.Modal.SetVisible(false)
		}
	}

	{
		sent, value := p.qrScanCamModal.Value()
		if sent {
			p.txtWalletAddr.SetValue(value)
		}
	}

	if p.buttonSetMax.Clicked() {
		wallet := wallet_manager.OpenedWallet
		balance, _ := wallet.Memory.Get_Balance_scid(p.token.GetHash())
		amount := utils.ShiftNumber{Number: balance, Decimals: int(p.token.Decimals)}.Format()
		p.txtAmount.SetValue(amount)
	}

	if build_tx_modal.Instance.TxSent() {
		p.ClearForm()
	}

	{
		changed, value := p.ringSizeSelector.Changed()
		if changed {
			settings.App.SendRingSize = value
			settings.Save()
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.tokenContainer.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.balanceContainer.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					v := utils.ShiftNumber{Number: 0, Decimals: int(p.token.Decimals)}
					return p.txtAmount.Layout(gtx, th, lang.Translate("Amount"), v.Format())
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, layout.Spacer{}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							p.buttonSetMax.Text = lang.Translate("SET MAX")
							p.buttonSetMax.Style.Colors = theme.Current.ModalButtonColors
							return p.buttonSetMax.Layout(gtx, th)
						}),
					)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate("Wallet Addr / Name"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							p.txtWalletAddr.Colors = theme.Current.InputColors
							return p.txtWalletAddr.Layout(gtx, th, "")
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							p.buttonAddr.Style.Colors = theme.Current.ButtonPrimaryColors
							return p.buttonAddr.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					addr := p.txtWalletAddr.Editor.Text()

					wallet := wallet_manager.OpenedWallet
					if wallet != nil {
						contact, _ := wallet.GetContact(addr)
						if contact != nil {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											lbl := material.Label(th, unit.Sp(16), lang.Translate("Matching contact:"))
											lbl.Color = theme.Current.TextMuteColor
											return lbl.Layout(gtx)
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(3)}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											lbl := material.Label(th, unit.Sp(16), contact.Name)
											lbl.Font.Weight = font.Bold
											return lbl.Layout(gtx)
										}),
									)
								}),
							)
						}
					}

					return layout.Dimensions{}
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.ringSizeSelector.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonOptions.Text = lang.Translate("OPTIONS")
			p.buttonOptions.Style.Colors = theme.Current.ButtonSecondaryColors
			return p.buttonOptions.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonBuildTx.Text = lang.Translate("BUILD TRANSACTION")
			p.buttonBuildTx.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonBuildTx.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtAmount.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	if p.txtWalletAddr.Clickable.Clicked() {
		p.list.ScrollTo(2)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageSendForm) ClearForm() {
	txtAmount := p.txtAmount
	txtWalletAddr := p.txtWalletAddr
	txtComment := page_instance.pageSendOptionsForm.txtComment
	txtDstPort := page_instance.pageSendOptionsForm.txtDstPort
	txtDescription := page_instance.pageSendOptionsForm.txtDescription

	txtWalletAddr.SetValue("")
	txtAmount.SetValue("")
	txtDescription.SetValue("")
	txtComment.SetValue("")
	txtDstPort.SetValue("")
	p.list.ScrollTo(0)
}

func (p *PageSendForm) prepareTx() error {
	txtAmount := p.txtAmount
	if txtAmount.Value() == "" {
		return fmt.Errorf(lang.Translate("Amount cannot be empty."))
	}

	amount, err := globals.ParseAmount(txtAmount.Value())
	if err != nil {
		return err
	}

	if amount == 0 {
		return fmt.Errorf(lang.Translate("Amount must be greater than 0."))
	}

	txtWalletAddr := p.txtWalletAddr
	if txtWalletAddr.Value() == "" {
		return fmt.Errorf(lang.Translate("Destination address is empty."))
	}

	txtComment := page_instance.pageSendOptionsForm.txtComment
	txtDstPort := page_instance.pageSendOptionsForm.txtDstPort

	var arguments rpc.Arguments

	address, err := rpc.NewAddress(txtWalletAddr.Value())
	if err != nil {
		return err
	}

	if address.IsIntegratedAddress() {
		err = address.Arguments.Validate_Arguments()
		if err != nil {
			return err
		}

		if !address.Arguments.Has(rpc.RPC_DESTINATION_PORT, rpc.DataUint64) {
			return fmt.Errorf(lang.Translate("The integrated address does not contain a destination port."))
		}

		destinationPort := address.Arguments.Value(rpc.RPC_DESTINATION_PORT, rpc.DataUint64).(uint64)
		arguments = append(arguments, rpc.Argument{Name: rpc.RPC_DESTINATION_PORT, DataType: rpc.DataUint64, Value: destinationPort})

		if address.Arguments.Has(rpc.RPC_COMMENT, rpc.DataString) {
			comment := address.Arguments.Value(rpc.RPC_COMMENT, rpc.DataString).(string)
			arguments = append(arguments, rpc.Argument{Name: rpc.RPC_COMMENT, DataType: rpc.DataString, Value: comment})
		}

		if address.Arguments.Has(rpc.RPC_EXPIRY, rpc.DataTime) {
			expireTime := address.Arguments.Value(rpc.RPC_EXPIRY, rpc.DataTime).(time.Time)
			if expireTime.Before(time.Now().UTC()) {
				return fmt.Errorf(lang.Translate("The integrated address has expired."))
			}
		}
	} else {
		comment := txtComment.Value()
		if len(comment) > 0 {
			arguments = append(arguments, rpc.Argument{Name: rpc.RPC_COMMENT, DataType: rpc.DataString, Value: comment})
		}

		destPortString := txtDstPort.Value()
		if len(destPortString) > 0 {
			destPort, err := strconv.ParseUint(destPortString, 10, 64)
			if err != nil {
				return err
			}

			arguments = append(arguments, rpc.Argument{Name: rpc.RPC_DESTINATION_PORT, DataType: rpc.DataUint64, Value: destPort})
		}
	}

	_, err = arguments.CheckPack(transaction.PAYLOAD0_LIMIT)
	if err != nil {
		return err
	}

	scId := p.token.GetHash()
	ringsize := uint64(p.ringSizeSelector.Value)

	wallet := wallet_manager.OpenedWallet
	balance, _ := wallet.Memory.Get_Balance_scid(scId)

	transfers := []rpc.Transfer{
		{
			SCID:        scId,
			Destination: address.String(),
			Amount:      amount,
			Payload_RPC: arguments,
		},
	}

	if scId.IsZero() && balance == amount {
		// sender is trying to send all Dero to another wallet
		// let's calculate fees before and deduct
		fees, err := wallet.CalculateFees(ringsize, transfers, arguments)
		if err != nil {
			return err
		}

		transfers[0].Amount = amount - fees
	}

	build_tx_modal.Instance.Open(build_tx_modal.TxPayload{
		Transfers: transfers,
		Ringsize:  ringsize,
		SCData:    rpc.Arguments{},
	})

	return nil
}

type TokenContainer struct {
	nameEditor *widget.Editor
	scIdEditor *widget.Editor
	tokenImage *components.Image
	token      *wallet_manager.Token
}

func NewTokenContainer() *TokenContainer {
	nameEditor := new(widget.Editor)
	nameEditor.ReadOnly = true
	nameEditor.SingleLine = true

	scIdEditor := new(widget.Editor)
	scIdEditor.ReadOnly = true
	scIdEditor.SingleLine = true

	tokenImage := &components.Image{
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}

	return &TokenContainer{
		scIdEditor: scIdEditor,
		nameEditor: nameEditor,
		tokenImage: tokenImage,
	}
}

func (t *TokenContainer) SetToken(token *wallet_manager.Token) {
	t.nameEditor.SetText(token.Name)
	scId := utils.ReduceTxId(token.SCID)

	if token.Symbol.Valid {
		scId = fmt.Sprintf("%s (%s)", scId, token.Symbol.String)
	}

	t.scIdEditor.SetText(scId)
	t.token = token
}

func (t *TokenContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				t.tokenImage.Src = t.token.LoadImageOp()
				gtx.Constraints.Max.X = gtx.Dp(50)
				gtx.Constraints.Max.Y = gtx.Dp(50)
				return t.tokenImage.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, t.nameEditor, "")
						editor.Font.Weight = font.Bold
						editor.TextSize = unit.Sp(22)
						return editor.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(th, t.scIdEditor, "")
						editor.Color = theme.Current.TextMuteColor
						return editor.Layout(gtx)
					}),
				)
			}),
		)
	})

	c := r.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor, clip.UniformRRect(
		image.Rectangle{Max: dims.Size},
		gtx.Dp(10),
	).Op(gtx.Ops))

	c.Add(gtx.Ops)
	return dims
}

type AddrMenuSelect struct {
	SelectModal *prefabs.SelectModal
}

func NewAddrMenuSelect() *AddrMenuSelect {
	var items []*prefabs.SelectListItem

	contactIcon, _ := widget.NewIcon(icons.SocialGroup)
	items = append(items, prefabs.NewSelectListItem("contact_list", prefabs.ListItemMenuItem{
		Icon:  contactIcon,
		Title: "Contact list", //@lang.Translate("Contact list")
	}.Layout))

	scanIcon, _ := widget.NewIcon(icons.HardwareScanner)
	items = append(items, prefabs.NewSelectListItem("scan_qrcode", prefabs.ListItemMenuItem{
		Icon:  scanIcon,
		Title: "Scan QR Code", //@lang.Translate("Scan QR Code")
	}.Layout))

	selectModal := prefabs.NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			selectModal.Layout(gtx, th, items)
		},
	})

	return &AddrMenuSelect{
		SelectModal: selectModal,
	}
}
