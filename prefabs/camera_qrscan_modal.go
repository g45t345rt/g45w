package prefabs

import (
	"math"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/camera"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/theme"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type CameraQRScanModal struct {
	cameraImage  *components.Image
	buttonCancel *components.Button
	buttonOk     *components.Button
	buttonRetry  *components.Button

	scanned  bool
	scanning bool
	value    string
	err      error
	send     bool

	Modal *components.Modal
}

func NewCameraQRScanModal() *CameraQRScanModal {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: false,
		CloseOnInsideClick:  false,
		Direction:           layout.Center,
		Rounded:             components.UniformRounded(unit.Dp(10)),
		Inset:               layout.UniformInset(25),
		Animation:           components.NewModalAnimationScaleBounce(),
	})

	buttonCancel := components.NewButton(components.ButtonStyle{
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Inset:    layout.UniformInset(unit.Dp(10)),
	})
	buttonCancel.Style.Font.Weight = font.Bold

	buttonOk := components.NewButton(components.ButtonStyle{
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Inset:    layout.UniformInset(unit.Dp(10)),
	})
	buttonOk.Style.Font.Weight = font.Bold

	buttonRetry := components.NewButton(components.ButtonStyle{
		Rounded:  components.UniformRounded(unit.Dp(5)),
		TextSize: unit.Sp(14),
		Inset:    layout.UniformInset(unit.Dp(10)),
	})
	buttonRetry.Style.Font.Weight = font.Bold

	cameraImage := &components.Image{
		Fit:     components.Contain,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}

	return &CameraQRScanModal{
		Modal:        modal,
		cameraImage:  cameraImage,
		buttonCancel: buttonCancel,
		buttonRetry:  buttonRetry,
		buttonOk:     buttonOk,
	}
}

func (w *CameraQRScanModal) Value() (bool, string) {
	if w.send {
		w.send = false
		return true, w.value
	}

	return false, ""
}

func (w *CameraQRScanModal) scan() {
	err := camera.Open()
	if err != nil {
		w.err = err
		return
	}

	w.err = nil
	w.scanned = false
	w.scanning = true
	go func() {
		for imageResult := range camera.GetImage() {
			if imageResult.Err != nil {
				continue
			}

			img := imageResult.Image
			w.cameraImage.Src = paint.NewImageOp(img)

			bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
			qrReader := qrcode.NewQRCodeReader()
			result, err := qrReader.Decode(bmp, nil)

			if err == nil {
				w.scanned = true
				w.value = result.String()
				w.scanning = false
				camera.Close()
				app_instance.Window.Invalidate()
				break
			}

			app_instance.Window.Invalidate()
		}
	}()
}

func (w *CameraQRScanModal) Show() {
	go w.scan()
	w.Modal.SetVisible(true)
}

func (w *CameraQRScanModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if w.buttonCancel.Clicked() {
		go camera.Close()
		w.Modal.SetVisible(false)
	}

	if w.buttonRetry.Clicked() {
		go w.scan()
	}

	if w.buttonOk.Clicked() {
		w.send = true
		w.Modal.SetVisible(false)
	}

	w.Modal.Style.Colors = theme.Current.ModalColors
	return w.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var childs []layout.FlexChild
			if w.scanning {
				childs = append(childs,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						w.cameraImage.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
							pt := dims.Size.Div(2)
							origin := f32.Pt(float32(pt.X), float32(pt.Y))
							rotate := float32(90 * (math.Pi / 180))
							return trans.Rotate(origin, rotate)
						}

						return w.cameraImage.Layout(gtx)
					}),
				)
			}

			if w.err != nil {
				childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(14), w.err.Error())
					return lbl.Layout(gtx)
				}))
			}

			if w.err != nil || w.scanning {
				childs = append(childs,
					layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							w.buttonCancel.Text = lang.Translate("CANCEL")
							w.buttonCancel.Style.Colors = theme.Current.ButtonPrimaryColors
							return w.buttonCancel.Layout(gtx, th)
						}))
					}),
				)
			}

			if w.scanned {
				childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(14), w.value)
							lbl.Alignment = text.Middle
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									w.buttonRetry.Text = lang.Translate("RETRY")
									w.buttonRetry.Style.Colors = theme.Current.ButtonPrimaryColors
									return w.buttonRetry.Layout(gtx, th)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									w.buttonOk.Text = lang.Translate("OK")
									w.buttonOk.Style.Colors = theme.Current.ButtonPrimaryColors
									return w.buttonOk.Layout(gtx, th)
								}),
							)
						}),
					)
				}))
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
		})
	})
}
