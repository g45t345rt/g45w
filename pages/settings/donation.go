package page_settings

import (
	"context"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_icons"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/build_tx_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type PageDonation struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonAmount1     *components.Button
	buttonAmount2     *components.Button
	buttonAmount3     *components.Button
	buttonDonate      *components.Button
	txtAmount         *prefabs.TextField
	anonymousDonation *widget.Bool

	donationResult DonationResult
	infoRows       []*prefabs.InfoRow

	list *widget.List
}

var _ router.Page = &PageDonation{}

func NewPageDonation() *PageDonation {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	buttonAmount1 := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonAmount1.Label.Alignment = text.Middle
	buttonAmount1.Style.Font.Weight = font.Bold

	buttonAmount2 := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonAmount2.Label.Alignment = text.Middle
	buttonAmount2.Style.Font.Weight = font.Bold

	buttonAmount3 := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Border: widget.Border{
			Color:        color.NRGBA{R: 0, G: 0, B: 0, A: 255},
			Width:        unit.Dp(2),
			CornerRadius: unit.Dp(5),
		},
	})
	buttonAmount3.Label.Alignment = text.Middle
	buttonAmount3.Style.Font.Weight = font.Bold

	donateIcon, _ := widget.NewIcon(app_icons.Donation)
	buttonDonate := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		TextSize:  unit.Sp(14),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
		Icon:      donateIcon,
		IconGap:   unit.Dp(10),
	})
	buttonDonate.Label.Alignment = text.Middle
	buttonDonate.Style.Font.Weight = font.Bold

	txtAmount := prefabs.NewTextField()
	txtAmount.Input.Editor.InputHint = key.HintNumeric
	txtAmount.Input.TextSize = unit.Sp(20)
	txtAmount.Input.FontWeight = font.Bold

	var infoRows []*prefabs.InfoRow
	for i := 0; i < 9; i++ {
		infoRows = append(infoRows, prefabs.NewInfoRow())
	}

	return &PageDonation{
		animationEnter: animationEnter,
		animationLeave: animationLeave,

		txtAmount:         txtAmount,
		buttonDonate:      buttonDonate,
		buttonAmount1:     buttonAmount1,
		buttonAmount2:     buttonAmount2,
		buttonAmount3:     buttonAmount3,
		anonymousDonation: new(widget.Bool),
		infoRows:          infoRows,

		list: list,
	}
}

func (p *PageDonation) IsActive() bool {
	return p.isActive
}

func (p *PageDonation) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Donation") }
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil

	if !page_instance.header.IsHistory(PAGE_DONATION) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageDonation) Leave() {
	p.animationEnter.Reset()
	p.animationLeave.Start()
}

type DonationResult struct {
	TotalDonated             uint64
	TotalDonations           uint64
	HighestDonation          uint64
	HighestDonationAddr      string
	HighestDonationTimestamp uint64
	LastDonation             uint64
	LastDonationAddr         string
	LastDonationTimestamp    uint64
	TotalAnonymouslyDonated  uint64
}

func (p *PageDonation) Load() error {
	var result rpc.GetSC_Result
	err := walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
		SCID:      "",
		Code:      false,
		Variables: false,
		KeysString: []string{
			"totalDonated", "totalDonations",
			"highestDonation", "highestDonationAddr", "highestDonationTimestamp",
			"lastDonation", "lastDonationAddr", "lastDonationTimestamp",
			"d_", // d_ = total donated anonymously or d_{signer} = total value donated by one specific addr
		},
	}, &result)
	if err != nil {
		return err
	}

	var donationResult DonationResult
	donationResult.TotalDonated, _ = strconv.ParseUint(result.ValuesString[0], 10, 64)
	donationResult.TotalDonations, _ = strconv.ParseUint(result.ValuesString[1], 10, 64)

	donationResult.HighestDonation, _ = strconv.ParseUint(result.ValuesString[2], 10, 64)
	donationResult.HighestDonationAddr, _ = utils.DecodeString(result.ValuesString[3])
	donationResult.HighestDonationTimestamp, _ = strconv.ParseUint(result.ValuesString[4], 10, 64)

	donationResult.LastDonation, _ = strconv.ParseUint(result.ValuesString[5], 10, 64)
	donationResult.LastDonationAddr, _ = utils.DecodeString(result.ValuesString[6])
	donationResult.LastDonationTimestamp, _ = strconv.ParseUint(result.ValuesString[7], 10, 64)

	donationResult.TotalAnonymouslyDonated, _ = strconv.ParseUint(result.ValuesString[8], 10, 64)

	p.donationResult = donationResult

	return nil
}

func (p *PageDonation) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonAmount1.Clicked() {
		p.txtAmount.SetValue("1")
	}

	if p.buttonAmount2.Clicked() {
		p.txtAmount.SetValue("10")
	}

	if p.buttonAmount3.Clicked() {
		p.txtAmount.SetValue("25")
	}

	if p.buttonDonate.Clicked() {
		go func() {
			donate := func() error {
				wallet := wallet_manager.OpenedWallet
				if wallet == nil {
					return fmt.Errorf("wallet is not opened")
				}

				amount := utils.ShiftNumber{Decimals: 5}
				err := amount.Parse(p.txtAmount.Value())
				if err != nil {
					return err
				}

				ringsize := uint64(2)
				if p.anonymousDonation.Value {
					ringsize = 16
				}

				build_tx_modal.Instance.OpenWithRandomAddr(crypto.ZEROHASH, func(randomAddr string, open func(txPayload build_tx_modal.TxPayload)) {
					open(build_tx_modal.TxPayload{
						Transfers: []rpc.Transfer{
							{SCID: crypto.ZEROHASH, Destination: randomAddr, Burn: amount.Number},
						},
						Ringsize: ringsize,
						SCArgs: rpc.Arguments{
							{Name: rpc.SCACTION, DataType: rpc.DataUint64, Value: uint64(rpc.SC_CALL)},
							{Name: rpc.SCID, DataType: rpc.DataHash, Value: crypto.HashHexToHash("")},
							{Name: "entrypoint", DataType: rpc.DataString, Value: "Donate"},
						},
					})
				})

				return nil
			}

			err := donate()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, 0)
			}
		}()
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), lang.Translate("Hello, I'm the anonymous developer, g45t345rt! Please never feel obligated to donate; only consider doing so if you can afford it. Thanks!"))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			txt := lang.Translate("As of now, there have been {0} donations with a total of {1} DERO.")
			txt = strings.Replace(txt, "{0}", fmt.Sprint(p.donationResult.TotalDonations), -1)
			totalDonated := globals.FormatMoney(p.donationResult.TotalDonated)
			txt = strings.Replace(txt, "{1}", totalDonated, -1)
			lbl := material.Label(th, unit.Sp(16), txt)
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					p.buttonAmount1.Text = "1 DERO"
					p.buttonAmount1.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonAmount1.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					p.buttonAmount2.Text = "10 DERO"
					p.buttonAmount2.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonAmount2.Layout(gtx, th)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					p.buttonAmount3.Text = "25 DERO"
					p.buttonAmount3.Style.Colors = theme.Current.ButtonSecondaryColors
					return p.buttonAmount3.Layout(gtx, th)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return p.txtAmount.Layout(gtx, th, lang.Translate("Amount"), "0.00000")
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					s := material.Switch(th, p.anonymousDonation, "")
					s.Color = theme.Current.SwitchColors
					return s.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					txt := lang.Translate("Anonymous donation? ({})")
					if p.anonymousDonation.Value {
						txt = strings.Replace(txt, "{}", lang.Translate("YES"), -1)
					} else {
						txt = strings.Replace(txt, "{}", lang.Translate("NO"), -1)
					}

					lbl := material.Label(th, unit.Sp(16), txt)
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			p.buttonDonate.Text = lang.Translate("DONATE")
			p.buttonDonate.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonDonate.Layout(gtx, th)
		},
		func(gtx layout.Context) layout.Dimensions {
			return prefabs.Divider(gtx, unit.Dp(5))
		},
		func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(14), lang.Translate("This donation system uses an on-chain smart contract to keep donations transparent."))
			lbl.Color = theme.Current.TextMuteColor
			return lbl.Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Donation Stats"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := globals.FormatMoney(p.donationResult.TotalDonated) + " DERO"
					return p.infoRows[0].Layout(gtx, th, lang.Translate("Total Donated"), value)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := fmt.Sprint(p.donationResult.TotalDonations)
					return p.infoRows[1].Layout(gtx, th, lang.Translate("Donation Count"), value)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := globals.FormatMoney(p.donationResult.TotalAnonymouslyDonated) + " DERO"
					return p.infoRows[2].Layout(gtx, th, lang.Translate("Anon Donated"), value)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				// Highest Donation
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Highest Donation"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := globals.FormatMoney(p.donationResult.HighestDonation) + " DERO"
					return p.infoRows[3].Layout(gtx, th, lang.Translate("Amount"), value)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := utils.ReduceAddr(p.donationResult.HighestDonationAddr)
					if value == "" {
						value = "?"
					}
					return p.infoRows[4].Layout(gtx, th, lang.Translate("Address"), value)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := "?"
					if p.donationResult.HighestDonationTimestamp > 0 {
						value = lang.TimeAgo(time.Unix(int64(p.donationResult.HighestDonationTimestamp), 0))
					}

					return p.infoRows[5].Layout(gtx, th, lang.Translate("Time"), value)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				// Last Donation
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Last Donation"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := globals.FormatMoney(p.donationResult.LastDonation) + " DERO"
					return p.infoRows[6].Layout(gtx, th, lang.Translate("Amount"), value)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := utils.ReduceAddr(p.donationResult.LastDonationAddr)
					if value == "" {
						value = "?"
					}
					return p.infoRows[7].Layout(gtx, th, lang.Translate("Address"), value)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					value := "?"
					if p.donationResult.LastDonationTimestamp > 0 {
						value = lang.TimeAgo(time.Unix(int64(p.donationResult.LastDonationTimestamp), 0))
					}

					return p.infoRows[8].Layout(gtx, th, lang.Translate("Time"), value)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
		},
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	if p.txtAmount.Input.Clickable.Clicked() {
		p.list.ScrollTo(3)
	}

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}
