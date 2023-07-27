package build_tx_modal

import (
	"fmt"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/globals"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/transaction"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/recent_txs_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type TxPayload struct {
	Transfers   []rpc.Transfer
	Ringsize    uint64
	SCData      rpc.Arguments
	Description string
}

type BuildTxModal struct {
	modal               *components.Modal
	buttonSend          *components.Button
	editorError         *widget.Editor
	modalWalletPassword *prefabs.PasswordModal
	loadingIcon         *widget.Icon
	animationLoading    *animation.Animation

	building bool
	buildTx  *transaction.Transaction
	buildErr error
	txSent   bool

	txPayload TxPayload
}

var Instance *BuildTxModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.N,
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Inset:               layout.UniformInset(unit.Dp(10)),
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Animation:           components.NewModalAnimationDown(),
		Backdrop:            components.NewModalBackground(),
	})

	sendIcon, _ := widget.NewIcon(icons.ContentSend)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonSend := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            sendIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		LoadingIcon:     loadingIcon,
	})
	buttonSend.Label.Alignment = text.Middle
	buttonSend.Style.Font.Weight = font.Bold

	editorError := new(widget.Editor)
	editorError.WrapPolicy = text.WrapGraphemes
	editorError.ReadOnly = true

	modalWalletPassword := prefabs.NewPasswordModal()

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 3,
		Layout: func(gtx layout.Context, th *material.Theme) {
			modalWalletPassword.Layout(gtx, th)
		},
	})

	animationLoading := animation.NewAnimation(false,
		gween.NewSequence(
			gween.New(0, 1, 1, ease.Linear),
		),
	)
	animationLoading.Sequence.SetLoop(-1)

	Instance = &BuildTxModal{
		modal:               modal,
		buttonSend:          buttonSend,
		editorError:         editorError,
		modalWalletPassword: modalWalletPassword,
		loadingIcon:         loadingIcon,
		animationLoading:    animationLoading,
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout:    Instance.layout,
	})
}

func (b *BuildTxModal) Open(txPayload TxPayload) {
	b.txSent = false
	b.txPayload = txPayload

	b.modal.SetVisible(true)
	b.animationLoading.Reset().Start()
	b.building = true
	wallet := wallet_manager.OpenedWallet
	tx, err := wallet.Memory.TransferPayload0(b.txPayload.Transfers, b.txPayload.Ringsize, false, b.txPayload.SCData, 0, false)
	b.buildErr = err
	if err != nil {
		b.editorError.SetText(err.Error())
	}

	b.building = false
	b.animationLoading.Pause()
	b.buildTx = tx
}

func (b *BuildTxModal) TxSent() bool {
	if b.txSent {
		b.txSent = false
		return true
	}

	return false
}

func (b *BuildTxModal) sendTx() error {
	b.buttonSend.SetLoading(true)
	wallet := wallet_manager.OpenedWallet
	tx := b.buildTx
	err := wallet.InsertOutgoingTx(tx, b.txPayload.Description)
	if err != nil {
		b.buttonSend.SetLoading(false)
		return err
	}

	err = wallet.Memory.SendTransaction(tx)
	if err != nil {
		b.buttonSend.SetLoading(false)
		return err
	}

	b.buttonSend.SetLoading(false)
	b.modal.SetVisible(false)
	recent_txs_modal.Instance.SetVisible(true)
	b.txSent = true
	return nil
}

func (b *BuildTxModal) layout(gtx layout.Context, th *material.Theme) {
	wallet := wallet_manager.OpenedWallet

	b.buttonSend.Text = lang.Translate("SEND TRANSACTION")
	if b.buttonSend.Clicked() {
		b.modalWalletPassword.Modal.SetVisible(true)
	}

	submitted, password := b.modalWalletPassword.Input.Submitted()
	if submitted {
		validPassword := wallet.Memory.Check_Password(password)

		if !validPassword {
			b.modalWalletPassword.StartWrongPassAnimation()
		} else {
			err := b.sendTx()
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			} else {
				b.modalWalletPassword.Modal.SetVisible(false)
			}
		}
	}

	b.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(15), Bottom: unit.Dp(15),
			Left: unit.Dp(15), Right: unit.Dp(15),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var childrens []layout.FlexChild

			if b.building {
				childrens = append(childrens,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								lbl := material.Label(th, unit.Sp(20), lang.Translate("Building transaction..."))
								lbl.Font.Weight = font.Bold
								return lbl.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								r := op.Record(gtx.Ops)
								dims := b.loadingIcon.Layout(gtx, color.NRGBA{A: 255})
								c := r.Stop()

								{
									gtx.Constraints.Min = dims.Size

									state := b.animationLoading.Update(gtx)
									if state.Active {
										defer animation.TransformRotate(gtx, state.Value).Push(gtx.Ops).Pop()
									}
								}

								c.Add(gtx.Ops)
								return dims
							}),
						)
					}))
			} else {
				if b.buildErr != nil {
					childrens = append(childrens,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("Error"))
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							editor := material.Editor(th, b.editorError, "")
							return editor.Layout(gtx)
						}))
				} else {
					childrens = append(childrens,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("TX fees"))
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							fees := globals.FormatMoney(b.buildTx.Fees())
							lbl := material.Label(th, unit.Sp(16), fmt.Sprintf("%s DERO", fees))
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("New balance"))
							lbl.Font.Weight = font.Bold
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							balance, _ := wallet.Memory.Get_Balance()
							newBalance := balance
							fees := b.buildTx.Fees()
							newBalance -= fees

							totalTransfer := uint64(0)
							for _, transfer := range b.txPayload.Transfers {
								if transfer.SCID.IsZero() {
									totalTransfer += transfer.Amount
								}
							}
							newBalance -= totalTransfer

							status := fmt.Sprintf("%s - %s - %s",
								globals.FormatMoney(balance),
								globals.FormatMoney(totalTransfer),
								globals.FormatMoney(fees),
							)

							return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(16), fmt.Sprintf("%s DERO", globals.FormatMoney(newBalance)))
									return lbl.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(12), status)
									return lbl.Layout(gtx)
								}),
							)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return b.buttonSend.Layout(gtx, th)
						}))
				}
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childrens...)
		})
	})
}
