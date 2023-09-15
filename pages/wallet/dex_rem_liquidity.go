package page_wallet

import (
	"fmt"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageDEXRemLiquidity struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	txtPercent         *prefabs.TextField
	buttonRemove       *components.Button
	infoRows           []*prefabs.InfoRow
	liquidityContainer *LiquidityContainer

	pair   dex_sc.Pair
	token1 *wallet_manager.Token
	token2 *wallet_manager.Token
	share  uint64

	list *widget.List
}

var _ router.Page = &PageDEXRemLiquidity{}

func NewPageDEXRemLiquidity() *PageDEXRemLiquidity {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	txtPercent := prefabs.NewNumberTextField()
	txtPercent.Input.TextSize = unit.Sp(18)
	txtPercent.Input.FontWeight = font.Bold

	removeIcon, _ := widget.NewIcon(icons.ContentClear)
	buttonRemove := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      removeIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonRemove.Label.Alignment = text.Middle
	buttonRemove.Style.Font.Weight = font.Bold

	return &PageDEXRemLiquidity{
		animationEnter:     animationEnter,
		animationLeave:     animationLeave,
		list:               list,
		txtPercent:         txtPercent,
		buttonRemove:       buttonRemove,
		infoRows:           prefabs.NewInfoRows(2),
		liquidityContainer: NewLiquidityContainer(),
	}
}

func (p *PageDEXRemLiquidity) IsActive() bool {
	return p.isActive
}

func (p *PageDEXRemLiquidity) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_DEX_REM_LIQUIDITY) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
	page_instance.header.Title = func() string {
		return lang.Translate("Remove Liquidity")
	}

	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), p.pair.Symbol)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.ButtonRight = nil
	p.liquidityContainer.SetPair(p.pair, p.token1, p.token2)
}

func (p *PageDEXRemLiquidity) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageDEXRemLiquidity) submitForm() error {
	percent, err := strconv.ParseFloat(p.txtPercent.Value(), 64)
	if err != nil {
		return err
	}

	if percent <= 0.0 || percent > 100.0 {
		return fmt.Errorf("amount must be > 0.0 and <= 100.0")
	}

	if p.share <= 0 {
		return fmt.Errorf("you don't have any liquidity")
	}

	remShares := uint64(float64(p.share) * percent / 100.0)
	pairSCID := crypto.HashHexToHash(p.pair.SCID)
	build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(addr string, open func(txPayload build_tx_modal.TxPayload)) {
		open(build_tx_modal.TxPayload{
			SCArgs: rpc.Arguments{
				{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
				{Name: rpc.SCID, DataType: rpc.DataHash, Value: pairSCID},
				{Name: "entrypoint", DataType: rpc.DataString, Value: "RemoveLiquidity"},
			},
			Transfers: []rpc.Transfer{
				rpc.Transfer{SCID: pairSCID, Burn: remShares, Destination: addr},
			},
			Ringsize:   2,
			TokensInfo: []*wallet_manager.Token{p.token1},
		})
	})

	return nil
}

func (p *PageDEXRemLiquidity) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonRemove.Clicked() {
		go func() {
			err := p.submitForm()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, 0)
			}
		}()
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.liquidityContainer.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.txtPercent.Layout(gtx, th, lang.Translate("Percentage"), "0.0")
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		percent, _ := strconv.ParseFloat(p.txtPercent.Value(), 64)

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				share := p.pair.CalcShare(p.pair.Asset1, p.share)
				share = uint64(float64(share) * percent / 100.0)
				amount := utils.ShiftNumber{Number: share, Decimals: int(p.token1.Decimals)}
				return p.infoRows[0].Layout(gtx, th, p.token1.Symbol.String, amount.Format())
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				share := p.pair.CalcShare(p.pair.Asset2, p.share)
				share = uint64(float64(share) * percent / 100.0)
				amount := utils.ShiftNumber{Number: share, Decimals: int(p.token2.Decimals)}
				return p.infoRows[1].Layout(gtx, th, p.token2.Symbol.String, amount.Format())
			}),
		)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonRemove.Style.Colors = theme.Current.ButtonPrimaryColors
		p.buttonRemove.Text = lang.Translate("REMOVE")
		return p.buttonRemove.Layout(gtx, th)
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
