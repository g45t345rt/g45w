package page_wallet

import (
	"context"
	"fmt"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_icons"

	// "github.com/g45t345rt/g45w/bridge_metamask"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageDEXSCBridgeIn struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	balanceContainer *BalanceContainer
	txtAmount        *prefabs.TextField
	buttonConnect    *components.Button

	token    *wallet_manager.Token
	infoRows []*prefabs.InfoRow

	bridgeOpened bool

	list *widget.List
}

var _ router.Page = &PageDEXSCBridgeIn{}

func NewPageDEXSCBridgeIn() *PageDEXSCBridgeIn {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	ethereumIcon, _ := widget.NewIcon(app_icons.Ethereum)
	buttonConnect := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      ethereumIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonConnect.Label.Alignment = text.Middle
	buttonConnect.Style.Font.Weight = font.Bold

	txtAmount := prefabs.NewNumberTextField()

	return &PageDEXSCBridgeIn{
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		list:             list,
		infoRows:         prefabs.NewInfoRows(1),
		buttonConnect:    buttonConnect,
		txtAmount:        txtAmount,
		balanceContainer: NewBalanceContainer(),
	}
}

func (p *PageDEXSCBridgeIn) IsActive() bool {
	return p.isActive
}

func (p *PageDEXSCBridgeIn) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_DEX_SC_BRIDGE_IN) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	page_instance.header.Title = func() string { return p.token.Name }
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), lang.Translate("Bridge to Stargate"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.ButtonRight = nil
}

func (p *PageDEXSCBridgeIn) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageDEXSCBridgeIn) SetToken(token *wallet_manager.Token) {
	p.token = token
	p.balanceContainer.SetToken(p.token)

	var result rpc.GetSC_Result
	err := walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
		SCID:       p.token.SCID,
		Code:       false,
		Variables:  false,
		KeysString: []string{"bridgeOpen"},
	}, &result)
	if err != nil {
		return
	}

	p.bridgeOpened, _ = strconv.ParseBool(result.ValuesString[0])
}

// func (p *PageDEXSCBridgeIn) submitForm() error {
// 	//browser.OpenUrl("wc:0254c043869952027016fd697a276e137039af5a6f077284c000a74bda702857@2?relay-protocol=irn&symKey=22f98d50e70540eb9407f40c8a081bf9a18e221305cc5bf7306cecc12b5def66")

// 	amount, err := strconv.ParseFloat(p.txtAmount.Value(), 64)
// 	if err != nil {
// 		return err
// 	}

// 	wallet := wallet_manager.OpenedWallet
// 	addr := wallet.Memory.GetAddress()

// 	symbol := strings.Replace(p.token.Symbol.String, "D", "", 1)

// 	url, err := bridge_metamask.Link(bridge_metamask.BridgeInData{
// 		WalletAddress: addr.String(),
// 		Symbol:        symbol,
// 		Amount:        amount,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return browser.OpenUrl(url)
// }

func (p *PageDEXSCBridgeIn) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	// if p.buttonConnect.Clicked() {
	// 	go p.submitForm()
	// }

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.balanceContainer.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.txtAmount.Layout(gtx, th, lang.Translate("Amount"), "")
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.infoRows[0].Layout(gtx, th, lang.Translate("Bridge Opened"), fmt.Sprint(p.bridgeOpened))
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonConnect.Text = lang.Translate("BRIDGE IN")
		p.buttonConnect.Style.Colors = theme.Current.ButtonPrimaryColors
		return p.buttonConnect.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Bottom: unit.Dp(20),
			Left:   unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return widgets[index](gtx)
		})
	})
}
