package page_wallet

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/g45t345rt/g45w/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageContacts struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	contactItems []*ContactListItem

	listStyle        material.ListStyle
	buttonAddContact *components.Button
}

var _ router.Page = &PageContacts{}

func NewPageContacts() *PageContacts {
	th := app_instance.Theme

	contactItems := []*ContactListItem{}

	for i := 0; i < 10; i++ {
		contactItems = append(contactItems, &ContactListItem{
			name:      fmt.Sprintf("asd %d", i),
			addr:      "dero1qy4egwhfjxtt96pdlkkc3d99yqkewkvljf50unqh5vz7xepl8w6yyqgklpjsz",
			Clickable: new(widget.Clickable),
		})
	}

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	addIcon, _ := widget.NewIcon(icons.SocialPersonAdd)
	buttonAddContact := components.NewButton(components.ButtonStyle{
		Icon:      addIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation: components.NewButtonAnimationScale(.98),
	})

	return &PageContacts{
		contactItems:     contactItems,
		animationEnter:   animationEnter,
		animationLeave:   animationLeave,
		listStyle:        listStyle,
		buttonAddContact: buttonAddContact,
	}
}

func (p *PageContacts) IsActive() bool {
	return p.isActive
}

func (p *PageContacts) Enter() {
	p.isActive = true
	page_instance.header.SetTitle("Contacts")
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = p.buttonAddContact

	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageContacts) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageContacts) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	//if len(p.contactItems) == 0 {
	return layout.Inset{
		Left: unit.Dp(30), Right: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), "You didn't add any contacts yet.")
		return lbl.Layout(gtx)
	})
	//}

	for _, item := range p.contactItems {
		widgets = append(widgets, item.Layout)
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
	})

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

type ContactListItem struct {
	name string
	addr string

	Clickable *widget.Clickable
}

func (item *ContactListItem) Layout(gtx layout.Context) layout.Dimensions {
	if item.Clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(10),
		Right: unit.Dp(30), Left: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		th := app_instance.Theme
		m := op.Record(gtx.Ops)
		dims := item.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(18), item.name)
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								addr := utils.ReduceAddr(item.addr)
								label := material.Label(th, unit.Sp(14), addr)
								label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
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

		return dims
	})
	/*
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return dims
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
		)


			if item.Clickable.Hovered() {
				pointer.CursorPointer.Add(gtx.Ops)
				bounds := image.Rect(0, 0, dims.Size.X, dims.Size.Y)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
					clip.UniformRRect(bounds, 10).Op(gtx.Ops),
				)
			}*/
}
