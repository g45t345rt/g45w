package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/transaction"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/image_modal"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/qrcode_scan_modal"
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

	txtAmount *prefabs.TextField

	buttonBuildTx    *components.Button
	buttonOptions    *components.Button
	buttonSetMax     *components.Button
	balanceContainer *BalanceContainer
	tokenContainer   *TokenContainer
	walletAddrInput  *WalletAddrInput

	token *wallet_manager.Token

	ringSizeSelector *prefabs.RingSizeSelector

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSendForm{}

func NewPageSendForm() *PageSendForm {
	buildIcon, _ := widget.NewIcon(icons.HardwareMemory)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonBuildTx := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        buildIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		LoadingIcon: loadingIcon,
		Animation:   components.NewButtonAnimationDefault(),
	})
	buttonBuildTx.Label.Alignment = text.Middle
	buttonBuildTx.Style.Font.Weight = font.Bold

	txtAmount := prefabs.NewNumberTextField()
	txtAmount.Input.TextSize = unit.Sp(26)
	txtAmount.Input.FontWeight = font.Bold

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

	buttonSetMax := components.NewButton(components.ButtonStyle{
		TextSize: unit.Sp(16),
	})
	buttonSetMax.Style.Font.Weight = font.Bold

	balanceContainer := NewBalanceContainer()
	tokenContainer := NewTokenContainer()
	walletAddrInput := NewWalletAddrInput()

	return &PageSendForm{
		txtAmount:        txtAmount,
		buttonBuildTx:    buttonBuildTx,
		ringSizeSelector: ringSizeSelector,
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		list:             list,
		buttonOptions:    buttonOptions,
		buttonSetMax:     buttonSetMax,
		balanceContainer: balanceContainer,
		tokenContainer:   tokenContainer,
		walletAddrInput:  walletAddrInput,
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
			p.buttonBuildTx.SetLoading(true)
			err := p.prepareTx()
			p.buttonBuildTx.SetLoading(false)
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
		if p.ringSizeSelector.Changed {
			settings.App.SendRingSize = p.ringSizeSelector.Size
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
					return p.walletAddrInput.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
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
		p.list.ScrollTo(2)
	}

	if p.walletAddrInput.txtWalletAddr.Clickable.Clicked() {
		p.list.ScrollTo(3)
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
	txtWalletAddr := p.walletAddrInput.txtWalletAddr
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
	wallet := wallet_manager.OpenedWallet

	txtAmount := p.txtAmount
	if txtAmount.Value() == "" {
		return fmt.Errorf(lang.Translate("Amount cannot be empty."))
	}

	amount := &utils.ShiftNumber{Decimals: int(p.token.Decimals)}
	err := amount.Parse(txtAmount.Value())
	if err != nil {
		return err
	}

	if amount.Number == 0 {
		return fmt.Errorf(lang.Translate("Amount must be greater than 0."))
	}

	txtWalletAddr := p.walletAddrInput.txtWalletAddr
	if txtWalletAddr.Value() == "" {
		return fmt.Errorf(lang.Translate("Destination address is empty."))
	}

	txtComment := page_instance.pageSendOptionsForm.txtComment
	txtDstPort := page_instance.pageSendOptionsForm.txtDstPort

	var arguments rpc.Arguments

	addrValue := txtWalletAddr.Value()
	address, err := rpc.NewAddress(addrValue)
	if err != nil {
		addrString, err := wallet.Memory.NameToAddress(addrValue)

		if err != nil {
			if utils.IsErrLeafNotFound(err) {
				return fmt.Errorf("address not found for [%s]", addrValue)
			}

			return err
		}

		address, err = rpc.NewAddress(addrString)
		if err != nil {
			return err
		}
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
	ringsize := uint64(p.ringSizeSelector.Size)
	deroBalance, _ := wallet.Memory.Get_Balance()

	transfers := []rpc.Transfer{
		{
			SCID:        scId,
			Destination: address.String(),
			Amount:      amount.Number,
			Payload_RPC: arguments,
		},
	}

	if scId.IsZero() && deroBalance == amount.Number {
		// sender is trying to send entire Dero balance to another wallet
		// let's calculate fees before and deduct

		transfers[0].Amount = 0 // set amount to 0 or transaction won't build because you don't have enough funds
		_, txFees, _, err := wallet.BuildTransaction(transfers, ringsize, nil, false)
		if err != nil {
			return err
		}

		transfers[0].Amount = amount.Number - txFees
	}

	build_tx_modal.Instance.Open(build_tx_modal.TxPayload{
		Transfers:  transfers,
		Ringsize:   ringsize,
		TokensInfo: []*wallet_manager.Token{p.token},
	})

	return nil
}

type TokenContainer struct {
	nameEditor       *widget.Editor
	scIdEditor       *widget.Editor
	tokenImagerHover *prefabs.ImageHoverClick
	token            *wallet_manager.Token
}

func NewTokenContainer() *TokenContainer {
	nameEditor := new(widget.Editor)
	nameEditor.ReadOnly = true
	nameEditor.SingleLine = true

	scIdEditor := new(widget.Editor)
	scIdEditor.ReadOnly = true
	scIdEditor.SingleLine = true

	tokenImagerHover := prefabs.NewImageHoverClick()

	return &TokenContainer{
		scIdEditor:       scIdEditor,
		nameEditor:       nameEditor,
		tokenImagerHover: tokenImagerHover,
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
				t.tokenImagerHover.Image.Src = t.token.LoadImageOp()

				if t.tokenImagerHover.Clickable.Clicked() {
					image_modal.Instance.Open(t.token.Name, t.tokenImagerHover.Image.Src)
				}

				gtx.Constraints.Max.X = gtx.Dp(50)
				gtx.Constraints.Max.Y = gtx.Dp(50)
				return t.tokenImagerHover.Layout(gtx)
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

type WalletAddrInput struct {
	txtWalletAddr       *prefabs.Input
	buttonAddrMenu      *components.Button
	newContactClickable *widget.Clickable

	txtDims layout.Dimensions
}

func NewWalletAddrInput() *WalletAddrInput {
	txtWalletAddr := prefabs.NewInput()

	addrIcon, _ := widget.NewIcon(icons.SocialPeople)
	buttonAddrMenu := components.NewButton(components.ButtonStyle{
		Rounded: components.UniformRounded(unit.Dp(5)),
		Icon:    addrIcon,
		/*Inset: layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(13), Right: unit.Dp(13),
		},*/
		Animation: components.NewButtonAnimationDefault(),
	})

	return &WalletAddrInput{
		txtWalletAddr:       txtWalletAddr,
		buttonAddrMenu:      buttonAddrMenu,
		newContactClickable: new(widget.Clickable),
	}
}

func (p *WalletAddrInput) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.buttonAddrMenu.Clicked() {
		go func() {
			contactIcon, _ := widget.NewIcon(icons.SocialGroup)
			scanIcon, _ := widget.NewIcon(app_icons.QRCodeScanner)

			keyChan := listselect_modal.Instance.Open([]*listselect_modal.SelectListItem{
				listselect_modal.NewSelectListItem("contact_list",
					listselect_modal.NewItemText(contactIcon, lang.Translate("Contact list")).Layout,
				),
				listselect_modal.NewSelectListItem("scan_qrcode",
					listselect_modal.NewItemText(scanIcon, lang.Translate("Scan QR Code")).Layout,
				),
			})

			for key := range keyChan {
				switch key {
				case "contact_list":
					page_instance.pageRouter.SetCurrent(PAGE_CONTACTS)
					page_instance.header.AddHistory(PAGE_CONTACTS)
				case "scan_qrcode":
					qrcode_scan_modal.Instance.Open()
				}
			}
		}()
	}

	{
		sent, value := qrcode_scan_modal.Instance.Value()
		if sent {
			p.txtWalletAddr.SetValue(value)
		}
	}

	var childs []layout.FlexChild

	childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(20), lang.Translate("Wallet Addr / Name"))
		lbl.Font.Weight = font.Bold
		return lbl.Layout(gtx)
	}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					p.txtWalletAddr.Colors = theme.Current.InputColors
					p.txtDims = p.txtWalletAddr.Layout(gtx, th, "")
					return p.txtDims
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					p.buttonAddrMenu.Style.Colors = theme.Current.ButtonPrimaryColors
					p.buttonAddrMenu.Flex = true
					size := image.Pt(p.txtDims.Size.Y, p.txtDims.Size.Y)
					gtx.Constraints.Min = size
					gtx.Constraints.Max = size
					return p.buttonAddrMenu.Layout(gtx, th)
				}),
			)
		}),
	)

	addr := p.txtWalletAddr.Editor.Text()

	wallet := wallet_manager.OpenedWallet
	if wallet != nil {
		contact, _ := wallet.GetContact(addr)
		if contact != nil {
			childs = append(childs,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
				}),
			)
		} else if addr != "" {
			if p.newContactClickable.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
			}

			if p.newContactClickable.Clicked() {
				page_instance.pageContactForm.ClearForm()
				page_instance.pageContactForm.txtAddr.SetValue(addr)
				page_instance.pageRouter.SetCurrent(PAGE_CONTACT_FORM)
				page_instance.header.AddHistory(PAGE_CONTACT_FORM)
			}

			childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return p.newContactClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("Create new contact?"))
							return lbl.Layout(gtx)
						})
					}),
				)
			}))
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
}
