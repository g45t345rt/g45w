package page_wallet

import (
	"fmt"
	"image"
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
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/containers/qrcode_scan_modal"
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

	buttonSave     *components.Button
	buttonDelete   *components.Button
	txtName        *prefabs.TextField
	txtAddr        *prefabs.Input
	txtNote        *prefabs.TextField
	buttonScanAddr *components.Button
	txtDims        layout.Dimensions

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

	scanIcon, _ := widget.NewIcon(app_icons.QRCodeScanner)
	buttonScanAddr := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      scanIcon,
		Animation: components.NewButtonAnimationDefault(),
	})

	txtName := prefabs.NewTextField()
	txtAddr := prefabs.NewInput()
	txtNote := prefabs.NewTextField()
	txtNote.Editor().SingleLine = false
	txtNote.Editor().Submit = false

	return &PageContactForm{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonSave:     buttonSave,
		buttonDelete:   buttonDelete,
		buttonScanAddr: buttonScanAddr,
		txtName:        txtName,
		txtAddr:        txtAddr,
		txtNote:        txtNote,

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
		p.txtNote.SetValue(p.contact.Note)
	} else {
		page_instance.header.Title = func() string { return lang.Translate("New Contact") }
	}

	page_instance.header.Subtitle = nil
	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = nil

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

	if p.buttonSave.Clicked(gtx) {
		err := p.submitForm()
		if err != nil {
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
		} else {
			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.SUCCESS,
				Title:      lang.Translate("Success"),
				Text:       lang.Translate("New contact added."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
			page_instance.pageContacts.Load()
			page_instance.header.GoBack()
			p.ClearForm()
		}
	}

	if p.buttonDelete.Clicked(gtx) {
		go func() {
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			if <-yesChan {
				wallet := wallet_manager.OpenedWallet
				err := wallet.DelContact(p.contact.Addr)
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				} else {
					notification_modal.Open(notification_modal.Params{
						Type:       notification_modal.SUCCESS,
						Title:      lang.Translate("Success"),
						Text:       lang.Translate("Contact deleted."),
						CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
					})
					page_instance.pageContacts.Load()
					page_instance.header.GoBack()
					p.ClearForm()
				}
			}
		}()
	}

	if p.buttonScanAddr.Clicked(gtx) {
		go qrcode_scan_modal.Instance.Open()
	}

	{
		sent, value := qrcode_scan_modal.Instance.Value()
		if sent {
			p.txtAddr.SetValue(value)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate("Address"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							p.txtAddr.Colors = theme.Current.InputColors
							p.txtDims = p.txtAddr.Layout(gtx, th, "")
							return p.txtDims
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							p.buttonScanAddr.Style.Colors = theme.Current.ButtonPrimaryColors
							p.buttonScanAddr.Flex = true
							size := image.Pt(p.txtDims.Size.Y, p.txtDims.Size.Y)
							gtx.Constraints.Min = size
							gtx.Constraints.Max = size
							return p.buttonScanAddr.Layout(gtx, th)
						}),
					)
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
			p.txtNote.Input.EditorMinY = gtx.Dp(75)
			return p.txtNote.Layout(gtx, th, lang.Translate("Note"), "")
		},
		func(gtx layout.Context) layout.Dimensions {
			if p.contact != nil {
				p.buttonSave.Text = lang.Translate("SAVE CONTACT")
			} else {
				p.buttonSave.Text = lang.Translate("ADD CONTACT")
			}

			p.buttonSave.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonSave.Layout(gtx, th)
		},
	}

	if p.contact != nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, unit.Dp(5))
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

	if p.txtName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtAddr.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(1)
	}

	if p.txtNote.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(2)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageContactForm) ClearForm() {
	p.contact = nil
	p.txtName.Editor().SetText("")
	p.txtAddr.Editor.SetText("")
	p.txtNote.Editor().SetText("")
}

func (p *PageContactForm) submitForm() error {
	txtName := p.txtName.Editor()
	txtAddr := p.txtAddr.Editor
	txtNote := p.txtNote.Editor()

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
		Note:      txtNote.Text(),
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	p.ClearForm()
	return nil
}
