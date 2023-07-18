package page_wallet

import (
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/lang"
	page_settings "github.com/g45t345rt/g45w/pages/settings"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageWalletInfo struct {
	isActive       bool
	animationLeave *animation.Animation
	animationEnter *animation.Animation
	list           *widget.List
	infoItems      []*page_settings.InfoListItem
}

var _ router.Page = &PageWalletInfo{}

func NewPageWalletInfo() *PageWalletInfo {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	return &PageWalletInfo{
		animationEnter: animationEnter,
		animationLeave: animationLeave,
		list:           list,
	}
}

func (p *PageWalletInfo) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Wallet Information"))

	if !page_instance.header.IsHistory(PAGE_WALLET_INFO) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	wallet := wallet_manager.OpenedWallet

	addr := wallet.Memory.GetAddress().String()
	seed := wallet.Memory.GetSeed()
	hexSeed := wallet.Memory.Get_Keys().Secret.Text(16)

	infoItems := []*page_settings.InfoListItem{
		page_settings.NewInfoListItem("Address", addr),     //@lang.Translate("Address")
		page_settings.NewInfoListItem("Seed", seed),        //@lang.Translate("Seed")
		page_settings.NewInfoListItem("Hex Seed", hexSeed), //@lang.Translate("Hex Seed")
	}

	p.infoItems = infoItems
}

func (p *PageWalletInfo) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

func (p *PageWalletInfo) IsActive() bool {
	return p.isActive
}

func (p *PageWalletInfo) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
			p.infoItems = make([]*page_settings.InfoListItem, 0)
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
