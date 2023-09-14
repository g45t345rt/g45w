package page_wallet

import (
	"fmt"
	"image"
	"strings"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageDEXPairs struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	list  *widget.List
	items []*DexPairItem
}

var _ router.Page = &PageDEXPairs{}

func NewPageDEXPairs() *PageDEXPairs {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageDEXPairs{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
	}
}

func (p *PageDEXPairs) IsActive() bool {
	return p.isActive
}

func (p *PageDEXPairs) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_DEX_PAIRS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
	page_instance.header.Title = func() string {
		return lang.Translate("DEX Pairs")
	}

	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil

	p.Load()
}

func (p *PageDEXPairs) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageDEXPairs) Load() error {
	p.items = make([]*DexPairItem, 0)
	// Keystore
	// 8088b0089725de1d323276a0daa1f25cfab9c0b68ccb9318cbf6bf83f5a127c1

	// dex.swap.registry
	// a6b36e8a23d153c5f09683183fc1059285476a1ce3f7f53952ab67b4fa34bcce

	var result rpc.GetSC_Result
	err := walletapi.RPC_Client.Call("DERO.GetSC", rpc.GetSC_Params{
		SCID:      "a6b36e8a23d153c5f09683183fc1059285476a1ce3f7f53952ab67b4fa34bcce",
		Code:      false,
		Variables: true,
	}, &result)
	if err != nil {
		return err
	}

	for key, value := range result.VariableStringKeys {
		k := strings.Split(key, ":")
		prefix := k[0]

		if prefix == "p" {
			//symbol1 := k[1]
			//symbol2 := k[2]

			scId := value.(string)

			var result rpc.GetSC_Result
			err := walletapi.RPC_Client.Call("DERO.GetSC", rpc.GetSC_Params{
				SCID:      scId,
				Code:      false,
				Variables: true,
			}, &result)
			if err != nil {
				continue
			}

			pair := dex_sc.Pair{}
			err = pair.Parse(scId, result.VariableStringKeys)
			if err != nil {
				continue
			}

			token1, err := wallet_manager.GetTokenBySCID(pair.Asset1)
			if err != nil {
				return err
			}

			token2, err := wallet_manager.GetTokenBySCID(pair.Asset2)
			if err != nil {
				return err
			}

			p.items = append(p.items, NewDexPairItem(pair, token1, token2))
		}
	}

	return err
}

func (p *PageDEXPairs) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(20), "TLV: 307182.20 USDT")
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		childs := []layout.FlexChild{}
		for i := range p.items {
			idx := i
			childs = append(childs,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					item := p.items[idx]
					if item.clickable.Clicked() {
						page_instance.pageDexSwap.pair = item.pair
						page_instance.pageDexSwap.token1 = item.token1
						page_instance.pageDexSwap.token2 = item.token2
						page_instance.pageRouter.SetCurrent(PAGE_DEX_SWAP)
						page_instance.header.AddHistory(PAGE_DEX_SWAP)
					}

					return item.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			)
		}

		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
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

type DexPairItem struct {
	pair        dex_sc.Pair
	token1      *wallet_manager.Token
	tokenImage1 *components.Image
	token2      *wallet_manager.Token
	tokenImage2 *components.Image
	clickable   *widget.Clickable
}

func NewDexPairItem(pair dex_sc.Pair, token1 *wallet_manager.Token, token2 *wallet_manager.Token) *DexPairItem {
	return &DexPairItem{
		pair:   pair,
		token1: token1,
		tokenImage1: &components.Image{
			Fit:     components.Cover,
			Rounded: components.UniformRounded(unit.Dp(5)),
		},
		token2: token2,
		tokenImage2: &components.Image{
			Fit:     components.Cover,
			Rounded: components.UniformRounded(unit.Dp(5)),
		},
		clickable: new(widget.Clickable),
	}
}

func (item *DexPairItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	m := op.Record(gtx.Ops)
	dims := layout.Inset{
		Top: unit.Dp(13), Bottom: unit.Dp(13),
		Left: unit.Dp(15), Right: unit.Dp(15),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
											item.tokenImage1.Src = item.token1.LoadImageOp()
											return item.tokenImage1.Layout(gtx)
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
										layout.Rigid(func(gtx layout.Context) layout.Dimensions {
											gtx.Constraints.Max = image.Pt(gtx.Dp(20), gtx.Dp(20))
											item.tokenImage2.Src = item.token2.LoadImageOp()
											return item.tokenImage2.Layout(gtx)
										}),
										layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
										layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
											lbl := material.Label(th, unit.Sp(18), item.pair.Symbol)
											lbl.Font.Weight = font.Bold
											return lbl.Layout(gtx)
										}),
									)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									d1 := utils.ShiftNumber{Number: item.pair.Liquidity1, Decimals: int(item.token1.Decimals)}
									d2 := utils.ShiftNumber{Number: item.pair.Liquidity2, Decimals: int(item.token2.Decimals)}
									txt := fmt.Sprintf("%s / %s", d1.Format(), d2.Format())
									lbl := material.Label(th, unit.Sp(14), txt)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				})
			}),
		)
	})
	c := m.Stop()

	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
		paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(10),
			).Op(gtx.Ops),
		)
	} else {
		paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
			clip.RRect{
				Rect: image.Rectangle{Max: dims.Size},
				NW:   gtx.Dp(10), NE: gtx.Dp(10),
				SE: gtx.Dp(10), SW: gtx.Dp(10),
			}.Op(gtx.Ops),
		)
	}

	c.Add(gtx.Ops)
	return dims
}
