package page_wallet

import (
	"context"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/token_manager"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAddSCForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonCheckSC *components.Button
	txtSCID       *components.TextField
	tokenManager  *token_manager.TokenManager

	listStyle material.ListStyle
}

var _ router.Page = &PageAddSCForm{}

func NewPageAddSCForm() *PageAddSCForm {
	th := app_instance.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	checkIcon, _ := widget.NewIcon(icons.ActionSearch)
	buttonCheckSC := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "CHECK SC",
		Icon:            checkIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonCheckSC.Label.Alignment = text.Middle
	buttonCheckSC.Style.Font.Weight = font.Bold

	txtSCID := components.NewTextField(th, "SCID", "Smart Contract ID")

	return &PageAddSCForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonCheckSC: buttonCheckSC,
		txtSCID:       txtSCID,

		listStyle: listStyle,
	}
}

func (p *PageAddSCForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddSCForm) Enter() {
	p.isActive = true
	page_instance.header.SetTitle("Add Smart Contract")
	page_instance.header.Subtitle = nil
	p.animationEnter.Start()
	p.animationLeave.Reset()
	addr := wallet_manager.Instance.OpenedWallet.Info.Addr
	p.tokenManager = token_manager.New(addr)
}

func (p *PageAddSCForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageAddSCForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonCheckSC.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtSCID.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonCheckSC.Layout(gtx, th)
		},
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddSCForm) submitForm() error {
	_, _, err := p.fetchSmartContract(p.txtSCID.Value())
	if err != nil {
		return err
	}

	return nil
}

func (p *PageAddSCForm) fetchSmartContract(scId string) (scType sc.SCType, result *rpc.GetSC_Result, err error) {
	err = walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
		SCID:      scId,
		Variables: true,
		Code:      true,
	}, &result)
	if err != nil {
		return sc.UNKNOWN_TYPE, nil, err
	}

	scType = sc.CheckType(result.Code)

	return scType, result, nil
}
