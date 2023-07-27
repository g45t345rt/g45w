package page_settings

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
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_data"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageIPFSGateways struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation
	buttonInfo     *components.Button
	modalInfo      *components.Modal

	gatewayList *GatewayList

	buttonAdd *components.Button
}

var _ router.Page = &PageIPFSGateways{}

func NewPageIPFSGateways() *PageIPFSGateways {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Icon:      addIcon,
		TextColor: color.NRGBA{A: 255},
	})

	gatewayList := NewGatewayList()

	infoIcon, _ := widget.NewIcon(icons.ActionInfo)
	buttonInfo := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            infoIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})

	modalInfo := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  true,
		Direction:           layout.Center,
		Inset:               layout.UniformInset(unit.Dp(30)),
		Rounded:             components.UniformRounded(unit.Dp(10)),
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Animation:           components.NewModalAnimationScaleBounce(),
		Backdrop:            components.NewModalBackground(),
	})

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
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

	return &PageIPFSGateways{
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		gatewayList:    gatewayList,
		buttonAdd:      buttonAdd,
		buttonInfo:     buttonInfo,
		modalInfo:      modalInfo,
	}
}

func (p *PageIPFSGateways) IsActive() bool {
	return p.isActive
}

func (p *PageIPFSGateways) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("IPFS Gateways"))
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("Interplanetary File System"))
		return lbl.Layout(gtx)
	}
	page_instance.header.ButtonRight = p.buttonAdd

	if !page_instance.header.IsHistory(PAGE_IPFS_GATEWAYS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.LoadGateways()
}

func (p *PageIPFSGateways) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageIPFSGateways) LoadGateways() error {
	items := make([]GatewayListItem, 0)

	gateways, err := app_data.GetIPFSGateways(app_data.GetIPFSGatewaysParams{})
	if err != nil {
		return err
	}

	for _, gateway := range gateways {
		items = append(items, NewGatewayListItem(gateway))
	}

	p.gatewayList.items = items
	return nil
}

func (p *PageIPFSGateways) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	{
		state := p.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := p.animationLeave.Update(gtx)
		if state.Finished {
			p.isActive = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	if p.buttonAdd.Clicked() {
		page_instance.pageRouter.SetCurrent(PAGE_ADD_IPFS_GATEWAY)
		page_instance.header.AddHistory(PAGE_ADD_IPFS_GATEWAY)
	}

	p.buttonInfo.Text = lang.Translate("Why use IPFS?")
	if p.buttonInfo.Clicked() {
		p.modalInfo.SetVisible(true)
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.buttonInfo.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.gatewayList.Layout(gtx, th, lang.Translate("You don't have any IPFS gateways available."))
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

type GatewayList struct {
	items []GatewayListItem
	list  *widget.List
}

func NewGatewayList() *GatewayList {
	list := new(widget.List)
	list.Axis = layout.Vertical

	return &GatewayList{
		list:  list,
		items: make([]GatewayListItem, 0),
	}
}

func (l *GatewayList) Layout(gtx layout.Context, th *material.Theme, emptyText string) layout.Dimensions {
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
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
			listStyle.Indicator.Color = black

			return listStyle.Layout(gtx, len(l.items), func(gtx layout.Context, i int) layout.Dimensions {
				return l.items[i].Layout(gtx, th)
			})
		}
	})
	c := r.Stop()

	paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		clip.UniformRRect(
			image.Rectangle{Max: dims.Size},
			gtx.Dp(unit.Dp(10)),
		).Op(gtx.Ops),
	)

	c.Add(gtx.Ops)
	return dims
}

type GatewayListItem struct {
	gateway   app_data.IPFSGateway
	clickable *widget.Clickable
	rounded   unit.Dp
	checkIcon *widget.Icon
}

func NewGatewayListItem(gateway app_data.IPFSGateway) GatewayListItem {
	checkIcon, _ := widget.NewIcon(icons.NavigationCheck)

	return GatewayListItem{
		gateway:   gateway,
		clickable: new(widget.Clickable),
		rounded:   unit.Dp(12),
		checkIcon: checkIcon,
	}
}

func (item *GatewayListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		dims := layout.UniformInset(item.rounded).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if item.gateway.Active {
				layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return item.checkIcon.Layout(gtx, color.NRGBA{A: 255})
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
									lbl.Color.A = 200
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				}),
			)
		})

		if item.clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(item.rounded),
				).Op(gtx.Ops),
			)
		}

		if item.clickable.Clicked() {
			page_instance.pageEditIPFSGateway.gateway = item.gateway
			page_instance.pageRouter.SetCurrent(PAGE_EDIT_IPFS_GATEWAY)
			page_instance.header.AddHistory(PAGE_EDIT_IPFS_GATEWAY)
		}

		return dims
	})
}
