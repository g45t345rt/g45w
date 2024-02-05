package page_wallet

import (
	"database/sql"
	"fmt"
	"image"
	"sort"
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
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/sc/g45_sc"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageScanCollection struct {
	isActive bool

	headerPageAnimation *prefabs.PageHeaderAnimation

	buttonFetchData              *components.Button
	txtSCID                      *prefabs.TextField
	scCollectionDetailsContainer *SCCollectionDetailsContainer

	list *widget.List
}

var _ router.Page = &PageScanCollection{}

func NewPageScanCollection() *PageScanCollection {

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

	list := new(widget.List)
	list.Axis = layout.Vertical
	scCollectionDetailsContainer := NewSCCollectionDetailsContainer()

	headerPageAnimation := prefabs.NewPageHeaderAnimation(PAGE_SCAN_COLLECTION)
	return &PageScanCollection{
		headerPageAnimation:          headerPageAnimation,
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
	p.isActive = p.headerPageAnimation.Enter(page_instance.header)

	page_instance.header.Title = func() string { return lang.Translate("Scan Collection") }
	page_instance.header.Subtitle = nil
	page_instance.header.RightLayout = nil
}

func (p *PageScanCollection) Leave() {
	p.isActive = p.headerPageAnimation.Leave(page_instance.header)
}

func (p *PageScanCollection) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	defer p.headerPageAnimation.Update(gtx, func() { p.isActive = false }).Push(gtx.Ops).Pop()

	if p.buttonFetchData.Clicked(gtx) {
		go func() {
			p.buttonFetchData.SetLoading(true)
			scId, scType, scResult, err := p.submitForm()
			if err == nil {
				err = p.scCollectionDetailsContainer.Set(scId, scType, scResult)
			}
			p.buttonFetchData.SetLoading(false)

			if err != nil {
				notification_modal.Open(notification_modal.Params{
					Type:  notification_modal.ERROR,
					Title: lang.Translate("Error"),
					Text:  err.Error(),
				})
			}
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
			Left: theme.PagePadding, Right: theme.PagePadding,
		}.Layout(gtx, widgets[index])
	})
}

func (p *PageScanCollection) submitForm() (scId string, scType sc.SCType, result *rpc.GetSC_Result, err error) {
	scId = strings.TrimSpace(p.txtSCID.Value())
	if scId == "" {
		return scId, sc.UNKNOWN_TYPE, nil, fmt.Errorf("scid is empty")
	}

	err = wallet_manager.RPCCall("DERO.GetSC", rpc.GetSC_Params{
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

type ScanTokenBalanceResult struct {
	SCID    string
	Balance uint64
	Token   *wallet_manager.Token
	err     error
}

type SCCollectionDetailsContainer struct {
	scIdEditor        *widget.Editor
	nameEditor        *widget.Editor
	totalAssetsEditor *widget.Editor
	dateEditor        *widget.Editor
	collection        *g45_sc.G45_C
	buttonScan        *components.Button
	buttonStop        *components.Button
	buttonAddTokens   *components.Button
	withBalanceOnly   *widget.Bool

	scanIdx       int
	scanning      bool
	tokenBalances []ScanTokenBalanceResult

	list     *widget.List
	scanList *widget.List
}

func NewSCCollectionDetailsContainer() *SCCollectionDetailsContainer {
	list := new(widget.List)
	list.Axis = layout.Vertical

	scanList := new(widget.List)
	scanList.Axis = layout.Vertical

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

	playIcon, _ := widget.NewIcon(icons.AVPlayArrow)
	buttonScan := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      playIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonScan.Label.Alignment = text.Middle
	buttonScan.Style.Font.Weight = font.Bold

	stopIcon, _ := widget.NewIcon(icons.AVPause)
	buttonStop := components.NewButton(components.ButtonStyle{
		Rounded:   components.UniformRounded(unit.Dp(5)),
		Icon:      stopIcon,
		TextSize:  unit.Sp(14),
		IconGap:   unit.Dp(10),
		Inset:     layout.UniformInset(unit.Dp(10)),
		Animation: components.NewButtonAnimationDefault(),
	})
	buttonStop.Label.Alignment = text.Middle
	buttonStop.Style.Font.Weight = font.Bold

	addIcon, _ := widget.NewIcon(icons.AVLibraryAdd)
	loadingIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	buttonAddTokens := components.NewButton(components.ButtonStyle{
		Rounded:     components.UniformRounded(unit.Dp(5)),
		Icon:        addIcon,
		TextSize:    unit.Sp(14),
		IconGap:     unit.Dp(10),
		Inset:       layout.UniformInset(unit.Dp(10)),
		Animation:   components.NewButtonAnimationDefault(),
		LoadingIcon: loadingIcon,
	})
	buttonAddTokens.Label.Alignment = text.Middle
	buttonAddTokens.Style.Font.Weight = font.Bold

	return &SCCollectionDetailsContainer{
		scIdEditor:        scIdEditor,
		nameEditor:        nameEditor,
		totalAssetsEditor: totalAssetsEditor,
		dateEditor:        dateEditor,
		buttonScan:        buttonScan,
		buttonStop:        buttonStop,
		buttonAddTokens:   buttonAddTokens,
		withBalanceOnly:   new(widget.Bool),

		tokenBalances: make([]ScanTokenBalanceResult, 0),

		list:     list,
		scanList: scanList,
	}
}

func (c *SCCollectionDetailsContainer) Set(scId string, scType sc.SCType, scResult *rpc.GetSC_Result) error {
	c.tokenBalances = make([]ScanTokenBalanceResult, 0)
	c.scanIdx = 0

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

func (c *SCCollectionDetailsContainer) scan() {
	c.scanning = true
	wallet := wallet_manager.OpenedWallet
	addr := wallet.Memory.GetAddress().String()

	var assets []string
	for scId := range c.collection.Assets {
		assets = append(assets, scId)
	}
	sort.Strings(assets) // make sure its sorted once again ordering is not guaranteed in map

	for i := c.scanIdx; i < int(c.collection.AssetCount); i++ {
		if !c.scanning {
			break
		}

		scId := assets[i]
		hash := crypto.HashHexToHash(scId)
		balance, _, _ := wallet.Memory.GetDecryptedBalanceAtTopoHeight(hash, -1, addr)
		tokenBalance := ScanTokenBalanceResult{SCID: scId, Balance: balance}
		var err error

		result, cached, err := wallet_manager.GetSC(scId)
		if err == nil {
			token := &wallet_manager.Token{}
			tokenBalance.err = token.Parse(scId, result)
			if err == nil {
				tokenBalance.Token = token
			}
		}

		c.tokenBalances = append(c.tokenBalances, tokenBalance)

		if !cached {
			time.Sleep(time.Millisecond * 100)
		}

		c.scanIdx++
		app_instance.Window.Invalidate()
	}

	c.scanning = false
}

func (c *SCCollectionDetailsContainer) storeTokens() error {
	wallet := wallet_manager.OpenedWallet
	for _, tokenBalance := range c.tokenBalances {
		if c.withBalanceOnly.Value && tokenBalance.Balance == 0 {
			continue
		}

		token := tokenBalance.Token
		if token != nil {
			currentFolder := page_instance.pageSCFolders.currentFolder
			if currentFolder != nil {
				token.FolderId = sql.NullInt64{Int64: currentFolder.ID, Valid: true}
			}

			err := wallet.InsertToken(*token)
			if err != nil {
				return err
			}
		}
	}
	c.buttonAddTokens.SetLoading(false)
	page_instance.header.GoBack()
	return nil
}

func (c *SCCollectionDetailsContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	var widgets []layout.Widget

	if c.buttonScan.Clicked(gtx) {
		go c.scan()
	}

	if c.buttonStop.Clicked(gtx) {
		c.scanning = false
	}

	if c.buttonAddTokens.Clicked(gtx) {
		c.buttonAddTokens.SetLoading(true)
		go c.storeTokens()
	}

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		return prefabs.Divider(gtx, unit.Dp(5))
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

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		if c.scanning {
			c.buttonStop.Text = lang.Translate("STOP SCAN")
			c.buttonStop.Style.Colors = theme.Current.ButtonDangerColors
			return c.buttonStop.Layout(gtx, th)
		}

		if !c.scanning && c.scanIdx > 0 {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							c.buttonScan.Text = lang.Translate("RESUME")
							c.buttonScan.Style.Colors = theme.Current.ButtonPrimaryColors
							return c.buttonScan.Layout(gtx, th)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							c.buttonAddTokens.Text = lang.Translate("ADD TOKENS")
							c.buttonAddTokens.Style.Colors = theme.Current.ButtonPrimaryColors
							return c.buttonAddTokens.Layout(gtx, th)
						}),
					)
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis:      layout.Horizontal,
						Alignment: layout.Middle,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							s := material.Switch(th, c.withBalanceOnly, "")
							s.Color = theme.Current.SwitchColors
							return s.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(16), lang.Translate("Add only with balance"))
							return lbl.Layout(gtx)
						}),
					)
				}),
			)
		}

		c.buttonScan.Text = lang.Translate("SCAN TOKENS")
		c.buttonScan.Style.Colors = theme.Current.ButtonPrimaryColors
		return c.buttonScan.Layout(gtx, th)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		status := fmt.Sprintf("%d / %d", len(c.tokenBalances), c.collection.AssetCount)
		lbl := material.Label(th, unit.Sp(16), status)
		lbl.Font.Weight = font.Bold
		return lbl.Layout(gtx)
	})

	widgets = append(widgets, func(gtx layout.Context) layout.Dimensions {
		var scanWidgets []layout.Widget

		for i := range c.tokenBalances {
			idx := len(c.tokenBalances) - 1 - i
			scanWidgets = append(scanWidgets, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(2)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					hashBalance := c.tokenBalances[idx]
					scId := utils.ReduceTxId(hashBalance.SCID)
					status := ""
					if hashBalance.err != nil {
						status = lang.Translate("error")
					} else {
						status = fmt.Sprint(hashBalance.Balance)
					}

					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(14), fmt.Sprintf("%d. ", idx+1))
							return lbl.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(14), scId)
							return lbl.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							lbl := material.Label(th, unit.Sp(14), status)
							return lbl.Layout(gtx)
						}),
					)
				})
			})
		}

		listStyle := material.List(th, c.scanList)
		gtx.Constraints.Max.Y = gtx.Dp(200)

		return listStyle.Layout(gtx, len(scanWidgets), func(gtx layout.Context, index int) layout.Dimensions {
			return scanWidgets[index](gtx)
		})
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
