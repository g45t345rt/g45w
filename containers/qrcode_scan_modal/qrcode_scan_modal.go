package qrcode_scan_modal

import (
	"errors"
	"math"
	"time"

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
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type CameraQRScanModal struct {
	cameraImage  *components.Image
	buttonCancel *components.Button
	buttonOk     *components.Button
	buttonRetry  *components.Button

	scanned           bool
	scanning          bool
	value             string
	err               error
	send              bool
	cameraOrientation int
	closing           bool

	Modal *components.Modal
}

var Instance *CameraQRScanModal

func LoadInstance() {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: false,
		CloseOnInsideClick:  false,
		KeepClickableArea:   true,
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

	Instance = &CameraQRScanModal{
		Modal:        modal,
		cameraImage:  cameraImage,
		buttonCancel: buttonCancel,
		buttonRetry:  buttonRetry,
		buttonOk:     buttonOk,
	}

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 2,
		Layout:    Instance.layout,
	})
}

func (w *CameraQRScanModal) Value() (bool, string) {
	if w.send {
		w.send = false
		return true, w.value
	}

	return false, ""
}

func (w *CameraQRScanModal) scan() {
	cameraId := ""
	ids, err := camera.GetIdList()
	if err != nil {
		w.err = err
		return
	}

	for _, id := range ids {
		lensFacing, err := camera.GetLensFacing(id)
		if err != nil {
			w.err = err
			return
		}

		if lensFacing == camera.LensFacingBack {
			cameraId = id
			break
		}
	}

	if cameraId == "" {
		w.err = errors.New("no back camera")
		return
	}

	w.cameraOrientation, err = camera.GetSensorOrientation(cameraId)
	if err != nil {
		w.err = err
		return
	}

	err = camera.OpenFeed(cameraId, 256, 256)
	if err != nil {
		w.err = err
		return
	}

	w.err = nil
	w.scanned = false
	w.scanning = true
	go func() {
		for imageResult := range camera.GetFeed() {
			if imageResult.Err != nil {
				w.err = imageResult.Err
				break
			}

			img := imageResult.Image
			w.cameraImage.Src = paint.NewImageOp(img)

			bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
			qrReader := qrcode.NewQRCodeReader()
			result, err := qrReader.Decode(bmp, nil)

			if err == nil {
				w.scanned = true
				w.value = result.String()
				camera.CloseFeed()
				app_instance.Window.Invalidate()
				break
			}

			app_instance.Window.Invalidate()
		}

		w.scanning = false
	}()
}

func (w *CameraQRScanModal) Open() {
	w.scan()
	w.Modal.SetVisible(true)
}

func (w *CameraQRScanModal) close() {
	if w.closing {
		return
	}

	// make sure you don't call closefeed multiple times by spamming button - it crash otherwise
	w.closing = true
	camera.CloseFeed()
	w.Modal.SetVisible(false)
	app_instance.Window.Invalidate()

	time.AfterFunc(1*time.Second, func() {
		w.closing = false
	})
}

func (w *CameraQRScanModal) layout(gtx layout.Context, th *material.Theme) {
	if w.buttonCancel.Clicked(gtx) {
		go w.close()
	}

	if w.buttonRetry.Clicked(gtx) {
		go w.scan()
	}

	if w.buttonOk.Clicked(gtx) {
		w.send = true
		w.Modal.SetVisible(false)
	}

	w.Modal.Style.Colors = theme.Current.ModalColors
	w.Modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			var childs []layout.FlexChild
			if w.scanning {
				childs = append(childs,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						w.cameraImage.Transform = func(dims layout.Dimensions, trans f32.Affine2D) f32.Affine2D {
							pt := dims.Size.Div(2)
							origin := f32.Pt(float32(pt.X), float32(pt.Y))
							rotate := float32(w.cameraOrientation) * (math.Pi / 180)
							return trans.Rotate(origin, rotate)
						}

						return w.cameraImage.Layout(gtx)
					}),
				)
			}

			if w.err != nil {
				childs = append(childs,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						lbl := material.Label(th, unit.Sp(14), w.err.Error())
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								w.buttonCancel.Text = lang.Translate("CANCEL")
								w.buttonCancel.Style.Colors = theme.Current.ButtonPrimaryColors
								return w.buttonCancel.Layout(gtx, th)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								w.buttonRetry.Text = lang.Translate("RETRY")
								w.buttonRetry.Style.Colors = theme.Current.ButtonPrimaryColors
								return w.buttonRetry.Layout(gtx, th)
							}),
						)
					}),
				)
			}

			if w.scanning {
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
