package page_wallet

import (
	"fmt"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	crypto "github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
)

type PageDEXSCBridgeOut struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	txtAmount        *prefabs.TextField
	txtEthAddr       *prefabs.TextField
	balanceContainer *BalanceContainer
	buttonBridge     *components.Button
	token            *wallet_manager.Token
	infoRows         []*prefabs.InfoRow
	ringSizeSelector *prefabs.RingSizeSelector

	deroBridgeFee uint64
	bridgeOpened  bool

	list *widget.List
}

var _ router.Page = &PageDEXSCBridgeOut{}

func NewPageDEXSCBridgeOut() *PageDEXSCBridgeOut {

	list := new(widget.List)
	list.Axis = layout.Vertical

	ethereumIcon, _ := widget.NewIcon(app_icons.Ethereum)
	buttonBridge := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      ethereumIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonBridge.Label.Alignment = text.Middle
	buttonBridge.Style.Font.Weight = font.Bold

	balanceContainer := NewBalanceContainer()
	txtAmount := prefabs.NewNumberTextField()
	txtEthAddr := prefabs.NewTextField()
	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_DEX_SC_BRIDGE_OUT)

	return &PageDEXSCBridgeOut{
		headerPageAnimation: headerPageAnimation,
		list:                list,
		buttonBridge:        buttonBridge,
		balanceContainer:    balanceContainer,
		txtAmount:           txtAmount,
		txtEthAddr:          txtEthAddr,
		infoRows:            prefabs.NewInfoRows(2),
		ringSizeSelector:    prefabs.NewRingSizeSelector(16),
	}
}

func (p *PageDEXSCBridgeOut) IsActive() bool {
	return p.isActive
}

func (p *PageDEXSCBridgeOut) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return p.token.Name }
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), lang.Translate("Bridge to Eth"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.RightLayout = nil
}

func (p *PageDEXSCBridgeOut) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageDEXSCBridgeOut) SetToken(token *wallet_manager.Token) {
	p.token = token
	p.balanceContainer.SetToken(p.token)

	var result rpc.GetSC_Result
	err := wallet_manager.RPCCall("DERO.GetSC", rpc.GetSC_Params{
		SCID:       p.token.SCID,
		Code:       false,
		Variables:  false,
		KeysString: []string{"bridgeFee", "bridgeOpen"},
	}, &result)
	if err != nil {
		return
	}

	p.deroBridgeFee, _ = strconv.ParseUint(result.ValuesString[0], 10, 64)
	p.bridgeOpened, _ = strconv.ParseBool(result.ValuesString[1])
}

func (p *PageDEXSCBridgeOut) submitForm() error {
	amount := utils.ShiftNumber{Decimals: int(p.token.Decimals)}
	err := amount.Parse(p.txtAmount.Value())
	if err != nil {
		return err
	}

	ethAddr := p.txtEthAddr.Value()
	if !eth_common.IsHexAddress(ethAddr) {
		return fmt.Errorf("not a valid eth address")
	}

	if ethAddr == strings.ToLower(ethAddr) || ethAddr == strings.ToUpper(ethAddr) {
		return fmt.Errorf("ethereum address must be in CamelCase (mixed case) not all lower or all upper")
	}

	build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(randomAddr string) build_tx_modal.TxPayload {
		return build_tx_modal.TxPayload{
			Transfer: rpc.Transfer_Params{
				SC_RPC: rpc.Arguments{
					{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
					{Name: rpc.SCID, DataType: rpc.DataHash, Value: crypto.HashHexToHash(p.token.SCID)},
					{Name: "entrypoint", DataType: rpc.DataString, Value: "Bridge"},
					{Name: "eth_addr", DataType: rpc.DataString, Value: ethAddr},
				},
				Transfers: []rpc.Transfer{
					rpc.Transfer{SCID: p.token.GetHash(), Burn: amount.Number, Destination: randomAddr},
					rpc.Transfer{SCID: crypto.ZEROHASH, Burn: p.deroBridgeFee, Destination: randomAddr},
				},
				Ringsize: uint64(p.ringSizeSelector.Size),
			},
			TokensInfo: []*wallet_manager.Token{p.token},
		}
	})

	return nil
}

func (p *PageDEXSCBridgeOut) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonBridge.Clicked(gtx) {
		go func() {
			err := p.submitForm()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			}
		}()
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.balanceContainer.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.txtAmount.Layout(gtx, th, lang.Translate("Amount"), "")
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.txtEthAddr.Layout(gtx, th, lang.Translate("Ethereum Address"), "")
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.ringSizeSelector.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return p.infoRows[0].Layout(gtx, th, lang.Translate("Bridge Opened"), fmt.Sprint(p.bridgeOpened))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				amount := utils.ShiftNumber{Number: p.deroBridgeFee, Decimals: 5}
				return p.infoRows[1].Layout(gtx, th, lang.Translate("Bridge Out Fee"), fmt.Sprintf("%s DERO", amount.Format()))
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonBridge.Style.Colors = theme.Current.ButtonPrimaryColors
		p.buttonBridge.Text = lang.Translate("BRIDGE OUT")
		return p.buttonBridge.Layout(gtx, th)
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
