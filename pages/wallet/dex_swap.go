package page_wallet

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/f32"
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
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
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

type PageDEXSwap struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonSwap              *components.Button
	infoRows                []*prefabs.InfoRow
	buttonOpenMenu          *components.Button
	pairTokenInputContainer *PairTokenInputContainer

	pair dex_sc.Pair

	slip         float64
	amountString string
	fee          uint64

	list *widget.List
}

var _ router.Page = &PageDEXSwap{}

func NewPageDEXSwap() *PageDEXSwap {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	buttonOpenMenu := components.NewButton(components.ButtonStyle{
		Icon:        navIcon,
		LoadingIcon: navIcon,
	})

	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	swapIcon, _ := widget.NewIcon(app_icons.Swap)
	buttonSwap := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        swapIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		LoadingIcon: loadingIcon,
		Animation:   components.NewButtonAnimationDefault(),
	})
	buttonSwap.Label.Alignment = text.Middle
	buttonSwap.Style.Font.Weight = font.Bold

	pairTokenInputContainer := NewPairTokenInputContainer()
	pairTokenInputContainer.txtAmount2.Editor().ReadOnly = true

	return &PageDEXSwap{
		animationEnter:          animationEnter,
		animationLeave:          animationLeave,
		list:                    list,
		infoRows:                prefabs.NewInfoRows(7),
		buttonOpenMenu:          buttonOpenMenu,
		buttonSwap:              buttonSwap,
		pairTokenInputContainer: pairTokenInputContainer,
	}
}

func (p *PageDEXSwap) IsActive() bool {
	return p.isActive
}

func (p *PageDEXSwap) Enter() {
	p.isActive = true

	if !page_instance.header.IsHistory(PAGE_DEX_SWAP) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
	page_instance.header.Title = func() string {
		return lang.Translate("DEX Swap")
	}

	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), p.pair.Symbol)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.ButtonRight = p.buttonOpenMenu
}

func (p *PageDEXSwap) SetPair(pair dex_sc.Pair, token1 *wallet_manager.Token, token2 *wallet_manager.Token) {
	p.pair = pair
	p.pairTokenInputContainer.SetTokens(token1, token2)

	page_instance.pageDEXAddLiquidity.SetPair(pair, token1, token2)

	page_instance.pageDEXRemLiquidity.token1 = token1
	page_instance.pageDEXRemLiquidity.token2 = token2
	page_instance.pageDEXRemLiquidity.pair = pair
}

func (p *PageDEXSwap) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageDEXSwap) submitForm() error {
	txtAmount1 := p.pairTokenInputContainer.txtAmount1
	token1 := p.pairTokenInputContainer.token1

	if p.pair.Liquidity1 == 0 || p.pair.Liquidity2 == 0 {
		return fmt.Errorf("pair has not liquidity")
	}

	if p.slip > 40.0 {
		return fmt.Errorf("slippage is too high")
	}

	amount := &utils.ShiftNumber{Decimals: int(token1.Decimals)}
	err := amount.Parse(txtAmount1.Value())
	if err != nil {
		return err
	}

	build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(addr string, open func(txPayload build_tx_modal.TxPayload)) {
		open(build_tx_modal.TxPayload{
			SCArgs: rpc.Arguments{
				{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
				{Name: rpc.SCID, DataType: rpc.DataHash, Value: crypto.HashHexToHash(p.pair.SCID)},
				{Name: "entrypoint", DataType: rpc.DataString, Value: "Swap"},
			},
			Transfers: []rpc.Transfer{
				rpc.Transfer{SCID: token1.GetHash(), Burn: amount.Number, Destination: addr},
			},
			Ringsize:   2,
			TokensInfo: []*wallet_manager.Token{token1},
		})
	})

	return nil
}

func (p *PageDEXSwap) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonOpenMenu.Clicked() {
		go func() {
			keyChan := listselect_modal.Instance.Open([]*listselect_modal.SelectListItem{
				listselect_modal.NewSelectListItem("add_liquidity",
					listselect_modal.NewItemText(nil, lang.Translate("Add Liquidity")).Layout,
				),
				listselect_modal.NewSelectListItem("rem_liquidity",
					listselect_modal.NewItemText(nil, lang.Translate("Remove Liquidity")).Layout,
				),
			})
			for key := range keyChan {
				switch key {
				case "add_liquidity":
					page_instance.pageRouter.SetCurrent(PAGE_DEX_ADD_LIQUIDITY)
					page_instance.header.AddHistory(PAGE_DEX_ADD_LIQUIDITY)
				case "rem_liquidity":
					page_instance.pageRouter.SetCurrent(PAGE_DEX_REM_LIQUIDITY)
					page_instance.header.AddHistory(PAGE_DEX_REM_LIQUIDITY)
				}
			}
		}()
	}

	if p.buttonSwap.Clicked() {
		go func() {
			err := p.submitForm()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, 0)
			}
		}()
	}

	txtAmount1 := p.pairTokenInputContainer.txtAmount1
	txtAmount2 := p.pairTokenInputContainer.txtAmount2
	token1 := p.pairTokenInputContainer.token1
	token2 := p.pairTokenInputContainer.token2

	if txtAmount1.Value() != p.amountString {
		p.amountString = txtAmount1.Value()
		amount1 := utils.ShiftNumber{Decimals: int(token1.Decimals)}
		amount1.Parse(p.amountString)

		receive, fee, slip := p.pair.CalcSwap(amount1.Number, p.pairTokenInputContainer.reversed)
		p.slip = slip
		p.fee = fee

		amount2 := utils.ShiftNumber{Number: receive, Decimals: int(token2.Decimals)}
		txtAmount2.SetValue(amount2.Format())
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.pairTokenInputContainer.Layout(gtx, th, lang.Translate("SEND ({})"), lang.Translate("RECEIVE ({})"))
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(20), Top: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.infoRows[0].Layout(gtx, th, lang.Translate("Slippage"), fmt.Sprintf("%.2f%%", p.slip))
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					amount := utils.ShiftNumber{Number: p.fee, Decimals: int(token2.Decimals)}
					return p.infoRows[1].Layout(gtx, th, lang.Translate("Fee"), fmt.Sprintf("%s %s", amount.Format(), token2.Symbol.String))
				}),
			)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := lang.Translate("Liquidity ({})")
					title = strings.Replace(title, "{}", token1.Symbol.String, -1)

					liquidity := p.pair.Liquidity1
					if token1.SCID != p.pair.Asset1 {
						liquidity = p.pair.Liquidity2
					}

					amount := utils.ShiftNumber{Number: liquidity, Decimals: int(token1.Decimals)}
					return p.infoRows[2].Layout(gtx, th, title, amount.Format())
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := lang.Translate("Liquidity ({})")
					title = strings.Replace(title, "{}", token2.Symbol.String, -1)

					liquidity := p.pair.Liquidity1
					if token2.SCID != p.pair.Asset1 {
						liquidity = p.pair.Liquidity2
					}

					amount := utils.ShiftNumber{Number: liquidity, Decimals: int(token2.Decimals)}
					return p.infoRows[3].Layout(gtx, th, title, amount.Format())
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					one := utils.ShiftNumber{Decimals: int(token1.Decimals)}
					one.Parse("1")
					amt, fee, _ := p.pair.CalcSwap(one.Number, p.pairTokenInputContainer.reversed)
					amount2 := utils.ShiftNumber{Number: amt + fee, Decimals: int(token2.Decimals)}
					txt := fmt.Sprintf("1 %s = %s %s", token1.Symbol.String, amount2.Format(), token2.Symbol.String)
					return p.infoRows[4].Layout(gtx, th, lang.Translate("Rate"), txt)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					txt := fmt.Sprintf("%.2f%%", float64(p.pair.Fee)/100)
					return p.infoRows[5].Layout(gtx, th, lang.Translate("Dex Fee"), txt)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.infoRows[6].Layout(gtx, th, lang.Translate("Swap Count"), fmt.Sprint(p.pair.SwapCount))
				}),
			)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonSwap.Style.Colors = theme.Current.ButtonPrimaryColors
		p.buttonSwap.Text = lang.Translate("SWAP")
		return p.buttonSwap.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return widgets[index](gtx)
		})
	})
}

type PairTokenInputContainer struct {
	token1      *wallet_manager.Token
	tokenImage1 *components.Image
	token2      *wallet_manager.Token
	tokenImage2 *components.Image

	txtAmount1   *prefabs.TextField
	txtAmount2   *prefabs.TextField
	buttonSwitch *components.Button

	reversed bool
}

func NewPairTokenInputContainer() *PairTokenInputContainer {
	txtAmount1 := prefabs.NewNumberTextField()
	txtAmount1.Input.TextSize = unit.Sp(18)
	txtAmount1.Input.FontWeight = font.Bold
	txtAmount2 := prefabs.NewNumberTextField()
	txtAmount2.Input.TextSize = unit.Sp(18)
	txtAmount2.Input.FontWeight = font.Bold
	//txtAmount2.Editor().ReadOnly = true

	tokenImage1 := &components.Image{
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}
	tokenImage2 := &components.Image{
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}

	arrowDown, _ := widget.NewIcon(icons.NavigationArrowDownward)
	buttonSwitch := components.NewButton(components.ButtonStyle{
		Icon:      arrowDown,
		Animation: components.NewButtonAnimationScale(.98),
		Inset: layout.Inset{
			Top: unit.Dp(8), Bottom: unit.Dp(8),
			Left: unit.Dp(8), Right: unit.Dp(8),
		},
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(1),
			CornerRadius: unit.Dp(5),
		},
		Rounded: components.UniformRounded(5),
	})

	return &PairTokenInputContainer{
		txtAmount1:   txtAmount1,
		txtAmount2:   txtAmount2,
		tokenImage1:  tokenImage1,
		tokenImage2:  tokenImage2,
		buttonSwitch: buttonSwitch,
	}
}

func (p *PairTokenInputContainer) SetTokens(token1 *wallet_manager.Token, token2 *wallet_manager.Token) {
	p.reversed = false
	p.token1 = token1
	p.token2 = token2
	p.txtAmount1.SetValue("0")
	p.txtAmount2.SetValue("0")
}

func (p *PairTokenInputContainer) Layout(gtx layout.Context, th *material.Theme, title1 string, title2 string) layout.Dimensions {
	if p.buttonSwitch.Clicked() {
		token1 := p.token1
		p.token1 = p.token2
		p.token2 = token1
		p.txtAmount1.SetValue("0")
		p.txtAmount2.SetValue("0")
		p.reversed = !p.reversed
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.token1.Symbol.Valid {
				title1 = strings.Replace(title1, "{}", p.token1.Symbol.String, -1)
			}

			dims := p.txtAmount1.Layout(gtx, th, title1, "")

			layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				trans := f32.Affine2D{}.Offset(f32.Pt(float32(gtx.Dp(-10)), float32(gtx.Dp(39))))
				defer op.Affine(trans).Push(gtx.Ops).Pop()
				gtx.Constraints.Max = image.Pt(gtx.Dp(35), gtx.Dp(35))
				p.tokenImage1.Src = p.token1.LoadImageOp()
				return p.tokenImage1.Layout(gtx)
			})
			return dims
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				trans := f32.Affine2D{}.Offset(f32.Pt(0, float32(gtx.Dp(14))))
				defer op.Affine(trans).Push(gtx.Ops).Pop()

				gtx.Constraints.Max = image.Pt(gtx.Dp(40), gtx.Dp(40))
				p.buttonSwitch.Style.Colors = theme.Current.ButtonSecondaryColors
				return p.buttonSwitch.Layout(gtx, th)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.token2.Symbol.Valid {
				title2 = strings.Replace(title2, "{}", p.token2.Symbol.String, -1)
			}

			dims := p.txtAmount2.Layout(gtx, th, title2, "")

			layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				trans := f32.Affine2D{}.Offset(f32.Pt(float32(gtx.Dp(-10)), float32(gtx.Dp(39))))
				defer op.Affine(trans).Push(gtx.Ops).Pop()
				gtx.Constraints.Max = image.Pt(gtx.Dp(35), gtx.Dp(35))
				p.tokenImage2.Src = p.token2.LoadImageOp()
				return p.tokenImage2.Layout(gtx)
			})

			return dims
		}),
	)
}
