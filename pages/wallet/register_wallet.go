package page_wallet

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"math"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/transaction"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/registration"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageRegisterWallet struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	registerWalletForm   *RegisterWalletForm
	sendRegistrationForm *SendRegistrationForm
}

var _ router.Page = &PageRegisterWallet{}

func NewPageRegisterWallet() *PageRegisterWallet {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	registerWalletForm := NewRegisterWalletForm()
	sendRegistrationForm := NewSendRegistrationForm()

	return &PageRegisterWallet{
		animationEnter:       animationEnter,
		animationLeave:       animationLeave,
		registerWalletForm:   registerWalletForm,
		sendRegistrationForm: sendRegistrationForm,
	}
}

func (p *PageRegisterWallet) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_REGISTER_WALLET) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageRegisterWallet) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageRegisterWallet) IsActive() bool {
	return p.isActive
}

func (p *PageRegisterWallet) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	regTxHex := wallet_manager.OpenedWallet.Info.RegistrationTxHex
	if regTxHex != "" {
		return p.sendRegistrationForm.Layout(gtx, th)
	}

	return p.registerWalletForm.Layout(gtx, th)
}

type RegisterWalletForm struct {
	list *widget.List

	txtThreadCount *components.TextField
	buttonStart    *components.Button
	buttonStop     *components.Button

	normalReg *registration.NormalReg
}

func NewRegisterWalletForm() *RegisterWalletForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	txtThreadCount := components.NewTextField()
	txtThreadCount.SetValue("1")

	buildIcon, _ := widget.NewIcon(icons.HardwareMemory)
	buttonStart := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            buildIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	stopIcon, _ := widget.NewIcon(icons.AVPause)
	buttonStop := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            stopIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	w := app_instance.Window
	normalReg := registration.NewNormalReg()
	normalReg.OnFound = func(tx *transaction.Transaction) {
		wallet := wallet_manager.OpenedWallet
		err := wallet_manager.StoreRegistrationTx(wallet.Info.Addr, tx)
		if err != nil {

		}
		w.Invalidate()
	}

	return &RegisterWalletForm{
		list:           list,
		txtThreadCount: txtThreadCount,
		buttonStart:    buttonStart,
		buttonStop:     buttonStop,

		normalReg: normalReg,
	}
}

func (p *RegisterWalletForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.buttonStart.Clickable.Clicked() {
		err := p.startRegistration()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	if p.buttonStop.Clickable.Clicked() {
		p.normalReg.Stop()
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(14), lang.Translate("The Dero blockchain is an account base model and requires a one time POW registration proccess to avoid spam."))
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.txtThreadCount.Layout(gtx, th, lang.Translate("Worker Count"), "")
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(13), lang.Translate("More workers is faster but takes more cpu ressources. You decide!"))
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Label(th, unit.Sp(16), lang.Translate("Progress"))
					label.Font.Weight = font.Bold
					return label.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					target := float64(3)

					// https://bitcoin.stackexchange.com/questions/114580/finding-hash-with-11-leading-zeroes
					value := 1 - math.Pow(1-math.Pow(16, -(target*2)), float64(p.normalReg.HashCount()))
					return components.ProgressBar{
						Value:   float32(value),
						Color:   color.NRGBA{A: 255},
						BgColor: color.NRGBA{R: 255, G: 255, B: 255, A: 200},
						Rounded: unit.Dp(5),
						Height:  unit.Dp(20),
					}.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					hashRate := utils.FormatHashRate(p.normalReg.HashRate())
					status := fmt.Sprintf("%d | %s", p.normalReg.HashCount(), hashRate)
					label := material.Label(th, unit.Sp(16), status)
					label.Font.Weight = font.Bold
					return label.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			if p.normalReg.Running {
				p.buttonStop.Text = lang.Translate("STOP")
				return p.buttonStop.Layout(gtx, th)
			}

			p.buttonStart.Text = lang.Translate("START")
			return p.buttonStart.Layout(gtx, th)
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

func (p *RegisterWalletForm) startRegistration() error {
	threadCount, err := strconv.ParseUint(p.txtThreadCount.Value(), 10, 64)
	if err != nil {
		return err
	}

	wallet := wallet_manager.OpenedWallet.Memory
	p.normalReg.Start(int(threadCount), wallet)
	return nil
}

type SendRegistrationForm struct {
	list       *widget.List
	buttonSend *components.Button
}

func NewSendRegistrationForm() *SendRegistrationForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	sendIcon, _ := widget.NewIcon(icons.ContentSend)
	buttonSend := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            sendIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSend.Style.Font.Weight = font.Bold

	return &SendRegistrationForm{
		list:       list,
		buttonSend: buttonSend,
	}
}

func (p *SendRegistrationForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.buttonSend.Clickable.Clicked() {
		err := p.sendTransaction()
		if err != nil {
			notification_modals.ErrorInstance.SetVisible(true, 0)
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
		} else {
			recent_txs_modal.Instance.SetVisible(true)
			page_instance.header.GoBack()
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("The registration POW has been completed succesfully. You can now send the solution to the network to finalize the registration process."))
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonSend.Text = lang.Translate("SEND TRANSACTION")
			return p.buttonSend.Layout(gtx, th)
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

func (p *SendRegistrationForm) sendTransaction() error {
	wallet := wallet_manager.OpenedWallet
	txHex := wallet.Info.RegistrationTxHex
	data, err := hex.DecodeString(txHex)
	if err != nil {
		return err
	}

	tx := new(transaction.Transaction)
	err = tx.Deserialize(data)
	if err != nil {
		return err
	}

	err = wallet.StoreOutgoingTx(tx, "")
	if err != nil {
		return err
	}

	err = wallet.Memory.SendTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}
