package page_settings

import (
	"fmt"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageEditIPFSGateway struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonEdit   *components.Button
	buttonDelete *components.Button
	txtEndpoint  *prefabs.TextField
	txtName      *prefabs.TextField
	switchActive *widget.Bool

	gateway app_data.IPFSGateway

	confirmDelete *prefabs.Confirm

	list *widget.List
}

var _ router.Page = &PageEditIPFSGateway{}

func NewPageEditIPFSGateway() *PageEditIPFSGateway {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	saveIcon, _ := widget.NewIcon(icons.ContentSave)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonEdit := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        saveIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonEdit.Label.Alignment = text.Middle
	buttonEdit.Style.Font.Weight = font.Bold

	txtName := prefabs.NewTextField()
	txtEndpoint := prefabs.NewTextField()

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

	confirmDelete := prefabs.NewConfirm(layout.Center)
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			confirmDelete.Layout(gtx, th, prefabs.ConfirmText{
				Prompt: lang.Translate("Are you sure?"),
				No:     lang.Translate("NO"),
				Yes:    lang.Translate("YES"),
			})
		},
	})

	return &PageEditIPFSGateway{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonEdit:   buttonEdit,
		buttonDelete: buttonDelete,
		txtName:      txtName,
		txtEndpoint:  txtEndpoint,
		switchActive: new(widget.Bool),

		confirmDelete: confirmDelete,

		list: list,
	}
}

func (p *PageEditIPFSGateway) IsActive() bool {
	return p.isActive
}

func (p *PageEditIPFSGateway) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Edit IPFS Gateway") }
	page_instance.header.ButtonRight = nil
	p.animationEnter.Start()
	p.animationLeave.Reset()

	p.txtEndpoint.SetValue(p.gateway.Endpoint)
	p.txtName.SetValue(p.gateway.Name)
	p.switchActive.Value = p.gateway.Active
}

func (p *PageEditIPFSGateway) Leave() {
	if page_instance.header.IsHistory(PAGE_EDIT_IPFS_GATEWAY) {
		p.animationEnter.Reset()
		p.animationLeave.Start()
	} else {
		p.isActive = false
	}
}

func (p *PageEditIPFSGateway) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonEdit.Clicked() {
		p.submitForm(gtx)
	}

	if p.buttonDelete.Clicked() {
		p.confirmDelete.SetVisible(true)
	}

	if p.confirmDelete.ClickedYes() {
		err := p.removeGateway()
		if err != nil {
			notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {

			notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Gateway deleted"))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			page_instance.header.GoBack()
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "deronfts.com")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th, lang.Translate("Endpoint"), "https://ipfs.deronfts.com/ipfs")
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(20), lang.Translate("Active"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					s := material.Switch(th, p.switchActive, lang.Translate("Set Active Gateway"))
					s.Color.Enabled = color.NRGBA{A: 255}
					return s.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Inactive gateway will not be used when fetching IPFS content."))
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonEdit.Text = lang.Translate("SAVE")
			p.buttonEdit.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonEdit.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, 5)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonDelete.Text = lang.Translate("DELETE GATEWAY")
			p.buttonDelete.Style.Colors = theme.Current.ButtonDangerColors
			return p.buttonDelete.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked() {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked() {
		p.list.ScrollTo(1)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageEditIPFSGateway) removeGateway() error {
	endpoint := p.gateway.Endpoint
	err := app_data.DelIPFSGateway(0)
	if err != nil {
		return err
	}

	if node_manager.CurrentNode != nil {
		if node_manager.CurrentNode.Endpoint == endpoint {
			node_manager.CurrentNode = nil
			walletapi.Connected = false

			settings.App.NodeEndpoint = ""
			err := settings.Save()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PageEditIPFSGateway) submitForm(gtx layout.Context) {
	p.buttonEdit.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonEdit.SetLoading(false)
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}

		txtName := p.txtName.Editor()
		txtEndpoint := p.txtEndpoint.Editor()

		if txtName.Text() == "" {
			setError(fmt.Errorf("enter name"))
			return
		}

		if txtEndpoint.Text() == "" {
			setError(fmt.Errorf("enter endpoint"))
			return
		}

		endpoint := txtEndpoint.Text()
		gateway := app_data.IPFSGateway{
			ID:       p.gateway.ID,
			Name:     txtName.Text(),
			Endpoint: endpoint,
			Active:   p.switchActive.Value,
		}

		err := gateway.TestFetch()
		if err != nil {
			setError(err)
			return
		}

		err = app_data.UpdateIPFSGateway(gateway)
		if err != nil {
			setError(err)
			return
		}

		p.buttonEdit.SetLoading(false)
		notification_modals.SuccessInstance.SetText("Success", lang.Translate("Data saved"))
		notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		page_instance.header.GoBack()
	}()
}
