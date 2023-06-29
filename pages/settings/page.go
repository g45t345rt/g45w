package page_settings

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/containers/bottom_bar"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
)

type Page struct {
	isActive bool
	list     *widget.List

	settingsItems []*SettingsItem

	langSelector *prefabs.LangSelector
}

var _ router.Page = &Page{}

func New() *Page {
	langSelector := prefabs.NewLangSelector("en")
	list := new(widget.List)
	list.Axis = layout.Vertical

	unix, _ := strconv.ParseUint(settings.BuildTime, 10, 64)
	buildTimeUnix := time.Unix(int64(unix), 0)
	buildTime := fmt.Sprintf("%s (%d)", buildTimeUnix.Local().String(), unix)

	// do not remove @lang.Translate comment
	// it's used by the python to generate language json dictionary
	// we don't use lang.Translate directly here because it needs to be inside the Layout func or the value won't be updated after language change
	settingsItems := []*SettingsItem{
		NewSettingsItem("App Directory", settings.AppDir),         //@lang.Translate("App Directory")
		NewSettingsItem("Node Directory", settings.NodeDir),       //@lang.Translate("Node Directory")
		NewSettingsItem("Wallets Directory", settings.WalletsDir), //@lang.Translate("Wallets Directory")
		NewSettingsItem("Version", settings.Version),              //@lang.Translate("Version")
		NewSettingsItem("Git Version", settings.GitVersion),       //@lang.Translate("Git Version")
		NewSettingsItem("Build Time", buildTime),                  //@lang.Translate("Build Time")
	}

	return &Page{
		list:          list,
		settingsItems: settingsItems,
		langSelector:  langSelector,
	}
}

func (p *Page) IsActive() bool {
	return p.isActive
}

func (p *Page) Enter() {
	bottom_bar.Instance.SetButtonActive(bottom_bar.BUTTON_SETTINGS)
	p.isActive = true
}

func (p *Page) Leave() {
	p.isActive = false
}

func (p *Page) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {

	var widgets []layout.ListElement

	for range p.settingsItems {
		widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
			return p.settingsItems[index].Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
		if p.langSelector.Changed() {
			settings.App.Language = p.langSelector.Value
			err := settings.Save()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			} else {
				notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Language applied."))
				notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			}
		}

		return p.langSelector.Layout(gtx, th)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return widgets[index](gtx, index)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return bottom_bar.Instance.Layout(gtx, th)
		}),
	)
}

type SettingsItem struct {
	title  string
	editor *widget.Editor
}

func NewSettingsItem(title string, value string) *SettingsItem {
	editor := new(widget.Editor)
	editor.WrapPolicy = text.WrapGraphemes
	editor.ReadOnly = true
	editor.SetText(value)

	return &SettingsItem{
		title:  title,
		editor: editor,
	}
}

func (s SettingsItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
