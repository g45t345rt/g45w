package page_wallet

import (
	"encoding/hex"
	"fmt"
	"image"
	"regexp"
	"sort"
	"strings"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type SCFunction struct {
	Name string
	Args []SCFunctionArg
}

type SCFunctionArg struct {
	Name string
	Type string
}

type PageSCExplorer struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation
	tabBars        *components.TabBars

	buttonMenu *components.Button
	scFuncs    []*SCFunctionItem
	scData     []*SCDataItem
	immutable  bool

	SCID string
	list *widget.List
}

var _ router.Page = &PageSCExplorer{}

func NewPageSCExplorer() *PageSCExplorer {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	tabBarsItems := []*components.TabBarsItem{
		components.NewTabBarItem("functions"),
		components.NewTabBarItem("data"),
	}

	tabBars := components.NewTabBars("functions", tabBarsItems)

	list := new(widget.List)
	list.Axis = layout.Vertical

	menuIcon, _ := widget.NewIcon(icons.NavigationMenu)
	buttonMenu := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      menuIcon,
		Animation: components.NewButtonAnimationDefault(),
	})

	return &PageSCExplorer{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		tabBars:        tabBars,
		buttonMenu:     buttonMenu,

		list: list,
	}
}

func (p *PageSCExplorer) IsActive() bool {
	return p.isActive
}

func (p *PageSCExplorer) LoadFunctions() error {
	var result rpc.GetSC_Result
	err := wallet_manager.RPCCall("DERO.GetSC", rpc.GetSC_Params{
		SCID:      p.SCID,
		Code:      true,
		Variables: false,
	}, &result)
	if err != nil {
		return err
	}

	matchFunctions, err := regexp.Compile(`Function ([A-Z]\w+)\(?(.+)\)`)
	if err != nil {
		return err
	}

	p.scFuncs = make([]*SCFunctionItem, 0)
	p.immutable = true

	values := matchFunctions.FindAllStringSubmatch(result.Code, -1)
	for _, value := range values {
		funcName := value[1]
		if funcName == "Initialize" || funcName == "InitializePrivate" {
			continue
		}

		if funcName == "UpdateCode" {
			p.immutable = false
		}

		scFunc := SCFunction{
			Name: funcName,
		}

		sArgs := value[2]
		if sArgs != "(" {
			args := strings.Split(sArgs, ",")

			for _, arg := range args {
				def := strings.Split(strings.Trim(arg, " "), " ")
				scFunc.Args = append(scFunc.Args, SCFunctionArg{
					Name: def[0],
					Type: def[1],
				})
			}
		}

		p.scFuncs = append(p.scFuncs, NewSCFunctionItem(scFunc))
	}

	// sort by name alphabetically
	sort.Slice(p.scFuncs, func(i, j int) bool {
		return p.scFuncs[i].scFunc.Name < p.scFuncs[j].scFunc.Name
	})

	return nil
}

func (p *PageSCExplorer) LoadData() error {
	p.scData = make([]*SCDataItem, 0)
	var result rpc.GetSC_Result
	err := wallet_manager.RPCCall("DERO.GetSC", rpc.GetSC_Params{
		SCID:      p.SCID,
		Code:      false,
		Variables: true,
	}, &result)
	if err != nil {
		return err
	}

	for key, data := range result.VariableStringKeys {
		if key == "C" {
			continue
		}

		p.scData = append(p.scData, NewSCDataItem(key, data))
	}

	// sort by keys alphabetically
	sort.Slice(p.scData, func(i, j int) bool {
		return p.scData[i].key < p.scData[j].key
	})

	return nil
}

func (p *PageSCExplorer) Enter() {
	p.isActive = true

	page_instance.header.Title = func() string {
		return lang.Translate("SC Explorer")
	}

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		p.buttonMenu.Style.Colors = theme.Current.ButtonIconPrimaryColors
		gtx.Constraints.Min.X = gtx.Dp(30)
		gtx.Constraints.Min.Y = gtx.Dp(30)

		if p.buttonMenu.Clicked(gtx) {
			go func() {
				codeIcon, _ := widget.NewIcon(icons.ActionCode)
				refreshIcon, _ := widget.NewIcon(icons.NavigationRefresh)

				keyChan := listselect_modal.Instance.Open([]*listselect_modal.SelectListItem{
					listselect_modal.NewSelectListItem("view_code",
						listselect_modal.NewItemText(codeIcon, lang.Translate("View code")).Layout,
					),
					listselect_modal.NewSelectListItem("reload_sc",
						listselect_modal.NewItemText(refreshIcon, lang.Translate("Reload SC")).Layout,
					),
				}, "")

				for key := range keyChan {
					switch key {
					case "view_code":
					case "reload_sc":
					}
				}
			}()
		}

		return p.buttonMenu.Layout(gtx, th)
	}

	if !page_instance.header.IsHistory(PAGE_SC_EXPLORER) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.LoadFunctions()
	p.LoadData()
}

func (p *PageSCExplorer) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSCExplorer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		txt := lang.Translate("This smart contract is mutable.")
		if p.immutable {
			txt = lang.Translate("This smart contract is immutable.")
		}

		lbl := material.Label(th, unit.Sp(16), txt)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		text := make(map[string]string)

		txt := lang.Translate("Functions ({})")
		txt = strings.Replace(txt, "{}", fmt.Sprint(len(p.scFuncs)), -1)
		text["functions"] = txt

		txt = lang.Translate("Data ({})")
		txt = strings.Replace(txt, "{}", fmt.Sprint(len(p.scData)), -1)
		text["data"] = txt

		p.tabBars.Colors = theme.Current.TabBarsColors
		return p.tabBars.Layout(gtx, th, unit.Sp(18), text)
	})

	if p.tabBars.Key == "functions" {
		if len(p.scFuncs) == 0 {
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(16), lang.Translate("This smart contract does not have any functions."))
				lbl.Color = theme.Current.TextMuteColor
				return lbl.Layout(gtx)
			})
		}

		for i := range p.scFuncs {
			item := p.scFuncs[i]
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return item.Layout(gtx, th)
			})
		}

		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
		})
	}

	if p.tabBars.Key == "data" {
		for i := range p.scData {
			item := p.scData[i]
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return item.Layout(gtx, th)
			})
		}
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
	})

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(10),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

type SCDataItem struct {
	key    string
	editor *widget.Editor
}

func NewSCDataItem(key string, data interface{}) *SCDataItem {
	editor := &widget.Editor{}
	editor.ReadOnly = true

	value := fmt.Sprintf("%v", data)

	decoded, err := hex.DecodeString(value)
	if err == nil {
		editor.SetText(string(decoded))
	} else {
		editor.SetText(value)
	}

	// check if address is raw
	p := new(crypto.Point)
	err = p.DecodeCompressed(decoded)
	if err == nil {
		addr := rpc.NewAddressFromKeys(p)
		editor.SetText(addr.String())
	}

	return &SCDataItem{
		key:    key,
		editor: editor,
	}
}

func (item *SCDataItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), item.key)
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(3)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			editor := material.Editor(th, item.editor, "")
			editor.Color = theme.Current.TextMuteColor
			editor.TextSize = unit.Sp(14)
			return editor.Layout(gtx)
		}),
	)
}

type SCFunctionItem struct {
	scFunc    SCFunction
	clickable *widget.Clickable
}

func NewSCFunctionItem(scFunc SCFunction) *SCFunctionItem {
	return &SCFunctionItem{
		scFunc:    scFunc,
		clickable: new(widget.Clickable),
	}
}

func (item *SCFunctionItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Clicked(gtx) {
		scid := page_instance.pageSCExplorer.SCID
		page_instance.pageSCFunction.SetData(scid, item.scFunc)
		page_instance.pageRouter.SetCurrent(PAGE_SC_FUNCTION)
		page_instance.header.AddHistory(PAGE_SC_FUNCTION)
	}

	m := op.Record(gtx.Ops)

	dims := item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		if item.clickable.Hovered() {
			pointer.CursorPointer.Add(gtx.Ops)
		}

		return layout.Inset{
			Top: unit.Dp(13), Bottom: unit.Dp(13),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), item.scFunc.Name)
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					count := len(item.scFunc.Args)
					txt := fmt.Sprintf("%d args", count)
					lbl := material.Label(th, unit.Sp(14), txt)
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		})
	})
	c := m.Stop()

	paint.FillShape(gtx.Ops, theme.Current.ListBgColor,
		clip.RRect{
			Rect: image.Rectangle{Max: dims.Size},
			NW:   gtx.Dp(10), NE: gtx.Dp(10),
			SE: gtx.Dp(10), SW: gtx.Dp(10),
		}.Op(gtx.Ops))

	c.Add(gtx.Ops)

	return dims
}
