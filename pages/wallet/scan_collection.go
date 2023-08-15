package page_wallet

import (
	"context"
	"fmt"
	"image"
	"strings"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/deroproject/derohe/rpc"
	"github.com/deroproject/derohe/walletapi"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/sc/g45_sc"
	"github.com/g45t345rt/g45w/theme"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageScanCollection struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	buttonFetchData              *components.Button
	txtSCID                      *prefabs.TextField
	scCollectionDetailsContainer *SCCollectionDetailsContainer

	list *widget.List
}

var _ router.Page = &PageScanCollection{}

func NewPageScanCollection() *PageScanCollection {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	checkIcon, _ := widget.NewIcon(icons.ActionSearch)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonFetchData := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        checkIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonFetchData.Label.Alignment = text.Middle
	buttonFetchData.Style.Font.Weight = font.Bold

	txtSCID := prefabs.NewTextField()

	list := new(widget.List)
	list.Axis = layout.Vertical
	scCollectionDetailsContainer := NewSCCollectionDetailsContainer()

	return &PageScanCollection{
		animationEnter:               animationEnter,
		animationLeave:               animationLeave,
		buttonFetchData:              buttonFetchData,
		txtSCID:                      txtSCID,
		scCollectionDetailsContainer: scCollectionDetailsContainer,

		list: list,
	}
}

func (p *PageScanCollection) IsActive() bool {
	return p.isActive
}

func (p *PageScanCollection) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Scan Collection") }
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil
	if !page_instance.header.IsHistory(PAGE_SCAN_COLLECTION) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageScanCollection) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageScanCollection) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonFetchData.Clicked() {
		p.scCollectionDetailsContainer.collection = nil
		p.buttonFetchData.SetLoading(true)
		scId, scType, scResult, err := p.submitForm()
		if err == nil {
			err = p.scCollectionDetailsContainer.Set(scId, scType, scResult)
		}
		p.buttonFetchData.SetLoading(false)

		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		}
	}

	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return p.txtSCID.Layout(gtx, th, "SCID", "Smart Contract ID")
		},
		func(gtx layout.Context) layout.Dimensions {
			if p.buttonFetchData.Loading {
				p.buttonFetchData.Text = lang.Translate("LOADING...")
			} else {
				p.buttonFetchData.Text = lang.Translate("FETCH DATA")
			}

			p.buttonFetchData.Style.Colors = theme.Current.ButtonPrimaryColors
			return p.buttonFetchData.Layout(gtx, th)
		},
	}

	if p.scCollectionDetailsContainer.collection != nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.scCollectionDetailsContainer.Layout(gtx, th)
		})
	}

	listStyle := material.List(th, p.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(30), Right: unit.Dp(30),
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageScanCollection) submitForm() (scId string, scType sc.SCType, result *rpc.GetSC_Result, err error) {
	scId = strings.TrimSpace(p.txtSCID.Value())
	if scId == "" {
		return scId, sc.UNKNOWN_TYPE, nil, fmt.Errorf("scid is empty")
	}

	err = walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
		SCID:      scId,
		Variables: true,
		Code:      true,
	}, &result)
	if err != nil {
		return scId, sc.UNKNOWN_TYPE, nil, err
	}

	if result.Code == "" {
		return scId, sc.UNKNOWN_TYPE, nil, fmt.Errorf("collection does not exists")
	}

	scType = sc.CheckType(result.Code)
	return scId, scType, result, nil
}

type SCCollectionDetailsContainer struct {
	scIdEditor        *widget.Editor
	nameEditor        *widget.Editor
	totalAssetsEditor *widget.Editor
	dateEditor        *widget.Editor
	collection        *g45_sc.G45_C

	list *widget.List
}

func NewSCCollectionDetailsContainer() *SCCollectionDetailsContainer {
	list := new(widget.List)
	list.Axis = layout.Vertical

	scIdEditor := new(widget.Editor)
	scIdEditor.WrapPolicy = text.WrapGraphemes
	scIdEditor.ReadOnly = true

	nameEditor := new(widget.Editor)
	nameEditor.WrapPolicy = text.WrapGraphemes
	nameEditor.ReadOnly = true

	totalAssetsEditor := new(widget.Editor)
	totalAssetsEditor.WrapPolicy = text.WrapGraphemes
	totalAssetsEditor.ReadOnly = true

	dateEditor := new(widget.Editor)
	dateEditor.WrapPolicy = text.WrapGraphemes
	dateEditor.ReadOnly = true

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAddToken := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      addIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonAddToken.Label.Alignment = text.Middle
	buttonAddToken.Style.Font.Weight = font.Bold

	return &SCCollectionDetailsContainer{
		scIdEditor:        scIdEditor,
		nameEditor:        nameEditor,
		totalAssetsEditor: totalAssetsEditor,
		dateEditor:        dateEditor,

		list: list,
	}
}

func (c *SCCollectionDetailsContainer) Set(scId string, scType sc.SCType, scResult *rpc.GetSC_Result) error {
	if scType != sc.G45_C_TYPE {
		return fmt.Errorf("not a valid G45_C smart contract")
	}

	collection := &g45_sc.G45_C{}
	err := collection.Parse(scId, scResult.VariableStringKeys)
	if err != nil {
		return err
	}

	metadata := g45_sc.CollectionMetadata{}
	err = metadata.Parse(collection.Metadata)
	if err != nil {
		return err
	}

	c.scIdEditor.SetText(collection.SCID)
	c.nameEditor.SetText(metadata.Name)
	c.totalAssetsEditor.SetText(fmt.Sprint(collection.AssetCount))

	date := time.Unix(int64(collection.Timestamp), 0)
	c.dateEditor.SetText(date.Format("2006-01-02 15:04:05"))

	c.collection = collection

	return nil
}

func (c *SCCollectionDetailsContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var widgets []layout.Widget

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return prefabs.Divider(gtx, 5)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("SCID"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, c.scIdEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Name"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, c.nameEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Asset Count"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, c.totalAssetsEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Created Date"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, c.dateEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
			)
		})
		c := r.Stop()

		paint.FillShape(
			gtx.Ops,
			theme.Current.ListBgColor,
			clip.UniformRRect(
				image.Rectangle{Max: dims.Size},
				gtx.Dp(10),
			).Op(gtx.Ops),
		)

		c.Add(gtx.Ops)
		return dims
	})

	/*widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		c.buttonAddToken.Text = lang.Translate("ADD TOKEN")
		c.buttonAddToken.Style.Colors = theme.Current.ButtonPrimaryColors
		return c.buttonAddToken.Layout(gtx, th)
	})*/

	listStyle := material.List(th, c.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(0), Right: unit.Dp(0),
		}.Layout(gtx, widgets[index])
	})
}
