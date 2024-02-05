package page_node

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/glue/rwc"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/app_db"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageAddNodeForm struct {
	isActive            bool
	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonAdd   *components.Button
	txtEndpoint *prefabs.TextField
	txtName     *prefabs.TextField

	list *widget.List
}

var _ router.Page = &PageAddNodeForm{}

func NewPageAddNodeForm() *PageAddNodeForm {
	list := new(widget.List)
	list.Axis = layout.Vertical

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonAdd := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        addIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonAdd.Label.Alignment = text.Middle
	buttonAdd.Style.Font.Weight = font.Bold

	txtName := prefabs.NewTextField()
	txtEndpoint := prefabs.NewTextField()

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_ADD_NODE_FORM)
	return &PageAddNodeForm{
		headerPageAnimation: headerPageAnimation,

		buttonAdd:   buttonAdd,
		txtName:     txtName,
		txtEndpoint: txtEndpoint,

		list: list,
	}
}

func (p *PageAddNodeForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddNodeForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)
	page_instance.header.Title = func() string { return lang.Translate("Add Node") }
}

func (p *PageAddNodeForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)

}

func (p *PageAddNodeForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonAdd.Clicked(gtx) {
		p.submitForm(gtx)
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("Here, you can add your own remote node. The endpoint connection must be a WebSocket connection, starting with ws:// or wss:// for TLS connection."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtName.Layout(gtx, th, lang.Translate("Name"), "Dero NFTs")
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtEndpoint.Layout(gtx, th, lang.Translate("Endpoint"), "wss://node.deronfts.com/ws")
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonAdd.Text = lang.Translate("ADD NODE")
			p.buttonAdd.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonAdd.Layout(gtx, th)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtName.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	if p.txtEndpoint.Input.Clickable.Clicked(gtx) {
		p.list.ScrollTo(0)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddNodeForm) submitForm(gtx layout.Context) {
	p.buttonAdd.SetLoading(true)
	go func() {
		setError := func(err error) {
			p.buttonAdd.SetLoading(false)
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
		}

		txtName := p.txtName.Editor()
		txtEndpoint := p.txtEndpoint.Editor()

		if txtName.Text() == "" {
			setError(fmt.Errorf("enter name"))
			return
		}

		if txtEndpoint.Text() == "" {
			setError(fmt.Errorf("enter endpoint"))
			return
		}

		_, err := TestConnect(txtEndpoint.Text())
		if err != nil {
			setError(err)
			return
		}

		err = app_db.InsertNodeConnection(app_db.NodeConnection{
			Name:     txtName.Text(),
			Endpoint: txtEndpoint.Text(),
		})
		if err != nil {
			setError(err)
			return
		}

		p.buttonAdd.SetLoading(false)
		notification_modal.Open(notification_modal.Params{
			Type:       notification_modal.SUCCESS,
			Title:      lang.Translate("Success"),
			Text:       lang.Translate("New noded added."),
			CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
		})
		page_instance.header.GoBack()
	}()
}

func TestConnect(endpoint string) (info rpc.GetInfo_Result, err error) {
	client := walletapi.Client{}
	ws, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return
	}

	client.WS = ws
	input_output := rwc.New(client.WS)
	client.RPC = jrpc2.NewClient(channel.RawJSON(input_output, input_output), &jrpc2.ClientOptions{})
	defer client.WS.Close()

	var result string
	err = client.Call("DERO.Echo", []string{"hello"}, &result)
	if err != nil {
		return
	}

	err = client.Call("DERO.GetInfo", nil, &info)
	if err != nil {
		return
	}

	if globals.IsMainnet() && info.Testnet {
		err = fmt.Errorf("this is not a Mainnet node")
	} else if !globals.IsMainnet() && !info.Testnet {
		err = fmt.Errorf("this is not a Testnet node")
	}

	return
}
