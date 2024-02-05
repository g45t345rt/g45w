package page_wallet

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/walletapi/xswd"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageXSWDManage struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation
	buttonStart         *components.Button
	buttonStop          *components.Button

	list *widget.List
	apps []*DAppItem
}

var _ router.Page = &PageXSWDManage{}

func NewPageXSWDManage() *PageXSWDManage {

	playIcon, _ := widget.NewIcon(icons.AVPlayArrow)
	buttonStart := components.NewButton(components.ButtonStyle{
		Icon: playIcon,
	})

	stopIcon, _ := widget.NewIcon(icons.AVPause)
	buttonStop := components.NewButton(components.ButtonStyle{
		Icon: stopIcon,
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_XSWD_MANAGE)
	return &PageXSWDManage{
		headerPageAnimation: headerPageAnimation,
		buttonStop:          buttonStop,
		buttonStart:         buttonStart,
		list:                list,
		apps:                make([]*DAppItem, 0),
	}
}

func (p *PageXSWDManage) IsActive() bool {
	return p.isActive
}

func (p *PageXSWDManage) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string {
		return lang.Translate("DApp connections")
	}

	page_instance.header.Subtitle = nil
	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		p.buttonStart.Style.Colors = theme.Current.ButtonIconPrimaryColors
		p.buttonStop.Style.Colors = theme.Current.ButtonIconPrimaryColors
		gtx.Constraints.Min.X = gtx.Dp(30)
		gtx.Constraints.Min.Y = gtx.Dp(30)

		wallet := wallet_manager.OpenedWallet
		xswd := wallet.ServerXSWD

		if xswd != nil && xswd.IsRunning() {
			if p.buttonStop.Clicked(gtx) {
				go page_instance.CloseXSWD()
			}

			return p.buttonStop.Layout(gtx, th)
		}

		if p.buttonStart.Clicked(gtx) {
			go func() {
				err := page_instance.OpenXSWD()
				if err != nil {
					notification_modal.Open(notification_modal.Params{
						Type:  notification_modal.ERROR,
						Title: lang.Translate("Error"),
						Text:  err.Error(),
					})
				}
				p.Load()
			}()
		}

		return p.buttonStart.Layout(gtx, th)
	}

	p.Load()
}

func (p *PageXSWDManage) Load() {
	xswd := wallet_manager.OpenedWallet.ServerXSWD
	if xswd != nil && xswd.IsRunning() {
		connectedApps := xswd.GetApplications()
		p.apps = make([]*DAppItem, 0)
		for _, app := range connectedApps {
			p.apps = append(p.apps, NewDAppItem(app))
		}
	}
}

func (p *PageXSWDManage) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageXSWDManage) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	widgets := []layout.Widget{}

	xswd := wallet_manager.OpenedWallet.ServerXSWD
	if xswd != nil && xswd.IsRunning() {
		if len(p.apps) > 0 {
			for i := range p.apps {
				app := p.apps[i]
				widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
					return app.Layout(gtx, th)
				})
			}
		} else {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, 16, lang.Translate("There are currently no dApp connections."))
				lbl.Color = theme.Current.TextMuteColor
				return lbl.Layout(gtx)
			})
		}
	} else {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, 16, lang.Translate("XSWD is not running."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(10),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

type DAppItem struct {
	app       xswd.ApplicationData
	clickable *widget.Clickable
}

func NewDAppItem(app xswd.ApplicationData) *DAppItem {
	return &DAppItem{
		app:       app,
		clickable: new(widget.Clickable),
	}
}

func (item *DAppItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Clicked(gtx) {
		page_instance.pageXSWDApp.app = item.app
		page_instance.pageRouter.SetCurrent(PAGE_XSWD_APP)
		page_instance.header.AddHistory(PAGE_XSWD_APP)
	}

	m := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(13), Bottom: unit.Dp(13),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {

					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(18), item.app.Name)
									lbl.Font.Weight = font.Bold
									return lbl.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(14), item.app.Url)
									lbl.Color = theme.Current.TextMuteColor
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				}),
			)
		})
	})
	c := m.Stop()

	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
		paint.FillShape(gtx.Ops, theme.Current.ListItemHoverBgColor,
			clip.UniformRRect(
				image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
				gtx.Dp(10),
			).Op(gtx.Ops),
		)
	} else {
		paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
			clip.RRect{
				Rect: image.Rectangle{Max: dims.Size},
				NW:   gtx.Dp(10), NE: gtx.Dp(10),
				SE: gtx.Dp(10), SW: gtx.Dp(10),
			}.Op(gtx.Ops),
		)
	}

	c.Add(gtx.Ops)

	return dims
}
