package page_wallet

import (
	"context"
	"database/sql"
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
	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/sc/g45_sc"
	"github.com/g45t345rt/g45w/sc/unknown_sc"
	"github.com/g45t345rt/g45w/theme"
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
	txtSCID         *prefabs.TextField

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
	scDetailsContainer := NewSCDetailsContainer()

	return &PageAddSCForm{
		animationEnter:     animationEnter,
		animationLeave:     animationLeave,
		txtSCID:            txtSCID,
		buttonFetchData:    buttonFetchData,
		scDetailsContainer: scDetailsContainer,
		list:               list,
	}
}

func (p *PageAddSCForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddSCForm) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Add Token") }
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
		p.scDetailsContainer.token = nil
		p.buttonFetchData.SetLoading(true)
		scId, scType, scResult, err := p.submitForm()
		if err == nil {
			err = p.scDetailsContainer.Set(scId, scType, scResult)
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

	if p.scDetailsContainer.token != nil {
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

func (p *PageAddSCForm) submitForm() (scId string, scType sc.SCType, result *rpc.GetSC_Result, err error) {
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
		return scId, sc.UNKNOWN_TYPE, nil, fmt.Errorf("token does not exists")
	}

	scType = sc.CheckType(result.Code)
	return scId, scType, result, nil
}

type SCDetailsContainer struct {
	scIdEditor         *widget.Editor
	decimalsEditor     *widget.Editor
	standardTypeEditor *widget.Editor
	maxSupplyEditor    *widget.Editor
	nameEditor         *widget.Editor
	symbolEditor       *widget.Editor
	dateEditor         *widget.Editor

	buttonAddToken *components.Button
	token          *wallet_manager.Token
	list           *widget.List
}

func NewSCDetailsContainer() *SCDetailsContainer {
	list := new(widget.List)
	list.Axis = layout.Vertical

	scIdEditor := new(widget.Editor)
	scIdEditor.WrapPolicy = text.WrapGraphemes
	scIdEditor.ReadOnly = true

	standardTypeEditor := new(widget.Editor)
	standardTypeEditor.WrapPolicy = text.WrapGraphemes
	standardTypeEditor.ReadOnly = true

	nameEditor := new(widget.Editor)
	nameEditor.WrapPolicy = text.WrapGraphemes
	nameEditor.ReadOnly = true

	maxSupplyEditor := new(widget.Editor)
	maxSupplyEditor.WrapPolicy = text.WrapGraphemes
	maxSupplyEditor.ReadOnly = true

	decimalsEditor := new(widget.Editor)
	decimalsEditor.WrapPolicy = text.WrapGraphemes
	decimalsEditor.ReadOnly = true

	symbolEditor := new(widget.Editor)
	symbolEditor.WrapPolicy = text.WrapGraphemes
	symbolEditor.ReadOnly = true

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

	return &SCDetailsContainer{
		scIdEditor:         scIdEditor,
		nameEditor:         nameEditor,
		decimalsEditor:     decimalsEditor,
		maxSupplyEditor:    maxSupplyEditor,
		standardTypeEditor: standardTypeEditor,
		symbolEditor:       symbolEditor,
		dateEditor:         dateEditor,

		buttonAddToken: buttonAddToken,
		list:           list,
	}
}

func (c *SCDetailsContainer) Set(scId string, scType sc.SCType, scResult *rpc.GetSC_Result) error {
	token := &wallet_manager.Token{
		SCID:         scId,
		StandardType: scType,
	}

	var timestamp uint64
	switch scType {
	case sc.G45_FAT_TYPE:
		fat := g45_sc.G45_FAT{}
		err := fat.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			return err
		}

		metadata := g45_sc.TokenMetadata{}
		err = metadata.Parse(fat.Metadata)
		if err != nil {
			return err
		}

		token.Name = metadata.Name
		token.Decimals = int64(fat.Decimals)
		token.MaxSupply = sql.NullInt64{Int64: int64(fat.MaxSupply), Valid: true}
		token.ImageUrl = sql.NullString{String: metadata.Image, Valid: true}
		token.Symbol = sql.NullString{String: metadata.Symbol, Valid: true}
		token.Metadata = sql.NullString{String: fat.Metadata, Valid: true}
		timestamp = fat.Timestamp
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
		token.ImageUrl = sql.NullString{String: metadata.Image, Valid: true}
		token.Symbol = sql.NullString{String: metadata.Symbol, Valid: true}
		token.Metadata = sql.NullString{String: at.Metadata, Valid: true}
		timestamp = at.Timestamp
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
		token.ImageUrl = sql.NullString{String: metadata.Image, Valid: true}
		token.Metadata = sql.NullString{String: nft.Metadata, Valid: true}
		timestamp = nft.Timestamp
	case sc.DEX_SC_TYPE:
		dex := dex_sc.SC{}
		err := dex.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = dex.Name
		token.Decimals = int64(dex.Decimals)
		token.ImageUrl = sql.NullString{String: dex.ImageUrl, Valid: true}
		token.Symbol = sql.NullString{String: dex.Symbol, Valid: true}
	case sc.UNKNOWN_TYPE:
		unknown := unknown_sc.SC{}
		err := unknown.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			fmt.Println(err)
		}

		token.Name = unknown.Name
		token.Decimals = int64(unknown.Decimals)
		token.ImageUrl = sql.NullString{String: unknown.ImageUrl, Valid: true}
		token.Symbol = sql.NullString{String: unknown.Symbol, Valid: true}
	}

	c.scIdEditor.SetText(token.SCID)
	c.standardTypeEditor.SetText(string(scType))
	c.nameEditor.SetText(token.Name)

	if timestamp != 0 {
		date := time.Unix(int64(timestamp), 0)
		c.dateEditor.SetText(date.Format("2006-01-02 15:04:05"))
	}

	maxSupply := "?"
	if token.MaxSupply.Valid {
		maxSupply = utils.ShiftNumber{Number: uint64(token.MaxSupply.Int64), Decimals: int(token.Decimals)}.LocaleString(language.English)
	}
	c.maxSupplyEditor.SetText(maxSupply)

	decimals := fmt.Sprintf("%d", token.Decimals)
	c.decimalsEditor.SetText(decimals)

	symbol := "?"
	if token.Symbol.Valid {
		symbol = token.Symbol.String
	}
	c.symbolEditor.SetText(symbol)
	c.token = token

	return nil
}

func (c *SCDetailsContainer) addToken() error {
	token := c.token
	wallet := wallet_manager.OpenedWallet
	currentFolder := page_instance.pageSCFolders.currentFolder
	if currentFolder != nil {
		token.FolderId = sql.NullInt64{Int64: currentFolder.ID, Valid: true}
	}

	return wallet.InsertToken(*token)
}

func (c *SCDetailsContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if c.buttonAddToken.Clicked() {
		err := c.addToken()
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
					lbl := material.Label(th, unit.Sp(16), lang.Translate("Max Supply"))
					lbl.Font.Weight = font.Bold
					return lbl.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					editor := material.Editor(th, c.maxSupplyEditor, "")
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
					editor := material.Editor(th, c.decimalsEditor, "")
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
					editor := material.Editor(th, c.symbolEditor, "")
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
					editor := material.Editor(th, c.standardTypeEditor, "")
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

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		c.buttonAddToken.Text = lang.Translate("ADD TOKEN")
		c.buttonAddToken.Style.Colors = theme.Current.ButtonPrimaryColors
		return c.buttonAddToken.Layout(gtx, th)
	})

	listStyle := material.List(th, c.list)
	listStyle.AnchorStrategy = material.Overlay

	return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(0), Bottom: unit.Dp(20),
			Left: unit.Dp(0), Right: unit.Dp(0),
		}.Layout(gtx, widgets[index])
	})
}
