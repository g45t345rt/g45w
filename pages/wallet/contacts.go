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
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/contact_manager"
	"github.com/g45t345rt/g45w/prefabs"
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

	p.contactItems = make([]*ContactListItem, 0)
	for _, contact := range page_instance.contactManager.Contacts {
		item := NewContactListItem(contact)
		p.contactItems = append(p.contactItems, item)
	}

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

	if p.buttonAddContact.Clickable.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_CONTACT_FORM)
	}

	widgets := []layout.Widget{}

	if len(p.contactItems) == 0 {
		return layout.Inset{
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), "You didn't add any contacts yet.")
			return lbl.Layout(gtx)
		})
	}

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
	contact        contact_manager.Contact
	listItemSelect *prefabs.ListItemSelectEdit
	clickable      *widget.Clickable
}

func NewContactListItem(contact contact_manager.Contact) *ContactListItem {
	return &ContactListItem{
		contact:        contact,
		listItemSelect: prefabs.NewListItemSelectEdit(),
		clickable:      new(widget.Clickable),
	}
}

func (item *ContactListItem) Layout(gtx layout.Context) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.listItemSelect.EditClicked() {
		page_instance.pageContactForm.contact = &item.contact
		page_instance.pageRouter.SetCurrent(PAGE_CONTACT_FORM)
	}

	if item.listItemSelect.SelectClicked() {
		page_instance.pageSendForm.txtWalletAddr.SetValue(item.contact.Addr)
		page_instance.pageRouter.SetCurrent(PAGE_SEND_FORM)
	}

	return layout.Inset{
		Top: unit.Dp(0), Bottom: unit.Dp(10),
		Right: unit.Dp(30), Left: unit.Dp(30),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		th := app_instance.Theme
		m := op.Record(gtx.Ops)
		dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			dims := layout.Inset{
				Top: unit.Dp(10), Bottom: unit.Dp(10),
				Left: unit.Dp(15), Right: unit.Dp(15),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				dims := layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								label := material.Label(th, unit.Sp(20), item.contact.Name)
								label.Font.Weight = font.Bold
								return label.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								addr := utils.ReduceAddr(item.contact.Addr)
								label := material.Label(th, unit.Sp(16), addr)
								label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 150}
								return label.Layout(gtx)
							}),
						)
					}),
				)

				item.listItemSelect.Layout(gtx, th)
				return dims
			})

			buttonEditHovered := item.listItemSelect.ButtonEdit.Clickable.Hovered()
			buttonSelectHovered := item.listItemSelect.ButtonSelect.Clickable.Hovered()
			if item.clickable.Hovered() && !buttonEditHovered && !buttonSelectHovered {
				pointer.CursorPointer.Add(gtx.Ops)
				paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
					clip.UniformRRect(
						image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
						gtx.Dp(10),
					).Op(gtx.Ops),
				)
			}

			if item.clickable.Clicked() && !buttonEditHovered && !buttonSelectHovered {
				item.listItemSelect.Toggle()
			}

			return dims
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
}
