package page_wallet

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/transaction"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/registration"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/shirou/gopsutil/cpu"
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

	// updated every second if normalreg running
	statusText       string
	cpuUsageText     string
	progressBarValue float64
	probabilityText  string
}

func NewRegisterWalletForm() *RegisterWalletForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	txtThreadCount := components.NewTextField()

	logicalCores, err := cpu.Counts(true)
	if err != nil {
		txtThreadCount.SetValue("1")
	} else {
		// recommend for normal reg is different than the fast reg
		recommendedWorkers := math.Floor(float64(logicalCores) * 6)
		txtThreadCount.SetValue(fmt.Sprint(recommendedWorkers))
	}

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

	page := &RegisterWalletForm{
		list:           list,
		txtThreadCount: txtThreadCount,
		buttonStart:    buttonStart,
		buttonStop:     buttonStop,

		normalReg: normalReg,
	}

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			if normalReg.Running {
				percent, err := cpu.Percent(0, false)
				if len(percent) == 1 && err == nil {
					page.cpuUsageText = fmt.Sprintf("CPU Usage: %.2f%%", percent[0])
				}

				hashRateText := utils.FormatHashRate(normalReg.HashRate())
				page.statusText = fmt.Sprintf("%d | %s", normalReg.HashCount(), hashRateText)

				// https://bitcoin.stackexchange.com/questions/114580/finding-hash-with-11-leading-zeroes
				target := float64(3)
				probability := 1 - math.Pow(1-math.Pow(16, -(target*2)), float64(normalReg.HashCount()))
				page.probabilityText = fmt.Sprintf("Probability: %.2f%%", probability*100)
				page.progressBarValue = probability

				w.Invalidate()
			} else {
				page.cpuUsageText = ""
				page.probabilityText = ""
				page.progressBarValue = 0
				hashRateText := utils.FormatHashRate(0)
				page.statusText = fmt.Sprintf("%d | %s", 0, hashRateText)
			}
		}
	}()

	return page
}

func (p *RegisterWalletForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if p.buttonStart.Clicked() {
		err := p.startRegistration()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	if p.buttonStop.Clicked() {
		p.normalReg.Stop()
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), lang.Translate("The Dero blockchain is an account base model and requires a one time POW registration proccess to avoid spam."))
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.txtThreadCount.Layout(gtx, th, lang.Translate("Worker Count"), "")
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(13), lang.Translate("By default, the worker count is set to the recommended value for your device. More workers is faster but takes more cpu ressources."))
				return lbl.Layout(gtx)
			}),
		)
	})

	if p.cpuUsageText != "" {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), p.cpuUsageText)
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), p.probabilityText)
					return lbl.Layout(gtx)
				}),
			)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(16), lang.Translate("Progress"))
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return components.ProgressBar{
					Value:   float32(p.progressBarValue),
					Color:   color.NRGBA{A: 255},
					BgColor: color.NRGBA{R: 255, G: 255, B: 255, A: 200},
					Rounded: unit.Dp(5),
					Height:  unit.Dp(20),
				}.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {

				label := material.Label(th, unit.Sp(16), p.statusText)
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		if p.normalReg.Running {
			p.buttonStop.Text = lang.Translate("STOP")
			return p.buttonStop.Layout(gtx, th)
		}

		p.buttonStart.Text = lang.Translate("START")
		return p.buttonStart.Layout(gtx, th)
	})

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

	wallet := wallet_manager.OpenedWallet
	p.normalReg.Start(int(threadCount), wallet.Memory)
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
	if p.buttonSend.Clicked() {
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

	err = wallet.InsertOutgoingTx(tx, "")
	if err != nil {
		return err
	}

	err = wallet.Memory.SendTransaction(tx)
	if err != nil {
		return err
	}

	return nil
}
