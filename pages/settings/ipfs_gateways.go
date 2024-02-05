package page_settings

import (
	"image"

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
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageIPFSGateways struct {
	isActive            bool
	list                *widget.List
	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonInfo             *components.Button
	modalInfo              *components.Modal
	buttonResetGatewayList *components.Button

	gatewayList *GatewayList

	buttonAdd *components.Button
}

var _ router.Page = &PageIPFSGateways{}

func NewPageIPFSGateways() *PageIPFSGateways {
	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Icon: addIcon,
	})

	gatewayList := NewGatewayList()

	infoIcon, _ := widget.NewIcon(icons.ActionInfo)
	buttonInfo := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      infoIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})

	modalInfo := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  true,
		Direction:           layout.Center,
		Inset:               layout.UniformInset(unit.Dp(30)),
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Animation:           components.NewModalAnimationScaleBounce(),
	})

	resetIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonResetGatewayList := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(16),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Icon:      resetIcon,
		IconGap:   unit.Dp(10),
	})
	buttonResetGatewayList.Label.Alignment = text.Middle
	buttonResetGatewayList.Style.Font.Weight = font.Bold

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			modalInfo.Style.Colors = theme.Current.ModalColors
			modalInfo.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(20), lang.Translate("Why use IPFS?"))
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("Storing data on Dero can be expensive. An alternative approach is to utilize IPFS for storing images, files, and other content, while saving only the corresponding links. This section let you add multiple IPFS gateways, ensuring seamless access to IPFS content within Dero smart contracts."))
							return lbl.Layout(gtx)
						}),
					)
				})
			})
		},
	})

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_IPFS_GATEWAYS)

	return &PageIPFSGateways{
		list:                   list,
		headerPageAnimation:    headerPageAnimation,
		gatewayList:            gatewayList,
		buttonAdd:              buttonAdd,
		buttonInfo:             buttonInfo,
		modalInfo:              modalInfo,
		buttonResetGatewayList: buttonResetGatewayList,
	}
}

func (p *PageIPFSGateways) IsActive() bool {
	return p.isActive
}

func (p *PageIPFSGateways) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("IPFS Gateways") }
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("Interplanetary File System"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}
	page_instance.header.RightLayout = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		p.buttonAdd.Style.Colors = theme.Current.ButtonIconPrimaryColors
		gtx.Constraints.Min.X = gtx.Dp(30)
		gtx.Constraints.Min.Y = gtx.Dp(30)

		if p.buttonAdd.Clicked(gtx) {
			page_instance.pageRouter.SetCurrent(PAGE_ADD_IPFS_GATEWAY)
			page_instance.header.AddHistory(PAGE_ADD_IPFS_GATEWAY)
		}

		return p.buttonAdd.Layout(gtx, th)
	}

	p.gatewayList.Load()
}

func (p *PageIPFSGateways) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageIPFSGateways) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonInfo.Clicked(gtx) {
		p.modalInfo.SetVisible(true)
	}

	if p.buttonResetGatewayList.Clicked(gtx) {
		go func() {
			yes := <-confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			if yes {
				err := app_db.ResetIPFSGateways()
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				} else {
					p.gatewayList.Load()
					notification_modal.Open(notification_modal.Params{
						Type:       notification_modal.SUCCESS,
						Title:      lang.Translate("Success"),
						Text:       lang.Translate("List reset to default."),
						CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
					})
				}
				app_instance.Window.Invalidate()
			}
		}()
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonInfo.Text = lang.Translate("Why use IPFS?")
		p.buttonInfo.Style.Colors = theme.Current.ButtonPrimaryColors
		return p.buttonInfo.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.gatewayList.Layout(gtx, th, lang.Translate("You don't have any IPFS gateways available."))
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonResetGatewayList.Text = lang.Translate("Reset gateway list")
		p.buttonResetGatewayList.Style.Colors = theme.Current.ButtonPrimaryColors
		return p.buttonResetGatewayList.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

type GatewayList struct {
	items     []GatewayListItem
	list      *widget.List
	dragItems *components.DragItems
}

func NewGatewayList() *GatewayList {
	list := new(widget.List)
	list.Axis = layout.Vertical

	return &GatewayList{
		list:      list,
		items:     make([]GatewayListItem, 0),
		dragItems: components.NewDragItems(),
	}
}

func (l *GatewayList) Load() error {
	items := make([]GatewayListItem, 0)

	gateways, err := app_db.GetIPFSGateways(app_db.GetIPFSGatewaysParams{})
	if err != nil {
		return err
	}

	for _, gateway := range gateways {
		items = append(items, NewGatewayListItem(gateway))
	}

	l.items = items
	return nil
}

func (l *GatewayList) Layout(gtx layout.Context, th *material.Theme, emptyText string) layout.Dimensions {
	{
		moved, cIndex, nIndex := l.dragItems.ItemMoved()
		if moved {
			go func() {
				updateIndex := func() error {
					gateway := l.items[cIndex].gateway
					gateway.OrderNumber = nIndex
					err := app_db.UpdateIPFSGateway(gateway)
					if err != nil {
						return err
					}

					return l.Load()
				}

				err := updateIndex()
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				}
				app_instance.Window.Invalidate()
			}()
		}
	}

	r := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if len(l.items) == 0 {
			lbl := material.Label(th, unit.Sp(16), emptyText)
			return lbl.Layout(gtx)
		} else {
			listStyle := material.List(th, l.list)
			listStyle.AnchorStrategy = material.Overlay
			listStyle.Indicator.MinorWidth = unit.Dp(10)
			listStyle.Indicator.CornerRadius = unit.Dp(5)
			listStyle.Indicator.Color = theme.Current.ListScrollBarBgColor

			return l.dragItems.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
				return listStyle.Layout(gtx, len(l.items), func(gtx layout.Context, index int) layout.Dimensions {
					l.dragItems.LayoutItem(gtx, index, func(gtx layout.Context) layout.Dimensions {
						return l.items[index].Layout(gtx, th, true)
					})

					return l.items[index].Layout(gtx, th, false)
				})
			})
		}
	})
	c := r.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
		clip.UniformRRect(
			image.Rectangle{Max: dims.Size},
			gtx.Dp(unit.Dp(10)),
		).Op(gtx.Ops),
	)

	c.Add(gtx.Ops)
	return dims
}

type GatewayListItem struct {
	gateway   app_db.IPFSGateway
	clickable *widget.Clickable
	rounded   unit.Dp
	checkIcon *widget.Icon
}

func NewGatewayListItem(gateway app_db.IPFSGateway) GatewayListItem {
	checkIcon, _ := widget.NewIcon(icons.NavigationCheck)

	return GatewayListItem{
		gateway:   gateway,
		clickable: new(widget.Clickable),
		rounded:   unit.Dp(12),
		checkIcon: checkIcon,
	}
}

func (item *GatewayListItem) Layout(gtx layout.Context, th *material.Theme, fill bool) layout.Dimensions {
	if item.clickable.Clicked(gtx) {
		page_instance.pageEditIPFSGateway.gateway = item.gateway
		page_instance.pageRouter.SetCurrent(PAGE_EDIT_IPFS_GATEWAY)
		page_instance.header.AddHistory(PAGE_EDIT_IPFS_GATEWAY)
	}

	return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		dims := layout.UniformInset(item.rounded).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if item.gateway.Active {
				layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return item.checkIcon.Layout(gtx, theme.Current.ListTextColor)
				})
			}

			return layout.Flex{Alignment: layout.Start}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(18), item.gateway.Name)
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(15), item.gateway.Endpoint)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				}),
			)
		})
		c := r.Stop()

		if item.clickable.Hovered() || fill {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(item.rounded),
				).Op(gtx.Ops),
			)
		}

		c.Add(gtx.Ops)
		return dims
	})
}
