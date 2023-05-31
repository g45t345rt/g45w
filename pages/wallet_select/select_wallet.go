package page_wallet_select

import (
	"fmt"
	"image"
	"image/color"

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
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSelectWallet struct {
	isActive   bool
	firstEnter bool
	clickable  *widget.Clickable

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonWalletCreate *components.Button
	walletList         *WalletList

	modalWalletPassword        *WalletPasswordModal
	modalCreateWalletSelection *CreateWalletSelectionModal
}

var _ router.Container = &PageSelectWallet{}

func NewPageSelectWallet() *PageSelectWallet {
	theme := app_instance.Current.Theme

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .5, ease.OutCubic),
	))

	walletList := NewWalletList(theme)

	//for i := 0; i < 20; i++ {
	//walletList.items = append(walletList.items,
	//	NewWalletListItem(theme, fmt.Sprintf("Wallet %d", i), "dero1qy...gr2j8u5"))
	//}

	//childRouter := router.NewRouter()
	//childRouter.Add("create_wallet_form", NewPageCreateWalletForm())
	//childRouter.Add("create_wallet_seed_form", NewPageCreateWalletSeedForm())
	//page_instance = &PageInstance{
	//	router: childRouter,
	//}

	modalWalletPassword := NewWalletPasswordModal(theme)
	modalCreateWalletSelection := NewCreateWalletSelectionModal(theme)

	router := app_instance.Current.Router
	router.PushLayout(func(gtx layout.Context, th *material.Theme) {
		modalWalletPassword.Layout(gtx)
		modalCreateWalletSelection.Layout(gtx, th)
	})

	return &PageSelectWallet{
		firstEnter: true,
		clickable:  new(widget.Clickable),

		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonWalletCreate: NewWalletCreateButton(),
		walletList:         walletList,

		modalWalletPassword:        modalWalletPassword,
		modalCreateWalletSelection: modalCreateWalletSelection,
	}
}

func (p *PageSelectWallet) IsActive() bool {
	return p.isActive
}

func (p *PageSelectWallet) Enter() {
	page_instance.header.LabelTitle.Text = "Select Wallet"
	p.isActive = true

	if !p.firstEnter {
		p.animationLeave.Reset()
		p.animationEnter.Start()
	}

	theme := app_instance.Current.Theme
	walletManager := wallet_manager.Instance
	p.walletList.items = make([]WalletListItem, 0)
	for _, wallet := range walletManager.Wallets {
		p.walletList.items = append(p.walletList.items,
			NewWalletListItem(theme,
				fmt.Sprintf("Wallet [%s]", wallet.Name),
				utils.ReduceString(wallet.Addr, 7, 7)),
		)
	}

	p.firstEnter = false
}

func (p *PageSelectWallet) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageSelectWallet) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(0), Bottom: unit.Dp(30),
				Left: unit.Dp(30), Right: unit.Dp(30),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if len(p.walletList.items) == 0 {
							labelNoWallet := material.Label(th, unit.Sp(16), "You didn't add a wallet yet.\nClick 'New Wallet' button to continue.")
							return labelNoWallet.Layout(gtx)
						} else {
							for _, item := range p.walletList.items {
								if item.Clickable.Clicked() {
									p.modalWalletPassword.modal.SetVisible(gtx, true)
								}
							}

							return p.walletList.Layout(gtx, th)
						}
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if p.buttonWalletCreate.Clickable.Clicked() {
							p.modalCreateWalletSelection.modal.SetVisible(gtx, true)
						}

						return p.buttonWalletCreate.Layout(gtx, th)
					}),
				)
			})
		}),
	)

	{
		submitted, text := p.modalWalletPassword.Submit()
		if submitted {
			if text == "test" {
				fmt.Println(text)
				//page.router.SetPage("wallet")
			} else {
				p.modalWalletPassword.StartWrongPassAnimation()
			}
		}
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func NewWalletCreateButton() *components.Button {
	addIcon, _ := widget.NewIcon(icons.ContentAddCircleOutline)

	var buttonStyle = components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "NEW WALLET",
		Icon:            addIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	}

	btn := components.NewButton(buttonStyle)
	btn.Label.Alignment = text.Middle
	btn.Style.Font.Weight = font.Bold
	return btn
}

type CreateWalletSelectionModal struct {
	modal     *components.Modal
	listStyle material.ListStyle
	items     []*CreateWalletListItem
}

func NewCreateWalletSelectionModal(th *material.Theme) *CreateWalletSelectionModal {
	modal := components.NewModal(th, layout.S, layout.UniformInset(25), components.NewModalAnimationDownUp())
	list := new(widget.List)
	list.Axis = layout.Vertical

	listStyle := material.List(th, list)

	fastIcon, _ := widget.NewIcon(icons.ImageFlashOn)
	newIcon, _ := widget.NewIcon(icons.ContentAddCircle)
	diskIcon, _ := widget.NewIcon(icons.FileFolder)
	seedIcon, _ := widget.NewIcon(icons.EditorShortText)

	items := []*CreateWalletListItem{
		NewCreateWalletListItem("Fast registration", fastIcon),
		NewCreateWalletListItem("Create new wallet", newIcon),
		NewCreateWalletListItem("Recover from Disk", diskIcon),
		NewCreateWalletListItem("Recover from Seed", seedIcon),
		NewCreateWalletListItem("Recover from Hex Seed", seedIcon),
	}

	return &CreateWalletSelectionModal{
		modal:     modal,
		listStyle: listStyle,
		items:     items,
	}
}

func (c *CreateWalletSelectionModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return c.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(10), Right: unit.Dp(0),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return c.listStyle.Layout(gtx, len(c.items), func(gtx layout.Context, index int) layout.Dimensions {
				if c.items[index].clickable.Clicked() {
					switch index {
					case 1:
						page_instance.router.SetCurrent("create_wallet_form")
					case 3:
						page_instance.router.SetCurrent("create_wallet_seed_form")
					case 4:
						page_instance.router.SetCurrent("create_wallet_hexseed_form")
					}

					c.modal.SetVisible(gtx, false)
				}

				return c.items[index].Layout(gtx, th)
			})
		})
	})
}

type CreateWalletListItem struct {
	text      string
	icon      *widget.Icon
	clickable *widget.Clickable
}

func NewCreateWalletListItem(text string, icon *widget.Icon) *CreateWalletListItem {
	return &CreateWalletListItem{
		text:      text,
		icon:      icon,
		clickable: new(widget.Clickable),
	}
}

func (c *CreateWalletListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	dims := c.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return c.icon.Layout(gtx, color.NRGBA{A: 255})
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					label := material.Label(th, unit.Sp(20), c.text)
					return label.Layout(gtx)
				}),
			)

		})
	})

	if c.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)

		paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(15),
			).Op(gtx.Ops),
		)
	}

	return dims
}

type WalletPasswordModal struct {
	editorStyle material.EditorStyle

	modal              *components.Modal
	animationWrongPass *animation.Animation

	iconLock *widget.Icon

	submitted  bool
	submitText string
}

func NewWalletPasswordModal(th *material.Theme) *WalletPasswordModal {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	editor.Mask = rune(42)
	editor.Focus()
	editorStyle := material.Editor(th, editor, "Enter password")
	editorStyle.TextSize = unit.Sp(20)

	animationWrongPass := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .05, ease.Linear),
		gween.New(1, -1, .05, ease.Linear),
		gween.New(-1, 0, .05, ease.Linear),
	))

	iconLock, _ := widget.NewIcon(icons.ActionLock)

	modal := components.NewModal(th, layout.Center, layout.UniformInset(25), components.NewModalAnimationScaleBounce())

	return &WalletPasswordModal{
		editorStyle:        editorStyle,
		modal:              modal,
		animationWrongPass: animationWrongPass,
		iconLock:           iconLock,
	}
}

func (w *WalletPasswordModal) Submit() (bool, string) {
	if w.submitted {
		w.submitted = false
		return true, w.submitText
	}

	return false, w.submitText
}

func (w *WalletPasswordModal) StartWrongPassAnimation() {
	w.animationWrongPass.Start()
}

func (w *WalletPasswordModal) Layout(gtx layout.Context) layout.Dimensions {
	for _, e := range w.editorStyle.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok {
			//w.animationWrongPass.Start()
			w.editorStyle.Editor.SetText("")
			w.submitText = e.Text
			w.submitted = true
		}
	}

	return w.modal.Layout(gtx,
		func(gtx layout.Context) {
			{
				state := w.animationWrongPass.Update(gtx)
				if state.Active {
					transform := f32.Affine2D{}.Offset(f32.Pt(state.Value*15, 0))
					op.Affine(transform).Add(gtx.Ops)
				}
			}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(25)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(25)
						gtx.Constraints.Max.Y = gtx.Dp(25)
						return w.iconLock.Layout(gtx, color.NRGBA{R: 0, G: 0, B: 0, A: 255})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
						return w.editorStyle.Layout(gtx)
					}),
				)
			})
		})
}

type WalletList struct {
	listStyle material.ListStyle
	items     []WalletListItem
}

func NewWalletList(th *material.Theme) *WalletList {
	list := new(widget.List)
	list.Axis = layout.Vertical

	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay
	listStyle.Indicator.MinorWidth = unit.Dp(10)
	listStyle.Indicator.CornerRadius = unit.Dp(5)
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	listStyle.Indicator.Color = black
	//listStyle.Indicator.HoverColor = f32color.Hovered(black)

	return &WalletList{
		listStyle: listStyle,
		items:     []WalletListItem{},
	}
}

func (l *WalletList) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		clip.UniformRRect(
			image.Rectangle{Max: gtx.Constraints.Max},
			gtx.Dp(unit.Dp(10)),
		).Op(gtx.Ops),
	)

	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return l.listStyle.Layout(gtx, len(l.items), func(gtx layout.Context, i int) layout.Dimensions {
			return l.items[i].Layout(gtx, th)
		})
	})

}

type WalletListItem struct {
	name      material.LabelStyle
	addr      material.LabelStyle
	Clickable *widget.Clickable

	rounded unit.Dp
}

func NewWalletListItem(th *material.Theme, name string, addr string) WalletListItem {
	namelbl := material.Label(th, unit.Sp(18), name)
	namelbl.Font.Weight = font.Bold
	addrlbl := material.Label(th, unit.Sp(15), addr)
	addrlbl.Color.A = 200

	return WalletListItem{
		name:      namelbl,
		addr:      addrlbl,
		Clickable: &widget.Clickable{},
		rounded:   unit.Dp(12),
	}
}

func (item *WalletListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		dims := layout.UniformInset(item.rounded).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Start}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(item.name.Layout),
						layout.Rigid(item.addr.Layout),
					)
				}),
			)
		})

		if item.Clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(item.rounded),
				).Op(gtx.Ops),
			)
		}

		return dims
	})
}
