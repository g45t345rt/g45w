package page_wallet

import (
	"encoding/json"
	"image"
	"strings"

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
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	gioui_hashicon "github.com/g45t345rt/gioui-hashicon"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageContacts struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	contactItems      []*ContactListItem
	txtFilterContacts *prefabs.Input
	buttonAddContact  *components.Button

	list              *widget.List
	buttonMenuContact *components.Button
	txtDims           layout.Dimensions
}

var _ router.Page = &PageContacts{}

func NewPageContacts() *PageContacts {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	menuIcon, _ := widget.NewIcon(icons.NavigationMenu)
	buttonMenuContact := components.NewButton(components.ButtonStyle{
		Icon:      menuIcon,
		Animation: components.NewButtonAnimationScale(.98),
	})

	txtFilterContacts := prefabs.NewInput()
	addContactIcon, _ := widget.NewIcon(icons.SocialPersonAdd)
	buttonAddContact := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      addContactIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonAddContact.Label.Alignment = text.Middle
	buttonAddContact.Style.Font.Weight = font.Bold

	return &PageContacts{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		list:              list,
		buttonMenuContact: buttonMenuContact,
		txtFilterContacts: txtFilterContacts,
		buttonAddContact:  buttonAddContact,
	}
}

func (p *PageContacts) IsActive() bool {
	return p.isActive
}

func (p *PageContacts) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Contacts") }
	page_instance.header.Subtitle = nil

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		p.buttonMenuContact.Style.Colors = theme.Current.ButtonIconPrimaryColors
		gtx.Constraints.Min.X = gtx.Dp(30)
		gtx.Constraints.Min.Y = gtx.Dp(30)

		if p.buttonMenuContact.Clicked(gtx) {
			go p.OpenMenu()
		}

		return p.buttonMenuContact.Layout(gtx, th)
	}

	if !page_instance.header.IsHistory(PAGE_CONTACTS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.Load()
}

func (p *PageContacts) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageContacts) Load() error {
	p.contactItems = make([]*ContactListItem, 0)

	wallet := wallet_manager.OpenedWallet
	contacts, err := wallet.GetContacts(wallet_manager.GetContactsParams{})
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		item := NewContactListItem(contact)
		p.contactItems = append(p.contactItems, item)
	}

	return nil
}

func (p *PageContacts) OpenMenu() {
	addContactIcon, _ := widget.NewIcon(icons.SocialPersonAdd)
	downIcon, _ := widget.NewIcon(icons.FileFileDownload)
	upIcon, _ := widget.NewIcon(icons.FileFileUpload)

	keyChan := listselect_modal.Instance.Open([]*listselect_modal.SelectListItem{
		listselect_modal.NewSelectListItem("add_contact",
			listselect_modal.NewItemText(addContactIcon, lang.Translate("Add contact")).Layout,
		),
		listselect_modal.NewSelectListItem("import_contacts",
			listselect_modal.NewItemText(downIcon, lang.Translate("Import contacts")).Layout,
		),
		listselect_modal.NewSelectListItem("export_contacts",
			listselect_modal.NewItemText(upIcon, lang.Translate("Export contacts")).Layout,
		),
	}, "")

	for key := range keyChan {
		switch key {
		case "add_contact":
			page_instance.pageContactForm.ClearForm()
			page_instance.pageRouter.SetCurrent(PAGE_CONTACT_FORM)
			page_instance.header.AddHistory(PAGE_CONTACT_FORM)
		case "export_contacts":
			exportContacts := func() error {
				file, err := app_instance.Explorer.CreateFile("contacts.json")
				if err != nil {
					return err
				}

				wallet := wallet_manager.OpenedWallet
				contacts, err := wallet.GetContacts(wallet_manager.GetContactsParams{})
				if err != nil {
					return err
				}

				data, err := json.MarshalIndent(contacts, "", " ")
				if err != nil {
					return err
				}

				_, err = file.Write(data)
				if err != nil {
					return err
				}
				defer file.Close()

				return nil
			}

			err := exportContacts()
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
					Text:       lang.Translate("Contacts exported."),
					CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
				})
			}
		case "import_contacts":
			importContacts := func() error {
				file, err := app_instance.Explorer.ChooseFile(".json")
				if err != nil {
					return err
				}

				reader := utils.ReadCloser{ReadCloser: file}
				data, err := reader.ReadAll()
				if err != nil {
					return err
				}

				var contacts []wallet_manager.Contact
				err = json.Unmarshal(data, &contacts)
				if err != nil {
					return err
				}

				wallet := wallet_manager.OpenedWallet
				for _, contact := range contacts {
					err = wallet.StoreContact(contact)
					if err != nil {
						return err
					}
				}

				return nil
			}

			err := importContacts()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			} else {
				p.Load()
				notification_modal.Open(notification_modal.Params{
					Type:       notification_modal.SUCCESS,
					Title:      lang.Translate("Success"),
					Text:       lang.Translate("Contacts imported."),
					CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
				})
			}
		}
	}
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

	if p.buttonAddContact.Clicked(gtx) {
		page_instance.pageContactForm.ClearForm()
		page_instance.pageRouter.SetCurrent(PAGE_CONTACT_FORM)
		page_instance.header.AddHistory(PAGE_CONTACT_FORM)
	}

	widgets := []layout.ListElement{}

	widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				p.txtFilterContacts.TextSize = unit.Sp(16)
				p.txtFilterContacts.Colors = theme.Current.InputColors
				// we store dims in a out of scope variable or it will be set to zero during the loop
				p.txtDims = p.txtFilterContacts.Layout(gtx, th, lang.Translate("Search contact..."))
				return p.txtDims
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				p.buttonAddContact.Style.Colors = theme.Current.ButtonPrimaryColors
				p.buttonAddContact.Flex = true
				size := image.Pt(p.txtDims.Size.Y, p.txtDims.Size.Y)
				gtx.Constraints.Min = size
				gtx.Constraints.Max = size
				return p.buttonAddContact.Layout(gtx, th)
			}),
		)
	})

	if len(p.contactItems) == 0 {
		widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
			txt := lang.Translate(`We see that you have no contacts. You can add a new contact here or import a list. You can find all these options in the menu button in the upper right.`)
			lbl := material.Label(th, unit.Sp(16), txt)
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})
	}

	for i := 0; i < len(p.contactItems); i++ {
		item := p.contactItems[i]
		searchText := p.txtFilterContacts.Editor.Text()
		display := true
		if searchText != "" {
			display = strings.Contains(item.contact.Name, searchText)
		}

		if display {
			widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
				return item.Layout(gtx, th)
			})
		}
	}

	widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Left: theme.PagePadding, Right: theme.PagePadding,
			Bottom: unit.Dp(10),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return widgets[index](gtx, index)
		})
	})
}

type ContactListItem struct {
	contact        wallet_manager.Contact
	buttonSelect   *components.Button
	buttonEdit     *components.Button
	listItemSelect *prefabs.ListItemSelect
	clickable      *widget.Clickable
}

func NewContactListItem(contact wallet_manager.Contact) *ContactListItem {
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Inset: layout.Inset{
			Top: unit.Dp(6), Bottom: unit.Dp(6),
			Left: unit.Dp(7), Right: unit.Dp(7),
		},
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	buttonEdit := components.NewButton(components.ButtonStyle{
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Inset: layout.Inset{
			Top: unit.Dp(6), Bottom: unit.Dp(6),
			Left: unit.Dp(7), Right: unit.Dp(7),
		},
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonEdit.Label.Alignment = text.Middle
	buttonEdit.Style.Font.Weight = font.Bold

	return &ContactListItem{
		contact:        contact,
		listItemSelect: prefabs.NewListItemSelect(),
		clickable:      new(widget.Clickable),
		buttonSelect:   buttonSelect,
		buttonEdit:     buttonEdit,
	}
}

func (item *ContactListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.buttonEdit.Clicked(gtx) {
		page_instance.pageContactForm.contact = &item.contact
		page_instance.pageRouter.SetCurrent(PAGE_CONTACT_FORM)
		page_instance.header.AddHistory(PAGE_CONTACT_FORM)
	}

	if item.buttonSelect.Clicked(gtx) {
		txtWalletAddr := page_instance.pageSendForm.walletAddrInput.txtWalletAddr
		txtWalletAddr.SetValue(item.contact.Addr)

		history := page_instance.header.History()
		page := history[len(history)-2]
		switch page {
		case PAGE_SEND_FORM:
			page_instance.header.GoBack()
		case PAGE_BALANCE_TOKENS:
			page_instance.pageSendForm.SetToken(wallet_manager.DeroToken())
			page_instance.pageRouter.SetCurrent(PAGE_SEND_FORM)
			page_instance.header.AddHistory(PAGE_SEND_FORM)
		}
	}

	if item.clickable.Clicked(gtx) {
		item.listItemSelect.Toggle()
	}

	r := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		dims := layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					size := float32(gtx.Dp(unit.Dp(30)))
					hashicon := gioui_hashicon.Hashicon{Config: gioui_hashicon.DefaultConfig}
					return hashicon.Layout(gtx, size, item.contact.Addr)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
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
							label.Color = theme.Current.TextMuteColor
							return label.Layout(gtx)
						}),
					)
				}),
			)
		})
		c := r.Stop()

		if item.clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(10),
				).Op(gtx.Ops),
			)
		}

		layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			item.buttonSelect.Text = lang.Translate("Select")
			item.buttonSelect.Style.Colors = theme.Current.ButtonPrimaryColors
			item.buttonEdit.Text = lang.Translate("Edit")
			item.buttonEdit.Style.Colors = theme.Current.ButtonPrimaryColors
			return item.listItemSelect.Layout(gtx, th, []*components.Button{item.buttonSelect, item.buttonEdit})
		})

		c.Add(gtx.Ops)
		return dims
	})
	c := r.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
		clip.RRect{
			Rect: image.Rectangle{Max: dims.Size},
			SE:   gtx.Dp(10), SW: gtx.Dp(10),
			NW: gtx.Dp(10), NE: gtx.Dp(10),
		}.Op(gtx.Ops))

	c.Add(gtx.Ops)

	return dims
}
