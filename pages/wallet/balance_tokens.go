package page_wallet

import (
	"image"
	"image/color"
	"log"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/assets"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageBalanceTokens struct {
	isActive   bool
	firstEnter bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	notRegistered   *NotRegistered
	displayBalance  *DisplayBalance
	tokensContainer *TokensContainer
	tokenItems      []TokenListItem
}

var _ router.Container = &PageBalanceTokens{}

func NewPageBalanceTokens() *PageBalanceTokens {
	th := app_instance.Current.Theme

	img, err := assets.GetImage("hobo48.jpg")
	if err != nil {
		log.Fatal(err)
	}

	tokenItems := []TokenListItem{}
	for i := 0; i < 10; i++ {
		tokenItems = append(tokenItems, TokenListItem{
			tokenImageItem: NewTokenImageItem(img),
			tokenName:      "Dero",
			tokenId:        "00000...00000",
			tokenBalance:   "342.35546",
			Clickable:      new(widget.Clickable),
			//ImageClickable: new(widget.Clickable),
		})
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	return &PageBalanceTokens{
		displayBalance:  NewDisplayBalance(th),
		tokensContainer: NewTokenContainer(th),
		notRegistered:   NewNotRegistered(),
		tokenItems:      tokenItems,
		firstEnter:      true,
		animationEnter:  animationEnter,
		animationLeave:  animationLeave,
	}
}

func (p *PageBalanceTokens) IsActive() bool {
	return p.isActive
}

func (p *PageBalanceTokens) Enter() {
	p.isActive = true

	if !p.firstEnter {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.firstEnter = false
}

func (p *PageBalanceTokens) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageBalanceTokens) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.notRegistered.Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return p.displayBalance.Layout(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return p.tokensContainer.Layout(gtx, th, p.tokenItems)
		}),
	)
}

type NotRegistered struct {
	iconWarning *widget.Icon
}

func NewNotRegistered() *NotRegistered {
	iconWarning, _ := widget.NewIcon(icons.AlertWarning)
	return &NotRegistered{
		iconWarning: iconWarning,
	}
}

func (n *NotRegistered) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	border := widget.Border{Color: color.NRGBA{A: 100}, CornerRadius: unit.Dp(5), Width: unit.Dp(1)}

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(20),
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return n.iconWarning.Layout(gtx, color.NRGBA{A: 100})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						label := material.Label(th, unit.Sp(14), "This wallet is not registered on the blockchain.")
						label.Color = color.NRGBA{A: 200}
						return label.Layout(gtx)
					}),
				)
			})
		})
	})
}

type DisplayBalance struct {
	buttonSend        *components.Button
	buttonReceive     *components.Button
	buttonHideBalance *components.Button
	labelTitle        material.LabelStyle
	labelAmount       material.LabelStyle

	hideBalanceIcon *widget.Icon
	showBalanceIcon *widget.Icon
	hiddenBalance   bool
}

func NewDisplayBalance(th *material.Theme) *DisplayBalance {
	labelTitle := material.Label(th, unit.Sp(14), "Available Balance")
	labelTitle.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}

	labelAmount := material.Label(th, unit.Sp(34), "--")
	labelAmount.Font.Weight = font.Bold

	sendIcon, _ := widget.NewIcon(icons.NavigationArrowUpward)
	buttonSend := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "SEND",
		Icon:            sendIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSend.Label.Alignment = text.Middle
	buttonSend.Style.Font.Weight = font.Bold

	receiveIcon, _ := widget.NewIcon(icons.NavigationArrowDownward)
	buttonReceive := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "RECEIVE",
		Icon:            receiveIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonReceive.Label.Alignment = text.Middle
	buttonReceive.Style.Font.Weight = font.Bold

	hideBalanceIcon, _ := widget.NewIcon(icons.ActionVisibility)
	showBalanceIcon, _ := widget.NewIcon(icons.ActionVisibilityOff)
	buttonHideBalance := components.NewButton(components.ButtonStyle{
		Icon:      hideBalanceIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	})

	return &DisplayBalance{
		labelTitle:        labelTitle,
		labelAmount:       labelAmount,
		buttonSend:        buttonSend,
		buttonReceive:     buttonReceive,
		buttonHideBalance: buttonHideBalance,
		hideBalanceIcon:   hideBalanceIcon,
		showBalanceIcon:   showBalanceIcon,
	}
}

func (d *DisplayBalance) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Left: unit.Dp(30), Right: unit.Dp(30),
		Top: unit.Dp(0), Bottom: unit.Dp(40),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return d.labelTitle.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						dims := d.labelAmount.Layout(gtx)

						if d.hiddenBalance {
							paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Rect{
								Max: dims.Size,
							}.Op())
						}

						return dims
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(50)
						gtx.Constraints.Max.X = gtx.Dp(50)

						if d.hiddenBalance {
							d.buttonHideBalance.Style.Icon = d.showBalanceIcon
						} else {
							d.buttonHideBalance.Style.Icon = d.hideBalanceIcon
						}

						if d.buttonHideBalance.Clickable.Clicked() {
							d.hiddenBalance = !d.hiddenBalance
							op.InvalidateOp{}.Add(gtx.Ops)
						}

						return d.buttonHideBalance.Layout(gtx, th)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(40)
						return d.buttonSend.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.Y = gtx.Dp(40)

						return d.buttonReceive.Layout(gtx, th)
					}),
				)
			}),
		)
	})
}

type TokensContainer struct {
	tokenList      *TokenList
	buttonAddToken *components.Button
}

func NewTokenContainer(th *material.Theme) *TokensContainer {
	textColor := color.NRGBA{R: 0, G: 0, B: 0, A: 100}
	textHoverColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

	walletIcon, _ := widget.NewIcon(icons.ContentAddBox)
	buttonAddToken := components.NewButton(components.ButtonStyle{
		Icon:           walletIcon,
		TextColor:      textColor,
		HoverTextColor: &textHoverColor,
		Animation:      components.NewButtonAnimationScale(.92),
	})

	return &TokensContainer{
		tokenList:      NewTokenList(th),
		buttonAddToken: buttonAddToken,
	}
}

func (t *TokensContainer) Layout(gtx layout.Context, th *material.Theme, items []TokenListItem) layout.Dimensions {
	dr := image.Rectangle{Max: gtx.Constraints.Min}
	paint.LinearGradientOp{
		Stop1:  f32.Pt(0, float32(dr.Min.Y)),
		Stop2:  f32.Pt(0, float32(dr.Max.Y)),
		Color1: color.NRGBA{R: 0, G: 0, B: 0, A: 5},
		Color2: color.NRGBA{R: 0, G: 0, B: 0, A: 50},
	}.Add(gtx.Ops)
	cl := clip.Rect(dr).Push(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	cl = clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, gtx.Dp(1))}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{R: 0, G: 0, B: 0, A: 50}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	return layout.Inset{
		Left: unit.Dp(30), Right: unit.Dp(30),
		Top: unit.Dp(30), Bottom: unit.Dp(0),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						labelTokens := material.Label(th, unit.Sp(16), "YOUR TOKENS")
						labelTokens.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 200}
						labelTokens.Font.Weight = font.Bold
						return labelTokens.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(30)
						gtx.Constraints.Min.Y = gtx.Dp(30)
						return t.buttonAddToken.Layout(gtx, th)
					}),
				)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(20),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return t.tokenList.Layout(gtx, th, items)
				})
			}),
		)
	})
}

type TokenList struct {
	listStyle material.ListStyle
}

func NewTokenList(th *material.Theme) *TokenList {
	list := new(widget.List)
	list.Axis = layout.Vertical

	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay
	listStyle.Indicator.MinorWidth = unit.Dp(10)
	listStyle.Indicator.CornerRadius = unit.Dp(5)
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	listStyle.Indicator.Color = black
	//listStyle.Indicator.HoverColor = f32color.Hovered(black)

	return &TokenList{
		listStyle: listStyle,
	}
}

func (l *TokenList) Layout(gtx layout.Context, th *material.Theme, items []TokenListItem) layout.Dimensions {
	return layout.UniformInset(unit.Dp(0)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return l.listStyle.Layout(gtx, len(items), func(gtx layout.Context, i int) layout.Dimensions {
			return items[i].Layout(gtx, th)
		})
	})

	/*
		bounds := image.Rect(0, 0, d.Size.X, d.Size.Y)
		rectPath := clip.UniformRRect(bounds, 10).Path(gtx.Ops)
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			clip.Stroke{
				Path:  rectPath,
				Width: 4,
			}.Op(),
		)*/

	//return d
}

type TokenImageItem struct {
	Image     *components.Image
	Clickable *widget.Clickable

	AnimationEnter   *animation.Animation
	AnimationLeave   *animation.Animation
	hoverSwitchState bool
}

func NewTokenImageItem(src image.Image) *TokenImageItem {
	image := &components.Image{
		Src:     paint.NewImageOp(src),
		Rounded: unit.Dp(10),
		Fit:     components.Cover,
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 1.1, .1, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1.1, 1, .1, ease.Linear),
	))

	return &TokenImageItem{
		Image:          image,
		Clickable:      new(widget.Clickable),
		AnimationEnter: animationEnter,
		AnimationLeave: animationLeave,
	}
}

func (item *TokenImageItem) Layout(gtx layout.Context) layout.Dimensions {
	return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		{
			state := item.AnimationEnter.Update(gtx)
			if state.Active {
				item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
					pt := dims.Size.Div(2)
					origin := f32.Pt(float32(pt.X), float32(pt.Y))
					return trans.Scale(origin, f32.Pt(state.Value, state.Value))
				}
			}
		}

		{
			state := item.AnimationLeave.Update(gtx)
			if state.Active {
				item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
					pt := dims.Size.Div(2)
					origin := f32.Pt(float32(pt.X), float32(pt.Y))
					return trans.Scale(origin, f32.Pt(state.Value, state.Value))
				}
			}
		}

		/*
			v1, a1, _ := item.AnimationEnter.Update(gtx)
			v2, a2, _ := item.AnimationLeave.Update(gtx)

			item.Image.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
				pt := dims.Size.Div(2)
				origin := f32.Pt(float32(pt.X), float32(pt.Y))

				if a1 {
					trans = trans.Scale(origin, f32.Pt(v1, v1))
				}

				if a2 {
					trans = trans.Scale(origin, f32.Pt(v2, v2))
				}

				return trans
			}*/

		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}

		if item.Clickable.Hovered() && !item.hoverSwitchState {
			item.hoverSwitchState = true
			item.AnimationEnter.Start()
			item.AnimationLeave.Reset()
		}

		if !item.Clickable.Hovered() && item.hoverSwitchState {
			item.hoverSwitchState = false
			item.AnimationLeave.Start()
			item.AnimationEnter.Reset()
		}

		return item.Image.Layout(gtx)
	})
}

type TokenListItem struct {
	tokenName    string
	tokenId      string
	tokenBalance string

	Clickable      *widget.Clickable
	tokenImageItem *TokenImageItem
}

func (item *TokenListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	m := op.Record(gtx.Ops)
	dims := item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(50)
					gtx.Constraints.Max.Y = gtx.Dp(50)
					return item.tokenImageItem.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Label(th, unit.Sp(18), item.tokenName)
									label.Font.Weight = font.Bold
									return label.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									label := material.Label(th, unit.Sp(14), item.tokenId)
									label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
									return label.Layout(gtx)
								}),
							)
						}),
						layout.Flexed(1, layout.Spacer{}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							label := material.Label(th, unit.Sp(18), item.tokenBalance)
							label.Font.Weight = font.Bold
							return label.Layout(gtx)
						}),
					)
				}),
			)
		})
	})
	c := m.Stop()

	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		clip.RRect{
			Rect: image.Rectangle{Max: dims.Size},
			SE:   gtx.Dp(10),
			NW:   gtx.Dp(10),
			NE:   gtx.Dp(10),
			SW:   gtx.Dp(10),
		}.Op(gtx.Ops))

	c.Add(gtx.Ops)

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return dims
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
	)

	/*
		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
			bounds := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
				clip.UniformRRect(bounds, 10).Op(gtx.Ops),
			)
		}*/
}
