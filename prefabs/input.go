package prefabs

import (
	"time"

	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/lang"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Input struct {
	*components.Input
	holdPress      *components.HoldPress
	clipboardValue string
}

func NewInput() *Input {
	return &Input{
		Input:     components.NewInput(),
		holdPress: components.NewHoldPress(500 * time.Millisecond),
	}
}

func NewNumberInput() *Input {
	return &Input{
		Input:     components.NewNumberInput(),
		holdPress: components.NewHoldPress(500 * time.Millisecond),
	}
}

func NewPasswordInput() *Input {
	return &Input{
		Input:     components.NewPasswordInput(),
		holdPress: components.NewHoldPress(500 * time.Millisecond),
	}
}

func (i *Input) Layout(gtx layout.Context, th *material.Theme, hint string) layout.Dimensions {
	for _, e := range gtx.Events(i) {
		switch e := e.(type) {
		case clipboard.Event:
			i.clipboardValue = e.Text
		}
	}

	if i.holdPress.Triggered && !i.Editor.ReadOnly {
		clipboard.ReadOp{Tag: i}.Add(gtx.Ops)

		go func() {
			pasteIcon, _ := widget.NewIcon(icons.ContentContentPaste)
			copyIcon, _ := widget.NewIcon(icons.ContentContentCopy)
			selectIcon, _ := widget.NewIcon(icons.ContentSelectAll)
			clearIcon, _ := widget.NewIcon(icons.ContentClear)

			keyChan := listselect_modal.Instance.Open([]*listselect_modal.SelectListItem{
				listselect_modal.NewSelectListItem("paste",
					listselect_modal.NewItemText(pasteIcon, lang.Translate("Paste")).Layout,
				),
				listselect_modal.NewSelectListItem("copy",
					listselect_modal.NewItemText(copyIcon, lang.Translate("Copy")).Layout,
				),
				listselect_modal.NewSelectListItem("select_all",
					listselect_modal.NewItemText(selectIcon, lang.Translate("Select All")).Layout,
				),
				listselect_modal.NewSelectListItem("clear",
					listselect_modal.NewItemText(clearIcon, lang.Translate("Clear")).Layout,
				),
			})

			for key := range keyChan {
				switch key {
				case "paste":
					//i.newValue = i.clipboardValue
					i.Input.SetValue(i.clipboardValue)
					///i.Input.SetValue(i.clipboardValue)
					//i.Input.Editor.SetCaret(len(i.Input.Value()), len(i.Input.Value()))
				case "copy":
					text := i.Input.Editor.SelectedText()
					if text == "" {
						text = i.Input.Value()
					}
					app_instance.Window.WriteClipboard(text)
				case "select_all":
					i.Input.Editor.SetCaret(len(i.Input.Value()), 0)
					i.Input.Editor.Focus()
				case "clear":
					i.Input.SetValue("")
					//i.Input.SetValue("")
				}
			}

			app_instance.Window.Invalidate()
		}()
	}

	return i.holdPress.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return i.Input.Layout(gtx, th, hint)
	})
}
