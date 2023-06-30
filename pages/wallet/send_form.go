package page_wallet

import (
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
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/app_instance"
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

	SCID             string
	txtAmount        *components.TextField
	txtWalletAddr    *components.Input
	txtComment       *components.TextField
	txtDescription   *components.TextField
	txtDstPort       *components.TextField
	buttonBuildTx    *components.Button
	accordionOptions *components.Accordion
	buttonContacts   *components.Button

	ringSizeSelector *prefabs.RingSizeSelector

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list *widget.List
}

var _ router.Page = &PageSendForm{}

func NewPageSendForm() *PageSendForm {
	th := app_instance.Theme
	buildIcon, _ := widget.NewIcon(icons.HardwareMemory)
	buttonBuildTx := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            buildIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonBuildTx.Label.Alignment = text.Middle
	buttonBuildTx.Style.Font.Weight = font.Bold

	txtAmount := components.NewTextField(th, lang.Translate("Amount"), "")
	txtWalletAddr := components.NewInput(th, "")
	txtComment := components.NewTextField(th, lang.Translate("Comment"), lang.Translate("The comment is natively encrypted."))
	txtComment.Editor().SingleLine = false
	txtComment.Editor().Submit = false
	txtDescription := components.NewTextField(th, lang.Translate("Description"), lang.Translate("Saved locally in your wallet."))
	txtDescription.Editor().SingleLine = false
	txtDescription.Editor().Submit = false
	txtDstPort := components.NewTextField(th, lang.Translate("Destination Port"), "")

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	ringSizeSelector := prefabs.NewRingSizeSelector("16")

	buttonOptions := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
	})
	buttonOptions.Label.Alignment = text.Middle
	buttonOptions.Style.Font.Weight = font.Bold
	accordionOptions := components.NewAccordion(components.AccordionStyle{
		Border: widget.Border{
			CornerRadius: unit.Dp(5),
			Color:        color.NRGBA{A: 255},
			Width:        unit.Dp(2),
		},
		Inset:  layout.UniformInset(unit.Dp(10)),
		Button: buttonOptions,
	}, false)

	contactIcon, _ := widget.NewIcon(icons.SocialPerson)
	buttonContacts := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            contactIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Inset: layout.Inset{
			Top: unit.Dp(14), Bottom: unit.Dp(13),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
		Animation: components.NewButtonAnimationDefault(),
	})

	return &PageSendForm{
		txtAmount:        txtAmount,
		txtWalletAddr:    txtWalletAddr,
		txtComment:       txtComment,
		txtDstPort:       txtDstPort,
		txtDescription:   txtDescription,
		buttonBuildTx:    buttonBuildTx,
		ringSizeSelector: ringSizeSelector,
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		list:             list,
		accordionOptions: accordionOptions,
		buttonContacts:   buttonContacts,
	}
}

func (p *PageSendForm) IsActive() bool {
	return p.isActive
}

func (p *PageSendForm) Enter() {
	p.isActive = true
	p.animationEnter.Start()
	p.animationLeave.Reset()
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

	if p.buttonBuildTx.Clickable.Clicked() {

	}

	if p.buttonContacts.Clickable.Clicked() {
		page_instance.header.AddHistory(PAGE_CONTACTS)
		page_instance.pageRouter.SetCurrent(PAGE_CONTACTS)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {

			r := op.Record(gtx.Ops)
			dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(16), "Send DERO")
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(14), "00000...00000")
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
			return p.txtAmount.Layout(gtx, th)
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
							return p.txtWalletAddr.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return p.buttonContacts.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					addr := p.txtWalletAddr.Editor().Text()
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
			p.accordionOptions.Style.Button.Text = lang.Translate("Options")
			return p.accordionOptions.Layout(gtx, th, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						p.txtComment.Input.EditorMinY = gtx.Dp(75)
						return p.txtComment.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return p.txtDstPort.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						p.txtDescription.Input.EditorMinY = gtx.Dp(75)
						return p.txtDescription.Layout(gtx, th)
					}),
				)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonBuildTx.Text = lang.Translate("BUILD TRANSACTION")
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

func (p *PageSendForm) submitForm() error {
	wallet := wallet_manager.OpenedWallet.Memory

	destination := p.txtWalletAddr.Value()

	var arguments rpc.Arguments

	comment := p.txtComment.Value()
	if len(comment) > 0 {
		arguments = append(arguments, rpc.Argument{Name: rpc.RPC_COMMENT, DataType: rpc.DataString, Value: comment})
	}

	destPortString := p.txtDstPort.Value()
	if len(destPortString) > 0 {
		destPort, err := strconv.ParseUint(destPortString, 10, 64)
		if err != nil {
			return err
		}

		arguments = append(arguments, rpc.Argument{Name: rpc.RPC_DESTINATION_PORT, DataType: rpc.DataUint64, Value: destPort})
	}

	scId := crypto.HashHexToHash(p.SCID)
	transfers := []rpc.Transfer{
		{SCID: scId, Destination: destination, Amount: 0, Payload_RPC: arguments},
	}

	ringsize, err := strconv.ParseUint(p.ringSizeSelector.Value, 10, 64)
	if err != nil {
		return err
	}

	tx, err := wallet.TransferPayload0(transfers, ringsize, false, rpc.Arguments{}, 0, false)
	if err != nil {
		return err
	}

	err = wallet.SendTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}
