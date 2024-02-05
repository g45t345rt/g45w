package page_wallet

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/lang"
	page_settings "github.com/g45t345rt/g45w/pages/settings"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/wallet_manager"
)

type PageWalletInfo struct {
	isActive            bool
	headerPageAnimation *prefabs.PageHeaderAnimation
	list                *widget.List
	infoItems           []*page_settings.InfoListItem
}

var _ router.Page = &PageWalletInfo{}

func NewPageWalletInfo() *PageWalletInfo {

	list := new(widget.List)
	list.Axis = layout.Vertical

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_WALLET_INFO)
	return &PageWalletInfo{
		headerPageAnimation: headerPageAnimation,
		list:                list,
	}
}

func (p *PageWalletInfo) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Wallet Information") }

	wallet := wallet_manager.OpenedWallet

	addr := wallet.Memory.GetAddress().String()
	seed := wallet.Memory.GetSeed()
	hexSeed := wallet.Memory.Get_Keys().Secret.Text(16)

	infoItems := []*page_settings.InfoListItem{
		page_settings.NewInfoListItem("Address", addr, text.WrapGraphemes),     //@lang.Translate("Address")
		page_settings.NewInfoListItem("Seed", seed, text.WrapWords),            //@lang.Translate("Seed")
		page_settings.NewInfoListItem("Hex Seed", hexSeed, text.WrapGraphemes), //@lang.Translate("Hex Seed")
	}

	p.infoItems = infoItems
}

func (p *PageWalletInfo) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)

	// clear from memory
	p.infoItems = make([]*page_settings.InfoListItem, 0)
}

func (p *PageWalletInfo) IsActive() bool {
	return p.isActive
}

func (p *PageWalletInfo) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay
	return listStyle.Layout(gtx, len(p.infoItems), func(gtx layout.Context, index int) layout.Dimensions {
		return p.infoItems[index].Layout(gtx, th)
	})
}
