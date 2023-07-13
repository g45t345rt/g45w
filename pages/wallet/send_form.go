package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
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

type PageSendForm struct {
	isActive bool

	SCID           string
	txtAmount      *components.TextField
	txtWalletAddr  *components.Input
	buttonSendTx   *components.Button
	buttonContacts *components.Button
	buttonOptions  *components.Button

	ringSizeSelector *prefabs.RingSizeSelector

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSendForm{}

func NewPageSendForm() *PageSendForm {
	sendIcon, _ := widget.NewIcon(icons.ContentSend)
	buttonSendTx := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            sendIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSendTx.Label.Alignment = text.Middle
	buttonSendTx.Style.Font.Weight = font.Bold

	txtAmount := components.NewTextField()
	txtWalletAddr := components.NewInput()

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	ringSizeSelector := prefabs.NewRingSizeSelector("16")

	optionIcon, _ := widget.NewIcon(icons.ActionSettingsEthernet)
	buttonOptions := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextSize:        unit.Sp(14),
		Icon:            optionIcon,
		IconGap:         unit.Dp(10),
		TextColor:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		BackgroundColor: color.NRGBA{A: 0},
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonOptions.Label.Alignment = text.Middle
	buttonOptions.Style.Font.Weight = font.Bold

	contactIcon, _ := widget.NewIcon(icons.SocialPerson)
	buttonContacts := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            contactIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Inset: layout.Inset{
			Top: unit.Dp(13), Bottom: unit.Dp(12),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
		Animation: components.NewButtonAnimationDefault(),
	})

	return &PageSendForm{
		txtAmount:        txtAmount,
		txtWalletAddr:    txtWalletAddr,
		buttonSendTx:     buttonSendTx,
		ringSizeSelector: ringSizeSelector,
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		list:             list,
		buttonContacts:   buttonContacts,
		buttonOptions:    buttonOptions,
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

	if p.buttonSendTx.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	if p.buttonOptions.Clickable.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_SEND_OPTIONS_FORM)
		page_instance.header.AddHistory(PAGE_SEND_OPTIONS_FORM)
	}

	if p.buttonContacts.Clickable.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_CONTACTS)
		page_instance.header.AddHistory(PAGE_CONTACTS)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {

			r := op.Record(gtx.Ops)
			dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(18), "Selected Asset [DERO]")
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), "00000...00000")
						lbl.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
						return lbl.Layout(gtx)
					}),
				)
			})
			c := r.Stop()

			paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.UniformRRect(
				image.Rectangle{Max: dims.Size},
				gtx.Dp(10),
			).Op(gtx.Ops))

			c.Add(gtx.Ops)
			return dims
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtAmount.Layout(gtx, th, lang.Translate("Amount"), "")
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
							return p.txtWalletAddr.Layout(gtx, th, "")
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return p.buttonContacts.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					addr := p.txtWalletAddr.Editor.Text()
					contact, ok := page_instance.contactManager.Contacts[addr]

					if ok {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										lbl := material.Label(th, unit.Sp(16), lang.Translate("Matching contact:"))
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

					return layout.Dimensions{}
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.ringSizeSelector.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonOptions.Text = lang.Translate("OPTIONS")
			return p.buttonOptions.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonSendTx.Text = lang.Translate("SEND TRANSACTION")
			return p.buttonSendTx.Layout(gtx, th)
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

func (p *PageSendForm) submitForm() error {
	wallet := wallet_manager.OpenedWallet

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
	txtDescription := page_instance.pageSendOptionsForm.txtDescription

	var arguments rpc.Arguments

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

	scId := crypto.HashHexToHash(p.SCID)
	transfers := []rpc.Transfer{
		{SCID: scId, Destination: txtWalletAddr.Value(), Amount: amount, Payload_RPC: arguments},
	}

	ringsize, err := strconv.ParseUint(p.ringSizeSelector.Value, 10, 64)
	if err != nil {
		return err
	}

	tx, err := wallet.Memory.TransferPayload0(transfers, ringsize, false, rpc.Arguments{}, 0, false)
	if err != nil {
		return err
	}

	err = wallet.StoreOutgoingTx(tx, txtDescription.Value())
	if err != nil {
		return err
	}

	err = wallet.Memory.SendTransaction(tx)
	if err != nil {
		return err
	}

	txtWalletAddr.SetValue("")
	txtAmount.SetValue("")
	txtDescription.SetValue("")
	txtComment.SetValue("")
	txtDstPort.SetValue("")

	recent_txs_modal.Instance.SetVisible(true)
	page_instance.header.GoBack()

	return nil
}
