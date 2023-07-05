package page_wallet

import (
	"errors"
	"image"
	"image/color"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/contact_manager"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageContactForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonCreate  *components.Button
	buttonDelete  *components.Button
	txtName       *components.TextField
	txtAddr       *components.TextField
	txtNote       *components.TextField
	confirmDelete *components.Confirm

	contact *contact_manager.Contact

	list *widget.List
}

var _ router.Page = &PageContactForm{}

func NewPageContactForm() *PageContactForm {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	buttonCreate := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            saveIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonCreate.Label.Alignment = text.Middle
	buttonCreate.Style.Font.Weight = font.Bold

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDelete := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            deleteIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 255, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonDelete.Label.Alignment = text.Middle
	buttonDelete.Style.Font.Weight = font.Bold

	txtName := components.NewTextField()
	txtAddr := components.NewTextField()
	txtNote := components.NewTextField()
	txtNote.Editor().SingleLine = false
	txtNote.Editor().Submit = false

	confirmDelete := components.NewConfirm(layout.Center)
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			confirmDelete.Prompt = lang.Translate("Are you sure?")
			confirmDelete.NoText = lang.Translate("NO")
			confirmDelete.YesText = lang.Translate("YES")
			confirmDelete.Layout(gtx, th)
		},
	})

	return &PageContactForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonCreate:  buttonCreate,
		buttonDelete:  buttonDelete,
		txtName:       txtName,
		txtAddr:       txtAddr,
		txtNote:       txtNote,
		confirmDelete: confirmDelete,

		list: list,
	}
}

func (p *PageContactForm) IsActive() bool {
	return p.isActive
}

func (p *PageContactForm) Enter() {
	p.isActive = true

	if p.contact != nil {
		page_instance.header.SetTitle(lang.Translate("Edit Contact"))
		p.txtName.SetValue(p.contact.Name)
		p.txtAddr.SetValue(p.contact.Addr)
		p.txtNote.SetValue(p.contact.Note)
	} else {
		page_instance.header.SetTitle(lang.Translate("New Contact"))
	}

	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil

	if !page_instance.header.IsHistory(PAGE_CONTACT_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageContactForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
	p.contact = nil
}

func (p *PageContactForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonCreate.Clickable.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("New contact added"))
			notification_modals.SuccessInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
			page_instance.pageRouter.SetCurrent(PAGE_CONTACTS)
			p.clearForm()
		}
	}

	if p.buttonDelete.Clickable.Clicked() {
		p.confirmDelete.SetVisible(gtx, true)
	}

	if p.confirmDelete.ClickedYes() {
		err := page_instance.contactManager.DelContact(p.contact.Addr)
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Contact deleted"))
			notification_modals.SuccessInstance.SetVisible(gtx, true, notification_modals.CLOSE_AFTER_DEFAULT)
			page_instance.pageRouter.SetCurrent(PAGE_CONTACTS)
			p.clearForm()
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtAddr.Layout(gtx, th, lang.Translate("Address"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.txtNote.Input.EditorMinY = gtx.Dp(75)
			return p.txtNote.Layout(gtx, th, lang.Translate("Note"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonCreate.Text = lang.Translate("ADD CONTACT")
			return p.buttonCreate.Layout(gtx, th)
		},
	}

	if p.contact != nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			max := image.Pt(gtx.Dp(unit.Dp(gtx.Constraints.Max.X)), 5)
			paint.FillShape(gtx.Ops, color.NRGBA{A: 150}, clip.Rect{
				Min: image.Pt(0, 0),
				Max: max,
			}.Op())
			return layout.Dimensions{Size: max}
		})

		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			p.buttonDelete.Text = lang.Translate("DELETE CONTACT")
			return p.buttonDelete.Layout(gtx, th)
		})
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtAddr.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	if p.txtNote.Input.Clickable.Clicked() {
		p.list.ScrollTo(2)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageContactForm) clearForm() {
	p.txtName.Editor().SetText("")
	p.txtAddr.Editor().SetText("")
	p.txtNote.Editor().SetText("")
}

func (p *PageContactForm) submitForm() error {
	txtName := p.txtName.Editor()
	txtAddr := p.txtAddr.Editor()
	txtNote := p.txtNote.Editor()

	_, err := globals.ParseValidateAddress(txtAddr.Text())
	if err != nil {
		return errors.New("invalid address")
	}

	err = page_instance.contactManager.SetContact(contact_manager.Contact{
		Name:      txtName.Text(),
		Addr:      txtAddr.Text(),
		Note:      txtNote.Text(),
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	p.clearForm()
	return nil
}
