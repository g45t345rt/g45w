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

	txtAmount1     *prefabs.TextField
	txtAmount2     *prefabs.TextField
	buttonSwap     *components.Button
	buttonSwitch   *components.Button
	infoRows       []*prefabs.InfoRow
	buttonOpenMenu *components.Button

	pair        dex_sc.Pair
	token1      *wallet_manager.Token
	tokenImage1 *components.Image
	token2      *wallet_manager.Token
	tokenImage2 *components.Image

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
		Icon:      navIcon,
		Animation: components.NewButtonAnimationScale(.98),
	})

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

	var infoRows []*prefabs.InfoRow
	for i := 0; i < 5; i++ {
		infoRows = append(infoRows, prefabs.NewInfoRow())
	}

	txtAmount1 := prefabs.NewNumberTextField()
	txtAmount1.Input.TextSize = unit.Sp(18)
	txtAmount1.Input.FontWeight = font.Bold
	txtAmount2 := prefabs.NewNumberTextField()
	txtAmount2.Input.TextSize = unit.Sp(18)
	txtAmount2.Input.FontWeight = font.Bold
	txtAmount2.Editor().ReadOnly = true

	tokenImage1 := &components.Image{
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}
	tokenImage2 := &components.Image{
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}

	return &PageDEXSwap{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		infoRows:       infoRows,
		buttonOpenMenu: buttonOpenMenu,
		txtAmount1:     txtAmount1,
		txtAmount2:     txtAmount2,
		buttonSwap:     buttonSwap,
		buttonSwitch:   buttonSwitch,
		tokenImage1:    tokenImage1,
		tokenImage2:    tokenImage2,
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
		lbl := material.Label(th, unit.Sp(14), lang.Translate("DERO:DST"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.ButtonRight = p.buttonOpenMenu

	p.Load()
}

func (p *PageDEXSwap) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageDEXSwap) Load() error {

	return nil
}

func (p *PageDEXSwap) CalcSwap(amt uint64) (receive uint64, fee uint64, slip float64) {
	pair := p.pair

	if pair.Asset1 == p.token2.SCID {
		receiveAmt := float64(amt) * float64(pair.Liquidity1) / float64(pair.Liquidity2+amt)
		receiveAmtMinusFee := receiveAmt * float64(10000-pair.Fee) / float64(10000)
		receive = uint64(receiveAmtMinusFee)
		fee = uint64(receiveAmt) - uint64(receiveAmtMinusFee)
		slip = 100.0 - (1.0 / (1.0 + float64(amt)/float64(pair.Liquidity2)) * 100.0)
	} else {
		receiveAmt := float64(amt) * float64(pair.Liquidity2) / float64(pair.Liquidity1+amt)
		receiveAmtMinusFee := receiveAmt * float64(10000-pair.Fee) / float64(10000)
		receive = uint64(receiveAmtMinusFee)
		fee = uint64(receiveAmt) - uint64(receiveAmtMinusFee)
		slip = 100.0 - (1.0 / (1.0 + float64(amt)/float64(pair.Liquidity1)) * 100.0)
	}

	return
}

func (p *PageDEXSwap) submitForm() error {
	if p.pair.Liquidity1 == 0 || p.pair.Liquidity2 == 0 {
		return fmt.Errorf("pair has not liquidity")
	}

	amount := &utils.ShiftNumber{Decimals: int(p.token1.Decimals)}
	err := amount.Parse(p.txtAmount1.Value())
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
				rpc.Transfer{SCID: p.token1.GetHash(), Burn: amount.Number, Destination: addr},
			},
			Ringsize:   2,
			TokensInfo: []*wallet_manager.Token{p.token1},
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
				case "rem_liquidity":
				}
			}
		}()
	}

	if p.buttonSwitch.Clicked() {
		token1 := p.token1
		p.token1 = p.token2
		p.token2 = token1
		p.txtAmount1.SetValue("0")
		p.txtAmount2.SetValue("0")
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

	if p.txtAmount1.Value() != p.amountString {
		p.amountString = p.txtAmount1.Value()
		amount1 := utils.ShiftNumber{Decimals: int(p.token1.Decimals)}
		amount1.Parse(p.amountString)

		receive, fee, slip := p.CalcSwap(amount1.Number)
		p.slip = slip
		p.fee = fee

		amount2 := utils.ShiftNumber{Number: receive, Decimals: int(p.token2.Decimals)}
		p.txtAmount2.SetValue(amount2.Format())
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		txt := lang.Translate("SEND ({})")
		if p.token1.Symbol.Valid {
			txt = strings.Replace(txt, "{}", p.token1.Symbol.String, -1)
		}

		dims := p.txtAmount1.Layout(gtx, th, txt, "")

		layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			trans := f32.Affine2D{}.Offset(f32.Pt(float32(gtx.Dp(-10)), float32(gtx.Dp(39))))
			defer op.Affine(trans).Push(gtx.Ops).Pop()
			gtx.Constraints.Max = image.Pt(gtx.Dp(35), gtx.Dp(35))
			p.tokenImage1.Src = p.token1.LoadImageOp()
			return p.tokenImage1.Layout(gtx)
		})
		return dims
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			trans := f32.Affine2D{}.Offset(f32.Pt(0, float32(gtx.Dp(14))))
			defer op.Affine(trans).Push(gtx.Ops).Pop()

			gtx.Constraints.Max = image.Pt(gtx.Dp(40), gtx.Dp(40))
			p.buttonSwitch.Style.Colors = theme.Current.ButtonSecondaryColors
			return p.buttonSwitch.Layout(gtx, th)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					txt := lang.Translate("RECEIVE ({})")
					if p.token2.Symbol.Valid {
						txt = strings.Replace(txt, "{}", p.token2.Symbol.String, -1)
					}

					txt = strings.Replace(txt, "{}", txt, -1)
					dims := p.txtAmount2.Layout(gtx, th, txt, "")

					layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						trans := f32.Affine2D{}.Offset(f32.Pt(float32(gtx.Dp(-10)), float32(gtx.Dp(39))))
						defer op.Affine(trans).Push(gtx.Ops).Pop()
						gtx.Constraints.Max = image.Pt(gtx.Dp(35), gtx.Dp(35))
						p.tokenImage2.Src = p.token2.LoadImageOp()
						return p.tokenImage2.Layout(gtx)
					})

					return dims
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					amount := utils.ShiftNumber{Number: p.fee, Decimals: int(p.token2.Decimals)}
					txt := lang.Translate("Slip: {0} Fee: {1}")
					txt = strings.Replace(txt, "{0}", fmt.Sprintf("%.2f%%", p.slip), -1)
					txt = strings.Replace(txt, "{1}", fmt.Sprintf("%s %s", amount.Format(), p.token2.Symbol.String), -1)
					lbl := material.Label(th, unit.Sp(16), txt)
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		})
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(20)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := lang.Translate("Liquidity ({})")
					title = strings.Replace(title, "{}", p.token1.Symbol.String, -1)

					liquidity := p.pair.Liquidity1
					if p.token1.SCID != p.pair.Asset1 {
						liquidity = p.pair.Liquidity2
					}

					amount := utils.ShiftNumber{Number: liquidity, Decimals: int(p.token1.Decimals)}
					return p.infoRows[0].Layout(gtx, th, title, amount.Format())
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					title := lang.Translate("Liquidity ({})")
					title = strings.Replace(title, "{}", p.token2.Symbol.String, -1)

					liquidity := p.pair.Liquidity1
					if p.token2.SCID != p.pair.Asset1 {
						liquidity = p.pair.Liquidity2
					}

					amount := utils.ShiftNumber{Number: liquidity, Decimals: int(p.token2.Decimals)}
					return p.infoRows[1].Layout(gtx, th, title, amount.Format())
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					amount1 := utils.ShiftNumber{Decimals: int(p.token1.Decimals)}
					amount1.Parse("1")
					amt, fee, _ := p.CalcSwap(amount1.Number)
					amount2 := utils.ShiftNumber{Number: amt + fee, Decimals: int(p.token2.Decimals)}
					txt := fmt.Sprintf("1 %s = %s %s", p.token1.Symbol.String, amount2.Format(), p.token2.Symbol.String)
					return p.infoRows[2].Layout(gtx, th, lang.Translate("Rate"), txt)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					txt := fmt.Sprintf("%.2f%%", float64(p.pair.Fee)/100)
					return p.infoRows[3].Layout(gtx, th, lang.Translate("Dex Fee"), txt)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.infoRows[4].Layout(gtx, th, lang.Translate("Swap Count"), fmt.Sprint(p.pair.SwapCount))
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
