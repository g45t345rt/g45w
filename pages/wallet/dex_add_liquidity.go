package page_wallet

import (
	"fmt"
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	crypto "github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageDEXAddLiquidity struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonAdd               *components.Button
	liquidityContainer      *LiquidityContainer
	pairTokenInputContainer *PairTokenInputContainer
	infoRows                []*prefabs.InfoRow

	pair   dex_sc.Pair
	amount string

	list *widget.List
}

var _ router.Page = &PageDEXAddLiquidity{}

func NewPageDEXAddLiquidity() *PageDEXAddLiquidity {

	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAddBox)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      addIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonAdd.Label.Alignment = text.Middle
	buttonAdd.Style.Font.Weight = font.Bold

	pairTokenInputContainer := NewPairTokenInputContainer()
	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_DEX_ADD_LIQUIDITY)
	return &PageDEXAddLiquidity{
		headerPageAnimation:     headerPageAnimation,
		list:                    list,
		pairTokenInputContainer: pairTokenInputContainer,
		buttonAdd:               buttonAdd,
		liquidityContainer:      NewLiquidityContainer(),
		infoRows:                prefabs.NewInfoRows(1),
	}
}

func (p *PageDEXAddLiquidity) IsActive() bool {
	return p.isActive
}

func (p *PageDEXAddLiquidity) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string {
		return lang.Translate("Add Liquidity")
	}

	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), p.pair.Symbol)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.RightLayout = nil
}

func (p *PageDEXAddLiquidity) SetPair(pair dex_sc.Pair, token1 *wallet_manager.Token, token2 *wallet_manager.Token) {
	p.pair = pair
	p.liquidityContainer.SetPair(pair, token1, token2)
	p.pairTokenInputContainer.SetTokens(token1, token2)

	if pair.SharesOutstanding > 0 {
		p.pairTokenInputContainer.txtAmount2.Editor().ReadOnly = true
	} else {
		p.pairTokenInputContainer.txtAmount2.Editor().ReadOnly = false
	}
}

func (p *PageDEXAddLiquidity) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageDEXAddLiquidity) submitForm() error {
	token1 := p.pairTokenInputContainer.token1
	token2 := p.pairTokenInputContainer.token2
	txtAmount1 := p.pairTokenInputContainer.txtAmount1
	txtAmount2 := p.pairTokenInputContainer.txtAmount2

	amount1 := utils.ShiftNumber{Decimals: int(token1.Decimals)}
	err := amount1.Parse(txtAmount1.Value())
	if err != nil {
		return err
	}

	amount2 := utils.ShiftNumber{Decimals: int(token2.Decimals)}
	err = amount2.Parse(txtAmount2.Value())
	if err != nil {
		return err
	}

	build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(randomAddr string) build_tx_modal.TxPayload {
		return build_tx_modal.TxPayload{
			Transfer: rpc.Transfer_Params{
				SC_RPC: rpc.Arguments{
					{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
					{Name: rpc.SCID, DataType: rpc.DataHash, Value: crypto.HashHexToHash(p.pair.SCID)},
					{Name: "entrypoint", DataType: rpc.DataString, Value: "AddLiquidity"},
				},
				Transfers: []rpc.Transfer{
					rpc.Transfer{SCID: token1.GetHash(), Burn: amount1.Number, Destination: randomAddr},
					rpc.Transfer{SCID: token2.GetHash(), Burn: amount2.Number, Destination: randomAddr},
				},
				Ringsize: 2,
			},
			TokensInfo: []*wallet_manager.Token{token1, token2},
		}
	})

	return nil
}

func (p *PageDEXAddLiquidity) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonAdd.Clicked(gtx) {
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

	if p.pair.SharesOutstanding == 0 {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("Looks like there is no shares. You will provide the initial liquidity to the pair."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})

		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("You have to enter the amount for both token and determine the rate."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})
	} else {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.liquidityContainer.Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		if p.pair.SharesOutstanding > 0 {
			txtAmount1 := p.pairTokenInputContainer.txtAmount1
			if txtAmount1.Value() != p.amount {
				p.amount = txtAmount1.Value()
				token1 := p.pairTokenInputContainer.token1
				amount1 := utils.ShiftNumber{Decimals: int(token1.Decimals)}
				amount1.Parse(p.amount)

				token2 := p.pairTokenInputContainer.token2
				txtAmount2 := p.pairTokenInputContainer.txtAmount2
				var value uint64
				if p.pairTokenInputContainer.reversed {
					value = utils.MultDiv(amount1.Number, p.pair.Liquidity1, p.pair.Liquidity2)
				} else {
					value = utils.MultDiv(amount1.Number, p.pair.Liquidity2, p.pair.Liquidity1)
				}

				amount2 := utils.ShiftNumber{Number: value, Decimals: int(token2.Decimals)}
				txtAmount2.SetValue(amount2.Format())
			}
		}

		return p.pairTokenInputContainer.Layout(gtx, th, lang.Translate("SEND {}"), lang.Translate("SEND {}"))
	})

	if p.pair.SharesOutstanding == 0 {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			token1 := p.pairTokenInputContainer.token1
			token2 := p.pairTokenInputContainer.token2

			liquidity1 := p.pairTokenInputContainer.txtAmount1.Value()
			amount1 := utils.ShiftNumber{Decimals: int(token1.Decimals)}
			amount1.Parse(liquidity1)

			one := utils.ShiftNumber{Decimals: int(token1.Decimals)}
			one.Parse("1")

			liquidity2 := p.pairTokenInputContainer.txtAmount2.Value()
			amount2 := utils.ShiftNumber{Decimals: int(token2.Decimals)}
			amount2.Parse(liquidity2)

			rate := uint64(0)
			if amount1.Number > 0 {
				rate = uint64(float64(one.Number) * float64(amount2.Number) / float64(amount1.Number))
			}

			rateAmount := utils.ShiftNumber{Number: rate, Decimals: int(token2.Decimals)}

			return p.infoRows[0].Layout(gtx, th, lang.Translate("Rate"), fmt.Sprintf("1 %s = %s %s", token1.Symbol.String, rateAmount.Format(), token2.Symbol.String))
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonAdd.Style.Colors = theme.Current.ButtonPrimaryColors
		p.buttonAdd.Text = lang.Translate("ADD")
		return p.buttonAdd.Layout(gtx, th)
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

type LiquidityContainer struct {
	pair   dex_sc.Pair
	token1 *wallet_manager.Token
	token2 *wallet_manager.Token
	share  uint64
}

func NewLiquidityContainer() *LiquidityContainer {
	return &LiquidityContainer{}
}

func (p *LiquidityContainer) SetPair(pair dex_sc.Pair, token1 *wallet_manager.Token, token2 *wallet_manager.Token) {
	p.pair = pair
	p.token1 = token1
	p.token2 = token2

	go func() {
		wallet := wallet_manager.OpenedWallet
		addr := wallet.Memory.GetAddress().String()
		p.share, _, _ = wallet.Memory.GetDecryptedBalanceAtTopoHeight(crypto.HashHexToHash(p.pair.SCID), -1, addr)
	}()
}

func (p *LiquidityContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(18), lang.Translate("Your liquidity"))
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						share := p.pair.CalcShare(p.share, false)
						amount := utils.ShiftNumber{Number: share, Decimals: int(p.token1.Decimals)}
						lbl := material.Label(th, unit.Sp(16), fmt.Sprintf("%s %s", amount.Format(), p.token1.Symbol.String))
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(1)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						share := p.pair.CalcShare(p.share, true)
						amount := utils.ShiftNumber{Number: share, Decimals: int(p.token2.Decimals)}
						lbl := material.Label(th, unit.Sp(16), fmt.Sprintf("%s %s", amount.Format(), p.token2.Symbol.String))
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(1)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						value := p.pair.CalcOwnership(p.share)
						lbl := material.Label(th, unit.Sp(16), fmt.Sprintf("Ownership: %.3f%%", value))
						lbl.Color = theme.Current.TextMuteColor
						return lbl.Layout(gtx)
					}),
				)
			}),
		)
	})
	c := r.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
		clip.UniformRRect(
			image.Rectangle{Max: dims.Size},
			gtx.Dp(15),
		).Op(gtx.Ops))

	c.Add(gtx.Ops)
	return dims
}
