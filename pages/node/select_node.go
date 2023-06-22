package page_node

import (
	"image"
	"image/color"
	"sort"

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
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/node_manager"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/g45t345rt/g45w/ui/components"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSelectNode struct {
	isActive                bool
	animationEnter          *animation.Animation
	animationLeave          *animation.Animation
	buttonSetIntegratedNode *components.Button
	buttonAddNode           *components.Button
	connecting              bool

	nodeList  *NodeList
	listStyle material.ListStyle
}

var _ router.Page = &PageSelectNode{}

func NewPageSelectNode() *PageSelectNode {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .5, ease.OutCubic),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .5, ease.OutCubic),
	))

	th := app_instance.Theme
	list := new(widget.List)
	list.Axis = layout.Vertical
	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay

	nodeList := NewNodeList(th, "You didn't add any remote nodes yet.")

	buttonSetIntegratedNode := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "Use Integrated Node",
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(16),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSetIntegratedNode.Label.Alignment = text.Middle
	buttonSetIntegratedNode.Style.Font.Weight = font.Bold

	addIcon, _ := widget.NewIcon(icons.ContentAddBox)
	buttonAddNode := components.NewButton(components.ButtonStyle{
		Icon:           addIcon,
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &color.NRGBA{A: 255},
		Animation:      components.NewButtonAnimationScale(.92),
	})

	return &PageSelectNode{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		listStyle:      listStyle,

		nodeList:                nodeList,
		buttonSetIntegratedNode: buttonSetIntegratedNode,
		buttonAddNode:           buttonAddNode,
	}
}

func (p *PageSelectNode) IsActive() bool {
	return p.isActive
}

func (p *PageSelectNode) Enter() {
	p.isActive = true

	page_instance.header.SetTitle("Select Node")
	p.animationLeave.Reset()
	p.animationEnter.Start()

	p.Load()
}

func (p *PageSelectNode) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageSelectNode) Load() {
	items := make([]NodeListItem, 0)
	for _, nodeInfo := range node_manager.Instance.NodeState.Nodes {
		items = append(items,
			NewNodeListItem(nodeInfo),
		)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].conn.ID < items[j].conn.ID
	})

	p.nodeList.items = items
}

func (p *PageSelectNode) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.buttonSetIntegratedNode.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(5)}.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(14), "Always use Integrated Node or your own remote node for full privacy and trust.")
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(20), "Remote Nodes")
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Dp(30)
							gtx.Constraints.Min.Y = gtx.Dp(30)
							return p.buttonAddNode.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return p.nodeList.Layout(gtx, th)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
		},
	}

	if p.buttonSetIntegratedNode.Clickable.Clicked() {
		node_manager.Instance.SelectNode(node_manager.INTEGRATED_NODE_ID, true)
		page_instance.pageRouter.SetCurrent(PAGE_INTEGRATED_NODE)
	}

	for _, item := range p.nodeList.items {
		if item.EditClicked() {
			page_instance.pageEditNodeForm.nodeConn = item.conn
			page_instance.pageRouter.SetCurrent(PAGE_EDIT_NODE_FORM)
		}

		if item.SelectClicked() {
			p.connect(item.conn)
		}
	}

	return p.listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(0),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageSelectNode) connect(conn node_manager.NodeConnection) {
	if p.connecting {
		return
	}

	p.connecting = true
	go func() {
		notification_modals.InfoInstance.SetText("Connecting...", conn.Endpoint)
		notification_modals.InfoInstance.SetVisible(true, 0)
		err := node_manager.Instance.SelectNode(conn.ID, true)
		p.connecting = false
		notification_modals.InfoInstance.SetVisible(false, 0)

		if err != nil {
			notification_modals.InfoInstance.SetVisible(false, 0)
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			page_instance.pageRemoteNode.useAnimationEnter = true
			page_instance.pageRouter.SetCurrent(PAGE_REMOTE_NODE)
			app_instance.Window.Invalidate()
		}
	}()
}

type NodeList struct {
	emptyText string
	items     []NodeListItem
	listStyle material.ListStyle
}

func NewNodeList(th *material.Theme, emptyText string) *NodeList {
	list := new(widget.List)
	list.Axis = layout.Vertical

	listStyle := material.List(th, list)
	listStyle.AnchorStrategy = material.Overlay
	listStyle.Indicator.MinorWidth = unit.Dp(10)
	listStyle.Indicator.CornerRadius = unit.Dp(5)
	black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	listStyle.Indicator.Color = black
	//listStyle.Indicator.HoverColor = f32color.Hovered(black)

	return &NodeList{
		listStyle: listStyle,
		emptyText: emptyText,
		items:     make([]NodeListItem, 0),
	}
}

func (l *NodeList) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	r := op.Record(gtx.Ops)
	dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if len(l.items) == 0 {
			lbl := material.Label(th, unit.Sp(14), l.emptyText)
			return lbl.Layout(gtx)
		} else {
			return l.listStyle.Layout(gtx, len(l.items), func(gtx layout.Context, i int) layout.Dimensions {
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

type NodeListItem struct {
	conn               node_manager.NodeConnection
	clickable          *widget.Clickable
	nodeListItemSelect *NodeListItemSelect

	rounded unit.Dp
}

func NewNodeListItem(conn node_manager.NodeConnection) NodeListItem {
	return NodeListItem{
		conn:               conn,
		clickable:          &widget.Clickable{},
		nodeListItemSelect: NewNodeListSelect(),
		rounded:            unit.Dp(12),
	}
}

func (item *NodeListItem) EditClicked() bool {
	return item.nodeListItemSelect.buttonEdit.Clickable.Clicked()
}

func (item *NodeListItem) SelectClicked() bool {
	return item.nodeListItemSelect.buttonSelect.Clickable.Clicked()
}

func (item *NodeListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		dims := layout.UniformInset(item.rounded).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Alignment: layout.Start}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(18), item.conn.Name)
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(15), item.conn.Endpoint)
									lbl.Color.A = 200
									return lbl.Layout(gtx)
								}),
							)
						}),
					)
				}),
			)
		})

		buttonEditHovered := item.nodeListItemSelect.buttonEdit.Clickable.Hovered()
		buttonSelectHovered := item.nodeListItemSelect.buttonSelect.Clickable.Hovered()
		if item.clickable.Hovered() && !buttonEditHovered && !buttonSelectHovered {
			pointer.CursorPointer.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 100},
				clip.UniformRRect(
					image.Rectangle{Max: image.Pt(dims.Size.X, dims.Size.Y)},
					gtx.Dp(item.rounded),
				).Op(gtx.Ops),
			)
		}

		item.nodeListItemSelect.Layout(gtx, th)

		if item.clickable.Clicked() && !buttonEditHovered && !buttonSelectHovered {
			item.nodeListItemSelect.Toggle()
		}

		return dims
	})
}

type NodeListItemSelect struct {
	buttonSelect   *components.Button
	buttonEdit     *components.Button
	visible        bool
	animationEnter *animation.Animation
	animationLeave *animation.Animation
}

func NewNodeListSelect() *NodeListItemSelect {
	buttonSelect := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "SELECT",
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonSelect.Label.Alignment = text.Middle
	buttonSelect.Style.Font.Weight = font.Bold

	buttonEdit := components.NewButton(components.ButtonStyle{
		Rounded:         unit.Dp(5),
		Text:            "EDIT",
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonEdit.Label.Alignment = text.Middle
	buttonEdit.Style.Font.Weight = font.Bold

	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .15, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .15, ease.Linear),
	))

	return &NodeListItemSelect{
		buttonEdit:     buttonEdit,
		buttonSelect:   buttonSelect,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (n *NodeListItemSelect) Toggle() {
	n.SetVisible(!n.visible)
}

func (n *NodeListItemSelect) SetVisible(visible bool) {
	if visible {
		n.visible = true
		n.animationEnter.Start()
		n.animationLeave.Reset()
	} else {
		n.animationEnter.Reset()
		n.animationLeave.Start()
	}
}

func (n *NodeListItemSelect) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if !n.visible {
		return layout.Dimensions{}
	}

	{
		state := n.animationEnter.Update(gtx)
		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}
	}

	{
		state := n.animationLeave.Update(gtx)

		if state.Active {
			defer animation.TransformX(gtx, state.Value).Push(gtx.Ops).Pop()
		}

		if state.Finished {
			n.visible = false
			op.InvalidateOp{}.Add(gtx.Ops)
		}
	}

	return layout.E.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return n.buttonSelect.Layout(gtx, th)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return n.buttonEdit.Layout(gtx, th)
			}),
		)
	})
}
