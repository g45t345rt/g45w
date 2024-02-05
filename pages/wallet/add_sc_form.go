package page_wallet

import (
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
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"golang.org/x/text/language"
)

type PageAddSCForm struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	scDetailsContainer *SCDetailsContainer

	buttonFetchData *components.Button
	txtSCID         *prefabs.TextField

	list *widget.List
}

var _ router.Page = &PageAddSCForm{}

func NewPageAddSCForm() *PageAddSCForm {

	list := new(widget.List)
	list.Axis = layout.Vertical

	searchIcon, _ := widget.NewIcon(icons.ActionSearch)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonFetchData := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        searchIcon,
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

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_ADD_SC_FORM)

	return &PageAddSCForm{
		headerPageAnimation: headerPageAnimation,
		txtSCID:             txtSCID,
		buttonFetchData:     buttonFetchData,
		scDetailsContainer:  scDetailsContainer,
		list:                list,
	}
}

func (p *PageAddSCForm) IsActive() bool {
	return p.isActive
}

func (p *PageAddSCForm) Enter() {
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Add Token") }
	page_instance.header.Subtitle = nil
	page_instance.header.RightLayout = nil
}

func (p *PageAddSCForm) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageAddSCForm) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonFetchData.Clicked(gtx) {
		go func() {
			p.scDetailsContainer.token = nil
			p.buttonFetchData.SetLoading(true)
			token, err := p.submitForm()
			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			}

			p.buttonFetchData.SetLoading(false)
			p.scDetailsContainer.Set(token)
		}()
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
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageAddSCForm) submitForm() (*wallet_manager.Token, error) {
	scId := strings.TrimSpace(p.txtSCID.Value())
	if scId == "" {
		return nil, fmt.Errorf("scid is empty")
	}

	result, _, err := wallet_manager.GetSC(scId)
	if err != nil {
		return nil, err
	}

	token := &wallet_manager.Token{}
	err = token.Parse(scId, result)
	if err != nil {
		return nil, err
	}

	return token, nil
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

func (c *SCDetailsContainer) Set(token *wallet_manager.Token) {
	c.scIdEditor.SetText(token.SCID)
	c.standardTypeEditor.SetText(string(token.StandardType))
	c.nameEditor.SetText(token.Name)

	if token.CreatedTimestamp.Valid && token.CreatedTimestamp.Int64 > 0 {
		date := time.Unix(token.CreatedTimestamp.Int64, 0)
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
	if c.buttonAddToken.Clicked(gtx) {
		err := c.addToken()
		if err != nil {
			notification_modal.Open(notification_modal.Params{
				Type:  notification_modal.ERROR,
				Title: lang.Translate("Error"),
				Text:  err.Error(),
			})
		} else {
			notification_modal.Open(notification_modal.Params{
				Type:       notification_modal.SUCCESS,
				Title:      lang.Translate("Success"),
				Text:       lang.Translate("New token added."),
				CloseAfter: notification_modal.CLOSE_AFTER_DEFAULT,
			})
			page_instance.header.GoBack()
		}
	}

	var widgets []layout.Widget

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return prefabs.Divider(gtx, unit.Dp(5))
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		r := op.Record(gtx.Ops)
		var childs []layout.FlexChild

		childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
		)

		childs = append(childs,
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
		)

		childs = append(childs,
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
		)

		childs = append(childs,
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
		)

		childs = append(childs,
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
		)

		childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
		)

		if c.dateEditor.Text() != "" {
			childs = append(childs, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
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
		}

		dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, childs...)
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
