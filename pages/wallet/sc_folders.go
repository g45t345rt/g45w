package page_wallet

import (
	"image"
	"image/color"

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
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCFolders struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	tokenItems []*TokenFolderItem

	list                 *widget.List
	buttonAddFolderToken *components.Button
}

var _ router.Page = &PageSCFolders{}

func NewPageSCFolders() *PageSCFolders {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	folderIcon, _ := widget.NewIcon(icons.FileCreateNewFolder)
	buttonAddFolderToken := components.NewButton(components.ButtonStyle{
		Icon:      folderIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation: components.NewButtonAnimationScale(.98),
	})

	return &PageSCFolders{
		animationEnter:       animationEnter,
		animationLeave:       animationLeave,
		list:                 list,
		buttonAddFolderToken: buttonAddFolderToken,
	}
}

func (p *PageSCFolders) IsActive() bool {
	return p.isActive
}

func (p *PageSCFolders) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Tokens"))
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = p.buttonAddFolderToken

	p.tokenItems = make([]*TokenFolderItem, 0)

	for i := 0; i < 47; i++ {
		p.tokenItems = append(p.tokenItems, NewTokenFolderItem())
	}

	if !page_instance.header.IsHistory(PAGE_SC_FOLDERS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageSCFolders) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSCFolders) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonAddFolderToken.Clickable.Clicked() {

	}

	widgets := []layout.ListElement{}

	if len(p.tokenItems) == 0 {
		return layout.Inset{
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("You didn't add any tokens yet."))
			return lbl.Layout(gtx)
		})
	}

	widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			//layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), lang.Translate("100 tokens - 15 folders"))
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		)
	})

	var itemIndex = 1
	for i := 0; i < len(p.tokenItems); i += 3 {
		widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
			dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return p.tokenItems[itemIndex].Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if itemIndex+1 < len(p.tokenItems) {
						return p.tokenItems[itemIndex+1].Layout(gtx, th)
					}
					return layout.Dimensions{}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if itemIndex+2 < len(p.tokenItems) {
						return p.tokenItems[itemIndex+2].Layout(gtx, th)
					}
					return layout.Dimensions{}
				}),
			)

			itemIndex = itemIndex + 3
			return dims
		})
	}

	widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{Left: unit.Dp(30), Right: unit.Dp(30)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return widgets[index](gtx, index)
		})
	})
}

type TokenFolderItem struct {
	folderIcon *widget.Icon
	Clickable  *widget.Clickable
}

func NewTokenFolderItem() *TokenFolderItem {
	folderIcon, _ := widget.NewIcon(icons.FileFolder)
	return &TokenFolderItem{
		folderIcon: folderIcon,
		Clickable:  new(widget.Clickable),
	}
}

func (item *TokenFolderItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.Clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Constraints.Max.X
				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.UniformRRect(image.Rectangle{
					Max: gtx.Constraints.Max,
				}, gtx.Dp(10)).Op(gtx.Ops))

				return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return item.folderIcon.Layout(gtx, color.NRGBA{A: 255})
				})
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), "wergasdfasdfasdfasdfasdfasdfasdf")
			lbl.Alignment = text.Middle
			lbl.Font.Weight = font.Bold
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(12), "10 tokens")
			lbl.Alignment = text.Middle
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
	)
}

type TokenMenuSelect struct {
	SelectModal *prefabs.SelectModal
	items       *prefabs.SelectListItem
}

func NewTokenMenuSelect() *TokenMenuSelect {
	//addTokenIcon, _ := widget.NewIcon(icons.ContentAdd)
	//addFolderIcon, _ := widget.NewIcon(icons.FileCreateNewFolder)

	items := []*prefabs.SelectListItem{}

	selectModal := prefabs.NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			selectModal.Layout(gtx, th, items)
		},
	})

	return &TokenMenuSelect{
		SelectModal: selectModal,
	}
}
