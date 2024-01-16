package page_settings

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
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

	gateway app_db.IPFSGateway

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

	return &PageEditIPFSGateway{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		buttonEdit:   buttonEdit,
		buttonDelete: buttonDelete,
		txtName:      txtName,
		txtEndpoint:  txtEndpoint,
		switchActive: new(widget.Bool),

		list: list,
	}
}

func (p *PageEditIPFSGateway) IsActive() bool {
	return p.isActive
}

func (p *PageEditIPFSGateway) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Edit IPFS Gateway") }
	page_instance.header.RightLayout = nil

	if !page_instance.header.IsHistory(PAGE_EDIT_IPFS_GATEWAY) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.txtEndpoint.SetValue(p.gateway.Endpoint)
	p.txtName.SetValue(p.gateway.Name)
	p.switchActive.Value = p.gateway.Active
}

func (p *PageEditIPFSGateway) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
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

	if p.buttonEdit.Clicked(gtx) {
		p.submitForm(gtx)
	}

	if p.buttonDelete.Clicked(gtx) {
		go func() {
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			if <-yesChan {
				err := p.removeGateway()
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
						Text:       lang.Translate("Gateway deleted."),
						CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
					})
					page_instance.header.GoBack()
				}
			}
		}()
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
					s := material.Switch(th, p.switchActive, "")
					s.Color = theme.Current.SwitchColors
					return s.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), lang.Translate("Inactive gateway will not be used when fetching IPFS content."))
					lbl.Color = theme.Current.TextMuteColor
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

	if p.txtName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked(gtx) {
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
	err := app_db.DelIPFSGateway(0)
	if err != nil {
		return err
	}

	if node_manager.CurrentNode != nil {
		if node_manager.CurrentNode.Endpoint == endpoint {
			node_manager.Set(nil, true)
			/*node_manager.CurrentNode = nil
			walletapi.Connected = false

			settings.App.NodeEndpoint = ""
			err := settings.Save()
			if err != nil {
				return err
			}*/
		}
	}

	return nil
}

func (p *PageEditIPFSGateway) submitForm(gtx layout.Context) {
	p.buttonEdit.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonEdit.SetLoading(false)
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
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

		gateway := app_db.IPFSGateway{
			ID:          p.gateway.ID,
			Name:        txtName.Text(),
			Endpoint:    txtEndpoint.Text(),
			Active:      p.switchActive.Value,
			OrderNumber: p.gateway.OrderNumber,
		}

		err := gateway.TestFetch()
		if err != nil {
			setError(err)
			return
		}

		err = app_db.UpdateIPFSGateway(gateway)
		if err != nil {
			setError(err)
			return
		}

		p.buttonEdit.SetLoading(false)
		notification_modal.Open(notification_modal.Params{
			Type:       notification_modal.SUCCESS,
			Title:      lang.Translate("Success"),
			Text:       lang.Translate("Data saved."),
			CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
		})
		page_instance.header.GoBack()
	}()
}
