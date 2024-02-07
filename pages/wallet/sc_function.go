package page_wallet

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	crypto "github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCFunction struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonExecute     *components.Button
	buttonAddTransfer *components.Button
	scArgItems        []*SCArgItem
	scTransferItems   []*SCTransferItem
	scFunction        SCFunction
	SCID              string
	list              *widget.List
}

var _ router.Page = &PageSCFunction{}

func NewPageSCFunction() *PageSCFunction {

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

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAddTransfer := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Icon:      addIcon,
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonExecute.Label.Alignment = text.Middle
	buttonExecute.Style.Font.Weight = font.Bold

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_SC_FUNCTION)
	return &PageSCFunction{
		headerPageAnimation: headerPageAnimation,
		buttonExecute:       buttonExecute,
		buttonAddTransfer:   buttonAddTransfer,

		list: list,
	}
}

func (p *PageSCFunction) IsActive() bool {
	return p.isActive
}

func (p *PageSCFunction) Enter() {
	p.scTransferItems = make([]*SCTransferItem, 0)
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string {
		return p.scFunction.Name
	}

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = nil
}

func (p *PageSCFunction) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageSCFunction) SetData(SCID string, scFunction SCFunction) {
	p.SCID = SCID
	p.scFunction = scFunction
	p.scArgItems = make([]*SCArgItem, 0)
	for _, arg := range p.scFunction.Args {
		p.scArgItems = append(p.scArgItems, NewSCArgItem(arg))
	}
}

func (p *PageSCFunction) addTransfer() {
	onDelete := func(index int) {
		p.scTransferItems = append(p.scTransferItems[:index], p.scTransferItems[index+1:]...)
		app_instance.Window.Invalidate()
	}

	item := NewSCTransferItem(onDelete)
	p.scTransferItems = append(p.scTransferItems, item)
	app_instance.Window.Invalidate()
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

	var transfers []rpc.Transfer

	formatTransfers := func() error {
		for _, item := range p.scTransferItems {
			burn, err := strconv.ParseUint(item.amountInput.Value(), 10, 64)
			if err != nil {
				return err
			}

			scId := item.scIdInput.Value()
			byteSlice, err := hex.DecodeString(scId)
			if err != nil {
				return err
			}

			if len(byteSlice) != 32 {
				return fmt.Errorf("invalid scid")
			}

			transfers = append(transfers, rpc.Transfer{
				SCID: crypto.HexToHash(scId),
				Burn: burn,
			})
		}
		return nil
	}

	err := formatTransfers()
	if err != nil {
		notification_modal.Open(notification_modal.Params{
			Type:  notification_modal.ERROR,
			Title: lang.Translate("Error"),
			Text:  err.Error(),
		})
		return
	}

	build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(randomAddr string) build_tx_modal.TxPayload {
		for _, transfer := range transfers {
			transfer.Destination = randomAddr
		}

		return build_tx_modal.TxPayload{
			SCArgs:    args,
			Transfers: transfers,
			Ringsize:  2,
			//TokensInfo: []*wallet_manager.Token{token1, token2},
		}
	})
}

func (p *PageSCFunction) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonExecute.Clicked(gtx) {
		go p.execute()
	}

	if p.buttonAddTransfer.Clicked(gtx) {
		p.addTransfer()
	}

	widgets := []layout.Widget{}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(18), lang.Translate("Arguments"))
		return lbl.Layout(gtx)
	})

	if len(p.scArgItems) == 0 {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("This function does not have any arguments."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})
	} else {
		for i := range p.scArgItems {
			item := p.scArgItems[i]
			widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
				return item.Layout(gtx, th)
			})
		}
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(18), lang.Translate("Transfers"))
		return lbl.Layout(gtx)
	})

	for i := range p.scTransferItems {
		index := i
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			if index < len(p.scTransferItems) {
				return p.scTransferItems[index].Layout(gtx, th, index)
			}

			return layout.Dimensions{}
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonAddTransfer.Text = lang.Translate("Add Transfer")
		p.buttonAddTransfer.Style.Colors = theme.Current.ButtonSecondaryColors
		return p.buttonAddTransfer.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return prefabs.Divider(gtx, unit.Dp(5))
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		)
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

type SCTransferItem struct {
	onDelete     func(index int)
	scIdInput    *prefabs.Input
	amountInput  *prefabs.Input
	buttonRemove *components.Button
}

func NewSCTransferItem(onDelete func(index int)) *SCTransferItem {
	scIdInput := prefabs.NewInput()
	amountInput := prefabs.NewNumberInput()

	removeIcon, _ := widget.NewIcon(icons.ActionDelete)
	buttonRemove := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Icon:      removeIcon,
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonRemove.Label.Alignment = text.Middle
	buttonRemove.Style.Font.Weight = font.Bold

	return &SCTransferItem{
		onDelete:     onDelete,
		scIdInput:    scIdInput,
		amountInput:  amountInput,
		buttonRemove: buttonRemove,
	}
}

func (item *SCTransferItem) Layout(gtx layout.Context, th *material.Theme, index int) layout.Dimensions {
	if item.buttonRemove.Clicked(gtx) {
		go item.onDelete(index)
	}

	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			item.scIdInput.Colors = theme.Current.InputColors
			return item.scIdInput.Layout(gtx, th, "SCID")
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			item.amountInput.Colors = theme.Current.InputColors
			return item.amountInput.Layout(gtx, th, lang.Translate("Amount"))
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			item.buttonRemove.Style.Colors = theme.Current.ButtonSecondaryColors
			return item.buttonRemove.Layout(gtx, th)
		}),
	)
}
