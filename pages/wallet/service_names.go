package page_wallet

import (
	"database/sql"
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageServiceNames struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation
	buttonRegister      *components.Button
	txtName             *prefabs.TextField

	entries []wallet_manager.Entry

	list *widget.List
}

var _ router.Page = &PageServiceNames{}

var SERVICE_NAME_SCID = crypto.HashHexToHash("0000000000000000000000000000000000000000000000000000000000000001")

func NewPageServiceNames() *PageServiceNames {

	addIcon, _ := widget.NewIcon(icons.ContentCreate)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonRegister := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        addIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonRegister.Label.Alignment = text.Middle
	buttonRegister.Style.Font.Weight = font.Bold

	txtName := prefabs.NewTextField()

	list := new(widget.List)
	list.Axis = layout.Vertical

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_SERVICE_NAMES)
	return &PageServiceNames{
		headerPageAnimation: headerPageAnimation,
		list:                list,
		buttonRegister:      buttonRegister,
		txtName:             txtName,
		entries:             make([]wallet_manager.Entry, 0),
	}
}

func (p *PageServiceNames) IsActive() bool {
	return p.isActive
}

func (p *PageServiceNames) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string {
		return lang.Translate("Service Names")
	}

	page_instance.header.Subtitle = nil
	p.Load()
}

func (p *PageServiceNames) Load() {
	p.entries = make([]wallet_manager.Entry, 0)
	wallet := wallet_manager.OpenedWallet

	p.entries = wallet.GetEntries(&crypto.ZEROHASH, wallet_manager.GetEntriesParams{
		SC_CALL: &wallet_manager.SCCallParams{
			SCID:       sql.NullString{String: SERVICE_NAME_SCID.String(), Valid: true},
			Entrypoint: sql.NullString{String: "Register", Valid: true},
		},
	})
}

func (p *PageServiceNames) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageServiceNames) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonRegister.Clicked(gtx) {
		go func() {
			err := p.submitForm()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			}
		}()
	}

	widgets := []layout.Widget{}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Label(th, unit.Sp(16), lang.Translate("On Dero, you can have multiple usernames for your wallet. You can use them instead of your Dero address for receiving payments."))
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return p.txtName.Layout(gtx, th, lang.Translate("Name"), "")
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		p.buttonRegister.Style.Colors = theme.Current.ButtonPrimaryColors
		p.buttonRegister.Text = lang.Translate("REGISTER")
		return p.buttonRegister.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return prefabs.Divider(gtx, unit.Dp(5))
	})

	if len(p.entries) > 0 {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			var childs []layout.FlexChild

			for i := range p.entries {
				idx := i
				childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					entry := p.entries[idx]
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							name := ""
							for _, arg := range entry.SCDATA {
								if arg.Name == "name" {
									value, ok := arg.Value.(string)
									if ok {
										name = value
									}
								}
							}

							lbl := material.Label(th, unit.Sp(16), name)
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), entry.Time.Format("2006-01-02"))
							lbl.Color = theme.Current.TextMuteColor
							return lbl.Layout(gtx)
						}),
					)
				}))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
		})
	} else {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("No registered names or the wallet is not synced properly. Try cleaning the wallet in settings page."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		})
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
	})

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageServiceNames) submitForm() error {
	wallet := wallet_manager.OpenedWallet
	name := p.txtName.Value()

	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(name) < 6 {
		return fmt.Errorf("name must be at least 6 characters")
	}

	addr, err := wallet.Memory.NameToAddress(name)
	if err != nil {
		if !utils.IsErrLeafNotFound(err) {
			return err
		}
	}

	if addr != "" {
		return fmt.Errorf("name already taken by [%s]", utils.ReduceAddr(addr))
	}

	build_tx_modal.Instance.Open(build_tx_modal.TxPayload{
		Ringsize: 2,
		SCArgs: rpc.Arguments{
			{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
			{Name: rpc.SCID, DataType: rpc.DataHash, Value: SERVICE_NAME_SCID},
			{Name: "entrypoint", DataType: rpc.DataString, Value: "Register"},
			{Name: "name", DataType: rpc.DataString, Value: name},
		},
	})

	return nil
}
