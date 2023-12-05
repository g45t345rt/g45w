package page_wallet

import (
	"fmt"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageContactForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonSave   *components.Button
	buttonDelete *components.Button
	txtName      *prefabs.TextField
	txtAddr      *prefabs.TextField

	contact *wallet_manager.Contact

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
	buttonSave := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      saveIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonSave.Label.Alignment = text.Middle
	buttonSave.Style.Font.Weight = font.Bold

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonDelete := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      deleteIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonDelete.Label.Alignment = text.Middle
	buttonDelete.Style.Font.Weight = font.Bold

	txtName := prefabs.NewTextField()
	txtAddr := prefabs.NewTextField()

	return &PageContactForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonSave:   buttonSave,
		buttonDelete: buttonDelete,
		txtName:      txtName,
		txtAddr:      txtAddr,

		list: list,
	}
}

func (p *PageContactForm) IsActive() bool {
	return p.isActive
}

func (p *PageContactForm) Enter() {
	p.isActive = true

	if p.contact != nil {
		page_instance.header.Title = func() string { return lang.Translate("Edit Contact") }
		p.txtName.SetValue(p.contact.Name)
		p.txtAddr.SetValue(p.contact.Addr)
	} else {
		page_instance.header.Title = func() string { return lang.Translate("New Contact") }
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

	if p.buttonSave.Clicked() {
		err := p.submitForm()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("New contact added"))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			page_instance.pageContacts.Load()
			page_instance.header.GoBack()
			p.ClearForm()
		}
	}

	if p.buttonDelete.Clicked() {
		go func() {
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			for yes := range yesChan {
				if yes {
					wallet := wallet_manager.OpenedWallet
					err := wallet.DelContact(p.contact.Addr)
					if err != nil {
						notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
						notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
					} else {
						notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Contact deleted"))
						notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
						page_instance.pageContacts.Load()
						page_instance.header.GoBack()
						p.ClearForm()
					}
				}
			}
		}()
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Contact Name"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.txtAddr.Layout(gtx, th, lang.Translate("DERO Address / Name"), "")
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					addr, err := globals.ParseValidateAddress(p.txtAddr.Value())
					if err == nil {
						status := lang.Translate("This is a valid address.")
						if addr.IsIntegratedAddress() {
							status = lang.Translate("This is an integrated address.")
						}

						return layout.Inset{Top: unit.Dp(3)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(14), status)
							lbl.Color = theme.Current.TextMuteColor
							return lbl.Layout(gtx)
						})
					}

					return layout.Dimensions{}
				}),
			)
		},

		func(gtx layout.Context) layout.Dimensions {
			if p.contact != nil {
				p.buttonSave.Text = lang.Translate("SAVE")
			} else {
				p.buttonSave.Text = lang.Translate("ADD")
			}

			p.buttonSave.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonSave.Layout(gtx, th)
		},
	}

	if p.contact != nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, 5)
		})

		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			p.buttonDelete.Text = lang.Translate("DELETE CONTACT")
			p.buttonDelete.Style.Colors = theme.Current.ButtonDangerColors
			return p.buttonDelete.Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtAddr.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageContactForm) ClearForm() {
	p.contact = nil
	p.txtName.Editor().SetText("")
	p.txtAddr.Editor().SetText("")
}

func (p *PageContactForm) submitForm() error {
	txtName := p.txtName.Editor()
	txtAddr := p.txtAddr.Editor()

	// don't validate addr because of service names
	//_, err := globals.ParseValidateAddress(txtAddr.Text())
	//	if err != nil {
	//	return errors.New("invalid address")
	//}

	if txtName.Text() == "" {
		return fmt.Errorf("name cannot be empty")
	}

	wallet := wallet_manager.OpenedWallet

	err := wallet.StoreContact(wallet_manager.Contact{
		Name:      txtName.Text(),
		Addr:      txtAddr.Text(),
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	p.ClearForm()
	return nil
}
