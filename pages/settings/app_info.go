package page_settings

import (
	"fmt"
	"image"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAppInfo struct {
	isActive       bool
	list           *widget.List
	animationEnter *animation.Animation
	animationLeave *animation.Animation
	infoItems      []*InfoListItem
}

var _ router.Page = &PageAppInfo{}

func NewPageAppInfo() *PageAppInfo {
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
	// it's used by the python script to generate language json dictionary
	// we don't use lang.Translate directly here because it needs to be inside the Layout func or the value won't be updated after language change
	infoItems := []*InfoListItem{
		NewInfoListItem("App Directory", settings.AppDir),             //@lang.Translate("App Directory")
		NewInfoListItem("Node Directory", settings.IntegratedNodeDir), //@lang.Translate("Node Directory")
		NewInfoListItem("Wallets Directory", settings.WalletsDir),     //@lang.Translate("Wallets Directory")
		NewInfoListItem("Cache Directory", settings.CacheDir),         //@lang.Translate("Cache Directory")
		NewInfoListItem("Version", settings.Version),                  //@lang.Translate("Version")
		NewInfoListItem("Git Version", settings.GitVersion),           //@lang.Translate("Git Version")
		NewInfoListItem("Build Time", buildTime),                      //@lang.Translate("Build Time")
		NewInfoListItem("Donation Address", settings.DonationAddress), //@lang.Translate("Donation")
	}

	return &PageAppInfo{
		infoItems:      infoItems,
		list:           list,
		animationEnter: animationEnter,
		animationLeave: animationLeave,
	}
}

func (p *PageAppInfo) IsActive() bool {
	return p.isActive
}

func (p *PageAppInfo) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("App Information") }

	if !page_instance.header.IsHistory(PAGE_APP_INFO) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageAppInfo) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageAppInfo) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	var widgets []layout.Widget

	for i := range p.infoItems {
		idx := i
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.infoItems[idx].Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay
	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return widgets[index](gtx)
	})
}

type InfoListItem struct {
	title      string
	editor     *widget.Editor
	buttonCopy *components.Button
}

func NewInfoListItem(title string, value string) *InfoListItem {
	editor := new(widget.Editor)
	editor.WrapPolicy = text.WrapGraphemes
	editor.ReadOnly = true
	editor.SetText(value)

	copyIcon, _ := widget.NewIcon(icons.ContentContentCopy)
	buttonCopy := components.NewButton(components.ButtonStyle{
		Icon: copyIcon,
	})

	return &InfoListItem{
		title:      title,
		editor:     editor,
		buttonCopy: buttonCopy,
	}
}

func (s InfoListItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if s.buttonCopy.Clicked() {
		clipboard.WriteOp{
			Text: s.editor.Text(),
		}.Add(gtx.Ops)
		notification_modals.InfoInstance.SetText(lang.Translate("Clipboard"), lang.Translate("Text copied to clipboard"))
		notification_modals.InfoInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
	}

	dims := layout.Inset{
		Top: unit.Dp(10), Bottom: unit.Dp(10),
		Left: unit.Dp(10), Right: unit.Dp(10),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(18), lang.Translate(s.title))
						lbl.Font.Weight = font.Bold
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Max.X = gtx.Dp(20)
						gtx.Constraints.Max.Y = gtx.Dp(20)
						s.buttonCopy.Style.Colors = theme.Current.ModalButtonColors
						return s.buttonCopy.Layout(gtx, th)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(th, s.editor, "")
				return editor.Layout(gtx)
			}),
		)
	})

	// Divider
	cl := clip.Rect{Max: image.Pt(dims.Size.X, gtx.Dp(1))}.Push(gtx.Ops)
	paint.ColorOp{Color: theme.Current.DividerColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	cl.Pop()

	return dims
}
