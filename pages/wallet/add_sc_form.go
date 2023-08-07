package page_wallet

import (
	"context"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"strings"

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
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/sc/g45_sc"
	"github.com/g45t345rt/g45w/sc/unknown_sc"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/text/language"
)

type PageAddSCForm struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	scDetailsContainer *SCDetailsContainer

	buttonFetchData *components.Button
	txtSCID         *components.TextField

	list *widget.List
}

var _ router.Page = &PageAddSCForm{}

func NewPageAddSCForm() *PageAddSCForm {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(-1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, -1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	checkIcon, _ := widget.NewIcon(icons.ActionSearch)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonFetchData := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            checkIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
		LoadingIcon:     loadingIcon,
	})
	buttonFetchData.Label.Alignment = text.Middle
	buttonFetchData.Style.Font.Weight = font.Bold

	txtSCID := components.NewTextField()

	return &PageAddSCForm{
		animationEnter:  animationEnter,
		animationLeave:  animationLeave,
		txtSCID:         txtSCID,
		buttonFetchData: buttonFetchData,
		list:            list,
	}
}

func (p *PageAddSCForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddSCForm) Enter() {
	p.isActive = true
	page_instance.header.SetTitle(lang.Translate("Add Token"))
	page_instance.header.Subtitle = nil
	page_instance.header.ButtonRight = nil
	if !page_instance.header.IsHistory(PAGE_ADD_SC_FORM) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}
}

func (p *PageAddSCForm) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageAddSCForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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
		p.scDetailsContainer = nil
		scId := p.txtSCID.Value()
		p.buttonFetchData.SetLoading(true)
		scType, scResult, err := p.submitForm()
		p.buttonFetchData.SetLoading(false)
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			p.scDetailsContainer = NewSCDetailsContainer(scId, scType, scResult)
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

			return p.buttonFetchData.Layout(gtx, th)
		},
	}

	if p.scDetailsContainer != nil {
		widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
			return p.scDetailsContainer.Layout(gtx, th)
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

func (p *PageAddSCForm) submitForm() (scType sc.SCType, result *rpc.GetSC_Result, err error) {
	scId := strings.TrimSpace(p.txtSCID.Value())
	if scId == "" {
		return sc.UNKNOWN_TYPE, nil, fmt.Errorf("scid is empty")
	}

	err = walletapi.RPC_Client.RPC.CallResult(context.Background(), "DERO.GetSC", rpc.GetSC_Params{
		SCID:      scId,
		Variables: true,
		Code:      true,
	}, &result)
	if err != nil {
		return sc.UNKNOWN_TYPE, nil, err
	}

	if result.Code == "" {
		return sc.UNKNOWN_TYPE, nil, fmt.Errorf("token does not exists")
	}

	scType = sc.CheckType(result.Code)
	return scType, result, nil
}

type SCDetailsContainer struct {
	scIdEditor         *widget.Editor
	decimalsEditor     *widget.Editor
	standardTypeEditor *widget.Editor
	maxSupplyEditor    *widget.Editor
	nameEditor         *widget.Editor
	symbolEditor       *widget.Editor

	buttonAddToken *components.Button
	token          wallet_manager.Token
	list           *widget.List
}

func NewSCDetailsContainer(scId string, scType sc.SCType, scResult *rpc.GetSC_Result) *SCDetailsContainer {
	list := new(widget.List)
	list.Axis = layout.Vertical

	token := wallet_manager.Token{
		SCID:         scId,
		StandardType: scType,
	}

	switch scType {
	case sc.G45_FAT_TYPE:
		fat := g45_sc.G45_FAT{}
		err := fat.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		metadata := g45_sc.TokenMetadata{}
		err = metadata.Parse(fat.Metadata)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = metadata.Name
		token.Decimals = int64(fat.Decimals)
		token.MaxSupply = sql.NullInt64{Int64: int64(fat.MaxSupply), Valid: true}
		token.Image = sql.NullString{String: metadata.Image, Valid: true}
		token.Symbol = sql.NullString{String: metadata.Symbol, Valid: true}
		token.Metadata = sql.NullString{String: fat.Metadata, Valid: true}
	case sc.G45_AT_TYPE:
		at := g45_sc.G45_AT{}
		err := at.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		metadata := g45_sc.TokenMetadata{}
		err = metadata.Parse(at.Metadata)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = metadata.Name
		token.Decimals = int64(at.Decimals)
		token.MaxSupply = sql.NullInt64{Int64: int64(at.MaxSupply), Valid: true}
		token.Image = sql.NullString{String: metadata.Image, Valid: true}
		token.Symbol = sql.NullString{String: metadata.Symbol, Valid: true}
		token.Metadata = sql.NullString{String: at.Metadata, Valid: true}
	case sc.G45_NFT_TYPE:
		nft := g45_sc.G45_NFT{}
		err := nft.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		metadata := g45_sc.NFTMetadata{}
		err = metadata.Parse(nft.Metadata)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = metadata.Name
		token.Decimals = 0
		token.MaxSupply = sql.NullInt64{Int64: 1, Valid: true}
		token.Image = sql.NullString{String: metadata.Image, Valid: true}
		token.Metadata = sql.NullString{String: nft.Metadata, Valid: true}
	case sc.DEX_SC_TYPE:
		dex := dex_sc.SC{}
		err := dex.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = dex.Name
		token.Decimals = int64(dex.Decimals)
		token.Image = sql.NullString{String: dex.ImageUrl, Valid: true}
		token.Symbol = sql.NullString{String: dex.Symbol, Valid: true}
	case sc.UNKNOWN_TYPE:
		unknown := unknown_sc.SC{}
		err := unknown.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = unknown.Name
		token.Decimals = int64(unknown.Decimals)
		token.Image = sql.NullString{String: unknown.ImageUrl, Valid: true}
		token.Symbol = sql.NullString{String: unknown.Symbol, Valid: true}
	}

	scIdEditor := new(widget.Editor)
	scIdEditor.SetText(token.SCID)
	scIdEditor.WrapPolicy = text.WrapGraphemes
	scIdEditor.ReadOnly = true

	standardTypeEditor := new(widget.Editor)
	standardTypeEditor.SetText(string(scType))
	standardTypeEditor.WrapPolicy = text.WrapGraphemes
	standardTypeEditor.ReadOnly = true

	nameEditor := new(widget.Editor)
	nameEditor.SetText(token.Name)
	nameEditor.WrapPolicy = text.WrapGraphemes
	nameEditor.ReadOnly = true

	maxSupplyEditor := new(widget.Editor)
	maxSupply := "?"
	if token.MaxSupply.Valid {
		maxSupply = utils.ShiftNumber{Number: uint64(token.MaxSupply.Int64), Decimals: int(token.Decimals)}.LocaleString(language.English)
	}
	maxSupplyEditor.SetText(maxSupply)
	maxSupplyEditor.WrapPolicy = text.WrapGraphemes
	maxSupplyEditor.ReadOnly = true

	decimalsEditor := new(widget.Editor)
	decimals := fmt.Sprintf("%d", token.Decimals)
	decimalsEditor.SetText(decimals)
	decimalsEditor.WrapPolicy = text.WrapGraphemes
	decimalsEditor.ReadOnly = true

	symbolEditor := new(widget.Editor)
	symbol := "?"
	if token.Symbol.Valid {
		symbol = token.Symbol.String
	}
	symbolEditor.SetText(symbol)
	symbolEditor.WrapPolicy = text.WrapGraphemes
	symbolEditor.ReadOnly = true

	addIcon, _ := widget.NewIcon(icons.ContentAdd)
	buttonAddToken := components.NewButton(components.ButtonStyle{
		Rounded:         components.UniformRounded(unit.Dp(5)),
		Icon:            addIcon,
		TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		TextSize:        unit.Sp(14),
		IconGap:         unit.Dp(10),
		Inset:           layout.UniformInset(unit.Dp(10)),
		Animation:       components.NewButtonAnimationDefault(),
	})
	buttonAddToken.Label.Alignment = text.Middle
	buttonAddToken.Style.Font.Weight = font.Bold

	return &SCDetailsContainer{
		scIdEditor:         scIdEditor,
		nameEditor:         nameEditor,
		decimalsEditor:     decimalsEditor,
		maxSupplyEditor:    maxSupplyEditor,
		standardTypeEditor: standardTypeEditor,
		symbolEditor:       symbolEditor,

		token:          token,
		buttonAddToken: buttonAddToken,
		list:           list,
	}
}

func (sc *SCDetailsContainer) addToken() error {
	token := sc.token
	wallet := wallet_manager.OpenedWallet
	currentFolder := page_instance.pageSCFolders.currentFolder
	if currentFolder != nil {
		token.FolderId = sql.NullInt64{Int64: currentFolder.ID, Valid: true}
	}

	return wallet.InsertToken(token)
}

func (sc *SCDetailsContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if sc.buttonAddToken.Clicked() {
		err := sc.addToken()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText("Success", lang.Translate("New token added."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			page_instance.header.GoBack()
		}
	}

	var widgets []layout.Widget

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Max.Y = gtx.Dp(5)
		paint.FillShape(gtx.Ops, color.NRGBA{A: 150}, clip.Rect{
			Max: gtx.Constraints.Max,
		}.Op())

		return layout.Dimensions{Size: gtx.Constraints.Max}
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
					editor := material.Editor(th, sc.scIdEditor, "")
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
					editor := material.Editor(th, sc.nameEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Max Supply"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, sc.maxSupplyEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Decimals"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, sc.decimalsEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Symbol"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, sc.symbolEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Standard Type"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, sc.standardTypeEditor, "")
					editor.TextSize = unit.Sp(14)
					return editor.Layout(gtx)
				}),
			)
		})
		c := r.Stop()

		paint.FillShape(
			gtx.Ops,
			color.NRGBA{R: 255, G: 255, B: 255, A: 255},
			clip.UniformRRect(
				image.Rectangle{Max: dims.Size},
				gtx.Dp(10),
			).Op(gtx.Ops),
		)

		c.Add(gtx.Ops)
		return dims
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		sc.buttonAddToken.Text = lang.Translate("ADD TOKEN")
		return sc.buttonAddToken.Layout(gtx, th)
	})

	listStyle := material.List(th, sc.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(0), Right: unit.Dp(0),
		}.Layout(gtx, widgets[index])
	})
}
