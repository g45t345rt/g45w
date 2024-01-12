package page_wallet

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/g45t345rt/g45w/animation"
	"github.com/g45t345rt/g45w/app_instance"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/confirm_modal"
	"github.com/g45t345rt/g45w/containers/listselect_modal"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/containers/prompt_modal"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/router"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
	"github.com/g45t345rt/g45w/utils"
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCFolders struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	items []*TokenFolderItem

	list               *widget.List
	buttonOpenMenu     *components.Button
	buttonFolderGoBack *components.Button

	currentFolder   *wallet_manager.TokenFolder // nil is root
	folderCount     int
	tokenCount      int
	folderPath      string
	gridColumnCount int
}

var _ router.Page = &PageSCFolders{}

func NewPageSCFolders() *PageSCFolders {
	animationEnter := animation.NewAnimation(false, gween.NewSequence(
		gween.New(1, 0, .25, ease.Linear),
	))

	animationLeave := animation.NewAnimation(false, gween.NewSequence(
		gween.New(0, 1, .25, ease.Linear),
	))

	list := new(widget.List)
	list.Axis = layout.Vertical

	navIcon, _ := widget.NewIcon(icons.NavigationMenu)
	buttonOpenMenu := components.NewButton(components.ButtonStyle{
		Icon:      navIcon,
		Animation: components.NewButtonAnimationScale(.98),
	})

	backIcon, _ := widget.NewIcon(icons.ContentBackspace)
	buttonFolderGoBack := components.NewButton(components.ButtonStyle{
		Icon: backIcon,
	})

	page := &PageSCFolders{
		animationEnter:     animationEnter,
		animationLeave:     animationLeave,
		list:               list,
		buttonOpenMenu:     buttonOpenMenu,
		buttonFolderGoBack: buttonFolderGoBack,
	}

	page.SetLayout(settings.App.FolderLayout)

	return page
}

func (p *PageSCFolders) IsActive() bool {
	return p.isActive
}

func (p *PageSCFolders) Enter() {
	p.isActive = true
	page_instance.header.Title = func() string { return lang.Translate("Tokens") }
	page_instance.header.Subtitle = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		folderName := "root"
		if p.currentFolder != nil {
			folderName = p.currentFolder.Name
		}

		lbl := material.Label(th, unit.Sp(16), folderName)
		lbl.Color = theme.Current.TextMuteColor
		return lbl.Layout(gtx)
	}

	page_instance.header.LeftLayout = nil
	page_instance.header.RightLayout = func(gtx layout.Context, th *material.Theme) layout.Dimensions {
		p.buttonOpenMenu.Style.Colors = theme.Current.ButtonIconPrimaryColors
		gtx.Constraints.Min.X = gtx.Dp(30)
		gtx.Constraints.Min.Y = gtx.Dp(30)

		if p.buttonOpenMenu.Clicked(gtx) {
			go p.OpenMenu()
		}

		return p.buttonOpenMenu.Layout(gtx, th)
	}

	if !page_instance.header.IsHistory(PAGE_SC_FOLDERS) {
		p.animationEnter.Start()
		p.animationLeave.Reset()
	}

	p.Load()
}

func (p *PageSCFolders) Leave() {
	p.animationLeave.Start()
	p.animationEnter.Reset()
}

func (p *PageSCFolders) SetLayout(layout string) {
	switch layout {
	case settings.FolderLayoutGrid:
		p.gridColumnCount = 3
	case settings.FolderLayoutList:
		p.gridColumnCount = 1
	}

	settings.App.FolderLayout = layout
	settings.Save()
}

func (p *PageSCFolders) Load() error {
	wallet := wallet_manager.OpenedWallet
	folder := p.currentFolder

	folderId := sql.NullInt64{}
	if folder != nil {
		folderId.Scan(folder.ID)
	}

	folderPath, err := wallet.GetTokenFolderPath(folderId)
	if err != nil {
		return err
	}
	p.folderPath = folderPath

	folders, err := wallet.GetTokenFolderFolders(folderId)
	if err != nil {
		return err
	}

	p.items = make([]*TokenFolderItem, 0)

	for _, folder := range folders {
		count, err := wallet.GetTokenCount(sql.NullInt64{Int64: folder.ID, Valid: true})
		if err != nil {
			return err
		}

		p.items = append(p.items, NewTokenFolderItemFolder(folder, count))
	}

	p.folderCount = len(folders)

	tokens, err := wallet.GetTokens(wallet_manager.GetTokensParams{
		FolderId: &folderId,
	})
	if err != nil {
		return err
	}

	for _, token := range tokens {
		p.items = append(p.items, NewTokenFolderItemToken(token))
	}

	p.tokenCount = len(tokens)

	app_instance.Window.Invalidate()
	return nil
}

func (p *PageSCFolders) OpenMenu() {
	addIcon, _ := widget.NewIcon(icons.ActionNoteAdd)
	scanIcon, _ := widget.NewIcon(icons.ActionSearch)
	folderIcon, _ := widget.NewIcon(icons.FileCreateNewFolder)
	editIcon, _ := widget.NewIcon(icons.EditorBorderColor)
	listIcon, _ := widget.NewIcon(icons.ActionList)
	gridIcon, _ := widget.NewIcon(icons.ActionViewModule)
	refreshIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)

	var items []*listselect_modal.SelectListItem

	items = append(items, listselect_modal.NewSelectListItem("add_token",
		listselect_modal.NewItemText(addIcon, lang.Translate("Add token")).Layout,
	))

	items = append(items, listselect_modal.NewSelectListItem("scan_collection",
		listselect_modal.NewItemText(scanIcon, lang.Translate("Scan collection")).Layout,
	))

	items = append(items, listselect_modal.NewSelectListItem("new_folder",
		listselect_modal.NewItemText(folderIcon, lang.Translate("New folder")).Layout,
	))

	if page_instance.pageSCFolders.currentFolder != nil {
		items = append(items, listselect_modal.NewSelectListItem("rename_folder",
			listselect_modal.NewItemText(editIcon, lang.Translate("Rename folder")).Layout,
		))
	}

	if settings.App.FolderLayout == settings.FolderLayoutGrid {
		items = append(items, listselect_modal.NewSelectListItem("view_list",
			listselect_modal.NewItemText(listIcon, lang.Translate("View list")).Layout,
		))
	}

	if settings.App.FolderLayout == settings.FolderLayoutList {
		items = append(items, listselect_modal.NewSelectListItem("view_grid",
			listselect_modal.NewItemText(gridIcon, lang.Translate("View grid")).Layout,
		))
	}

	items = append(items, listselect_modal.NewSelectListItem("refresh_cache",
		listselect_modal.NewItemText(refreshIcon, lang.Translate("Refresh cache")).Layout,
	))

	if page_instance.pageSCFolders.currentFolder != nil {
		items = append(items, listselect_modal.NewSelectListItem("delete_folder",
			listselect_modal.NewItemText(deleteIcon, lang.Translate("Delete this folder")).Layout,
		))
	}

	items = append(items, listselect_modal.NewSelectListItem("remove_tokens",
		listselect_modal.NewItemText(deleteIcon, lang.Translate("Remove tokens")).Layout,
	))

	keyChan := listselect_modal.Instance.Open(items)
	for sKey := range keyChan {
		switch sKey {
		case "add_token":
			page_instance.pageRouter.SetCurrent(PAGE_ADD_SC_FORM)
			page_instance.header.AddHistory(PAGE_ADD_SC_FORM)
		case "scan_collection":
			page_instance.pageRouter.SetCurrent(PAGE_SCAN_COLLECTION)
			page_instance.header.AddHistory(PAGE_SCAN_COLLECTION)
		case "new_folder":
			wallet := wallet_manager.OpenedWallet
			currentFolder := page_instance.pageSCFolders.currentFolder

			txtChan := prompt_modal.Instance.Open("", lang.Translate("Enter folder name"), key.HintText)
			for folderName := range txtChan {
				tokenFolder := wallet_manager.TokenFolder{Name: folderName}
				if currentFolder != nil {
					parentId := sql.NullInt64{Int64: currentFolder.ID, Valid: true}
					tokenFolder.ParentId = parentId
				}

				err := wallet.InsertFolderToken(tokenFolder)
				if err != nil {
					notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
					notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				} else {
					page_instance.pageSCFolders.Load()
					notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("New folder created."))
					notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				}
			}
		case "rename_folder":
			wallet := wallet_manager.OpenedWallet
			currentFolder := page_instance.pageSCFolders.currentFolder

			txtChan := prompt_modal.Instance.Open(currentFolder.Name, lang.Translate("Rename folder"), key.HintText)
			for folderName := range txtChan {
				err := wallet.UpdateFolderToken(wallet_manager.TokenFolder{
					ID:       currentFolder.ID,
					Name:     folderName,
					ParentId: currentFolder.ParentId,
				})
				if err != nil {
					notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
					notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				} else {
					page_instance.pageSCFolders.Load()
					notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("Folder renamed."))
					notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				}

				currentFolder.Name = folderName
			}
		case "view_list":
			p.SetLayout(settings.FolderLayoutList)
		case "view_grid":
			p.SetLayout(settings.FolderLayoutGrid)
		case "refresh_cache":
			wallet := wallet_manager.OpenedWallet

			for _, item := range p.items {
				if item.token != nil {
					wallet.Memory.TokenAdd(item.token.GetHash())
					wallet.ResetBalanceResult(item.token.SCID)
				}
			}

			notification_modals.SuccessInstance.SetText("Success", lang.Translate("Cache refreshed."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		case "remove_tokens":
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			if <-yesChan {
				wallet := wallet_manager.OpenedWallet

				for _, item := range p.items {
					if item.token != nil { // not a folder
						wallet.DelToken(item.token.ID)
					}
				}

				notification_modals.SuccessInstance.SetText("Success", lang.Translate("Tokens removed."))
				notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				p.Load()
			}
		case "delete_folder":
			yesChan := confirm_modal.Instance.Open(confirm_modal.ConfirmText{})

			if <-yesChan {
				err := p.deleteCurrentFolder()
				if err != nil {
					notification_modals.ErrorInstance.SetText("Error", err.Error())
					notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
				} else {
					notification_modals.SuccessInstance.SetText("Success", lang.Translate("Folder and subfolders deleted."))
					notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
					p.changeFolder(p.currentFolder.ParentId)
				}
			}
		}
	}
}

func (p *PageSCFolders) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
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

	if p.buttonFolderGoBack.Clicked(gtx) {
		p.changeFolder(p.currentFolder.ParentId)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Bottom: unit.Dp(10),
						Left:   unit.Dp(30), Right: unit.Dp(30),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						r := op.Record(gtx.Ops)
						dims := layout.Inset{
							Top: unit.Dp(10), Bottom: unit.Dp(10),
							Left: unit.Dp(10), Right: unit.Dp(10),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if p.currentFolder != nil {
										p.buttonFolderGoBack.Style.Colors = theme.Current.ButtonIconPrimaryColors
										return p.buttonFolderGoBack.Layout(gtx, th)
									}
									return layout.Dimensions{}
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if p.currentFolder != nil {
										return layout.Spacer{Width: unit.Dp(10)}.Layout(gtx)
									}
									return layout.Dimensions{}
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									lbl := material.Label(th, unit.Sp(16), p.folderPath)
									return lbl.Layout(gtx)
								}),
							)
						})
						c := r.Stop()

						paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 50}, clip.UniformRRect(image.Rectangle{
							Max: dims.Size,
						}, gtx.Dp(10)).Op(gtx.Ops))

						c.Add(gtx.Ops)
						return dims
					})
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			widgets := []layout.ListElement{}

			if len(p.items) == 0 {
				widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
					lbl := material.Label(th, unit.Sp(16), lang.Translate("This folder is empty. Click the menu button to add folders or more tokens."))
					lbl.Color = theme.Current.TextMuteColor
					return lbl.Layout(gtx)
				})
			} else {
				widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							status := lang.Translate("{0} tokens - {1} folders")
							status = strings.Replace(status, "{0}", fmt.Sprint(p.tokenCount), -1)
							status = strings.Replace(status, "{1}", fmt.Sprint(p.folderCount), -1)
							lbl := material.Label(th, unit.Sp(16), status)
							lbl.Alignment = text.Middle
							return lbl.Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
					)
				})

				for i := 0; i < len(p.items); i += p.gridColumnCount {
					widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
						rowIndex := index - 1

						var childs []layout.FlexChild
						for a := 0; a < p.gridColumnCount; a++ {
							columnIndex := a
							itemIndex := (rowIndex * p.gridColumnCount) + columnIndex

							childs = append(childs,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									if itemIndex < len(p.items) {
										return p.items[itemIndex].Layout(gtx, th)
									}
									return layout.Dimensions{}
								}),
							)

							if columnIndex < p.gridColumnCount-1 {
								childs = append(childs, layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout))
							}
						}

						return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, childs...)
					})
				}
			}

			widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(30)}.Layout(gtx)
			})

			listStyle := material.List(th, p.list)
			listStyle.AnchorStrategy = material.Overlay

			return listStyle.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(30), Right: unit.Dp(30)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return widgets[index](gtx, index)
				})
			})
		}),
	)
}

func (p *PageSCFolders) deleteCurrentFolder() error {
	if p.currentFolder == nil {
		return fmt.Errorf("can't delete root folder")
	}

	wallet := wallet_manager.OpenedWallet
	return wallet.DelTokenFolder(p.currentFolder.ID)
}

func (p *PageSCFolders) changeFolder(id sql.NullInt64) error {
	if id.Valid {
		tokenFolder, _ := wallet_manager.OpenedWallet.GetTokenFolder(id.Int64)
		p.currentFolder = tokenFolder
	} else {
		p.currentFolder = nil
	}

	return p.Load()
}

type TokenFolderItem struct {
	clickable *widget.Clickable

	token      *wallet_manager.Token
	tokenImage *components.Image

	folder     *wallet_manager.TokenFolder
	folderIcon *widget.Icon

	name   string
	status string
}

func NewTokenFolderItemToken(token wallet_manager.Token) *TokenFolderItem {
	status := utils.ReduceTxId(token.SCID)

	tokenImage := &components.Image{
		Fit:     components.Cover,
		Rounded: components.UniformRounded(unit.Dp(10)),
	}

	return &TokenFolderItem{
		token:      &token,
		clickable:  new(widget.Clickable),
		tokenImage: tokenImage,
		name:       token.Name,
		status:     status,
	}
}

func NewTokenFolderItemFolder(folder wallet_manager.TokenFolder, tokenCount int) *TokenFolderItem {
	folderIcon, _ := widget.NewIcon(icons.FileFolder)

	status := lang.Translate("{} tokens")
	status = strings.Replace(status, "{}", fmt.Sprint(tokenCount), -1)

	return &TokenFolderItem{
		folder:     &folder,
		folderIcon: folderIcon,
		clickable:  new(widget.Clickable),
		name:       folder.Name,
		status:     status,
	}
}

func (item *TokenFolderItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.clickable.Clicked(gtx) {
		if item.folder != nil {
			id := sql.NullInt64{Int64: item.folder.ID, Valid: true}
			page_instance.pageSCFolders.changeFolder(id)
		}

		if item.token != nil {
			page_instance.pageSCToken.SetToken(item.token)
			page_instance.pageRouter.SetCurrent(PAGE_SC_TOKEN)
			page_instance.header.AddHistory(PAGE_SC_TOKEN)
			app_instance.Window.Invalidate()
		}
	}

	viewLayout := settings.App.FolderLayout

	switch viewLayout {
	case settings.FolderLayoutGrid:
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.Y = gtx.Constraints.Max.X
					paint.FillShape(gtx.Ops, theme.Current.ListBgColor, clip.UniformRRect(image.Rectangle{
						Max: gtx.Constraints.Max,
					}, gtx.Dp(10)).Op(gtx.Ops))

					if item.folderIcon != nil {
						return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return item.folderIcon.Layout(gtx, th.Fg)
						})
					}

					if item.tokenImage != nil {
						item.tokenImage.Src = item.token.LoadImageOp()
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return item.tokenImage.Layout(gtx)
						})
					}

					return layout.Dimensions{}
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				name := item.name
				size := unit.Sp(16)
				if len(name) > 20 {
					size = unit.Sp(14)
				}

				if len(name) > 30 {
					name = utils.ReduceString(name, 30, 0)
				}

				lbl := material.Label(th, size, name)
				lbl.Alignment = text.Middle
				lbl.Font.Weight = font.Bold
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(1)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Label(th, unit.Sp(14), item.status)
				lbl.Color = theme.Current.TextMuteColor
				lbl.Alignment = text.Middle
				return lbl.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		)
	case settings.FolderLayoutList:
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					r := op.Record(gtx.Ops)
					dims := layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints.Max = image.Point{X: gtx.Dp(50), Y: gtx.Dp(50)}
								gtx.Constraints.Min = gtx.Constraints.Max

								if item.folderIcon != nil {
									return item.folderIcon.Layout(gtx, th.Fg)
								}

								if item.tokenImage != nil {
									item.tokenImage.Src = item.token.LoadImageOp()
									return item.tokenImage.Layout(gtx)
								}

								return layout.Dimensions{}
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										lbl := material.Label(th, unit.Sp(16), item.name)
										lbl.Font.Weight = font.Bold
										return lbl.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Height: unit.Dp(1)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										lbl := material.Label(th, unit.Sp(16), item.status)
										lbl.Color = theme.Current.TextMuteColor
										return lbl.Layout(gtx)
									}),
								)
							}),
						)
					})
					c := r.Stop()

					paint.FillShape(gtx.Ops, theme.Current.ListBgColor, clip.UniformRRect(image.Rectangle{
						Max: dims.Size,
					}, gtx.Dp(10)).Op(gtx.Ops))

					c.Add(gtx.Ops)
					return dims
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
		)
	}

	return layout.Dimensions{}
}
