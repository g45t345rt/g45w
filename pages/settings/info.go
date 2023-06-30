package page_settings

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/ui/animation"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageInfo struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation
	infoItems      []*InfoListItem
}

var _ router.Page = &PageInfo{}

func NewPageInfo() *PageInfo {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	unix, _ := strconv.ParseUint(settings.BuildTime, 10, 64)
	buildTimeUnix := time.Unix(int64(unix), 0)
	buildTime := fmt.Sprintf("%s (%d)", buildTimeUnix.Local().String(), unix)

	// do not remove @lang.Translate comment
	// it's used by the python to generate language json dictionary
	// we don't use lang.Translate directly here because it needs to be inside the Layout func or the value won't be updated after language change
	infoItems := []*InfoListItem{
		NewInfoListItem("App Directory", settings.AppDir),         //@lang.Translate("App Directory")
		NewInfoListItem("Node Directory", settings.NodeDir),       //@lang.Translate("Node Directory")
		NewInfoListItem("Wallets Directory", settings.WalletsDir), //@lang.Translate("Wallets Directory")
		NewInfoListItem("Version", settings.Version),              //@lang.Translate("Version")
		NewInfoListItem("Git Version", settings.GitVersion),       //@lang.Translate("Git Version")
		NewInfoListItem("Build Time", buildTime),                  //@lang.Translate("Build Time")
	}

	return &PageInfo{
		infoItems:      infoItems,
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (p *PageInfo) IsActive() bool {
	return p.isActive
}

func (p *PageInfo) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("App Information"))
	p.animationEnter.Start()
	p.animationLeave.Reset()
}

func (p *PageInfo) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageInfo) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay
	return listStyle.Layout(gtx, len(p.infoItems), func(gtx layout.Context, index int) layout.Dimensions {
		return p.infoItems[index].Layout(gtx, th)
	})
}

type InfoListItem struct {
	title  string
	editor *widget.Editor
}

func NewInfoListItem(title string, value string) *InfoListItem {
	editor := new(widget.Editor)
	editor.WrapPolicy = text.WrapGraphemes
	editor.ReadOnly = true
	editor.SetText(value)

	return &InfoListItem{
		title:  title,
		editor: editor,
	}
}

func (s InfoListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	dims := layout.Inset{
		Top: unit.Dp(10), Bottom: unit.Dp(10),
		Left: unit.Dp(10), Right: unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Label(th, unit.Sp(18), lang.Translate(s.title))
				label.Font.Weight = font.Bold
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, s.editor, "")
				return editor.Layout(gtx)
			}),
		)
	})

	cl := clip.Rect{Max: image.Pt(dims.Size.X, gtx.Dp(1))}.Push(gtx.Ops)
	paint.ColorOp{Color: color.NRGBA{A: 50}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	return dims
}
