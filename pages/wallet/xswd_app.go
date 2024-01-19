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
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageXSWDApp struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonRemove *components.Button
	list         *widget.List
	app          xswd.ApplicationData
	permissions  []*DAppPermissionItem
}

var _ router.Page = &PageXSWDApp{}

func NewPageXSWDApp() *PageXSWDApp {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	removeIcon, _ := widget.NewIcon(icons.NavigationCancel)
	buttonRemove := components.NewButton(components.ButtonStyle{
		Icon: removeIcon,
	})

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageXSWDApp{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
		buttonRemove:   buttonRemove,
	}
}

func (p *PageXSWDApp) IsActive() bool {
	return p.isActive
}

func (p *PageXSWDApp) Load() {
	p.permissions = make([]*DAppPermissionItem, 0)
	for name, perm := range p.app.Permissions {
		p.permissions = append(p.permissions, NewDAppPermissionItem(name, perm))
	}
}

func (p *PageXSWDApp) Enter() {
	p.isActive = true

	page_instance.header.Title = func() string {
		return p.app.Name
	}

	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(14), p.app.Url)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		p.buttonRemove.Style.Colors = theme.Current.ButtonIconPrimaryColors
		gtx.Constraints.Min.X = gtx.Dp(30)
		gtx.Constraints.Min.Y = gtx.Dp(30)

		xswd := wallet_manager.OpenedWallet.ServerXSWD
		if p.buttonRemove.Clicked(gtx) {
			go func() {
				xswd.RemoveApplication(&p.app)
				page_instance.header.GoBack()
				page_instance.pageXSWDManage.Load()
				notification_modal.Open(notification_modal.Params{
					Type:       notification_modal.SUCCESS,
					Title:      lang.Translate("Success"),
					Text:       lang.Translate("Connection was terminated."),
					CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
				})
			}()
		}

		return p.buttonRemove.Layout(gtx, th)
	}

	if !page_instance.header.IsHistory(PAGE_XSWD_APP) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.Load()
}

func (p *PageXSWDApp) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageXSWDApp) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	for i := range p.permissions {
		perm := p.permissions[i]
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return perm.Layout(gtx, th)
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

type DAppPermissionItem struct {
	name      string
	perm      xswd.Permission
	clickable *widget.Clickable
}

func NewDAppPermissionItem(name string, perm xswd.Permission) *DAppPermissionItem {
	return &DAppPermissionItem{
		name:      name,
		perm:      perm,
		clickable: new(widget.Clickable),
	}
}

func (item *DAppPermissionItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Clicked(gtx) {
		// TODO
	}

	m := op.Record(gtx.Ops)
	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(13), Bottom: unit.Dp(13),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), item.name)
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), item.perm.String())
					return lbl.Layout(gtx)
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
