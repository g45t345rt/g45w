package page_wallet

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	crypto "github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCFunction struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonExecute *components.Button
	scArgItems    []*SCArgItem
	scFunction    SCFunction
	SCID          string
	list          *widget.List
}

var _ router.Page = &PageSCFunction{}

func NewPageSCFunction() *PageSCFunction {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	validIcon, _ := widget.NewIcon(icons.ActionCheckCircle)
	buttonExecute := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      validIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonExecute.Label.Alignment = text.Middle
	buttonExecute.Style.Font.Weight = font.Bold

	return &PageSCFunction{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		buttonExecute:  buttonExecute,

		list: list,
	}
}

func (p *PageSCFunction) IsActive() bool {
	return p.isActive
}

func (p *PageSCFunction) Enter() {
	p.isActive = true

	page_instance.header.Title = func() string {
		return p.scFunction.Name
	}

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = nil

	if !page_instance.header.IsHistory(PAGE_SC_FUNCTION) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageSCFunction) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSCFunction) SetData(SCID string, scFunction SCFunction) {
	p.SCID = SCID
	p.scFunction = scFunction
	p.scArgItems = make([]*SCArgItem, 0)
	for _, arg := range p.scFunction.Args {
		p.scArgItems = append(p.scArgItems, NewSCArgItem(arg))
	}
}

func (p *PageSCFunction) execute() {
	args := rpc.Arguments{
		{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
		{Name: rpc.SCID, DataType: rpc.DataHash, Value: crypto.HashHexToHash(p.SCID)},
		{Name: "entrypoint", DataType: rpc.DataString, Value: p.scFunction.Name},
	}

	for _, item := range p.scArgItems {
		dataType := rpc.DataString
		if item.arg.Type == "Uint64" {
			dataType = rpc.DataUint64
		}
		name := item.arg.Name
		value := item.txtValue.Value()

		args = append(args, rpc.Argument{
			Name: name, DataType: dataType, Value: value,
		})
	}

	build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(randomAddr string) build_tx_modal.TxPayload {
		return build_tx_modal.TxPayload{
			SCArgs:    args,
			Transfers: []rpc.Transfer{
				//rpc.Transfer{SCID: token1.GetHash(), Burn: amount1.Number, Destination: randomAddr},
				//rpc.Transfer{SCID: token2.GetHash(), Burn: amount2.Number, Destination: randomAddr},
			},
			Ringsize: 2,
			//TokensInfo: []*wallet_manager.Token{token1, token2},
		}
	})
}

func (p *PageSCFunction) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonExecute.Clicked(gtx) {
		go p.execute()
	}

	widgets := []layout.Widget{}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if len(p.scArgItems) == 0 {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("This function does not have any arguments."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("Arguments"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	})

	for i := range p.scArgItems {
		item := p.scArgItems[i]
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return item.Layout(gtx, th)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("Transfers"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("TODO!!"))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{}
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonExecute.Style.Colors = theme.Current.ButtonPrimaryColors
		p.buttonExecute.Text = lang.Translate("VALIDATE FUNCTION")
		return p.buttonExecute.Layout(gtx, th)
	})

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

type SCArgItem struct {
	arg      SCFunctionArg
	txtValue *prefabs.TextField
}

func NewSCArgItem(arg SCFunctionArg) *SCArgItem {
	txtValue := prefabs.NewTextField()

	if arg.Type == "Uint64" {
		txtValue = prefabs.NewNumberTextField()
	}

	return &SCArgItem{
		arg:      arg,
		txtValue: txtValue,
	}
}

func (item *SCArgItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	title := fmt.Sprintf("%s (%s)", item.arg.Name, item.arg.Type)
	return item.txtValue.Layout(gtx, th, title, "")
}
