package page_wallet

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"strings"

	"gioui.org/font"
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
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
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

	list                *widget.List
	buttonOpenMenu      *components.Button
	folderMenuSelect    *FolderMenuSelect
	createFolderModal   *CreateFolderModal
	confirmDeleteFolder *prefabs.Confirm
	confirmRemoveTokens *prefabs.Confirm
	buttonFolderGoBack  *components.Button

	currentFolder *wallet_manager.TokenFolder // nil is root
	folderCount   int
	tokenCount    int
	folderPath    string
	//viewLayout      string // "grid", "list"
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

	addIcon, _ := widget.NewIcon(icons.NavigationMenu)
	buttonOpenMenu := components.NewButton(components.ButtonStyle{
		Icon:      addIcon,
		Animation: components.NewButtonAnimationScale(.98),
	})

	backIcon, _ := widget.NewIcon(icons.ContentBackspace)
	buttonFolderGoBack := components.NewButton(components.ButtonStyle{
		Icon: backIcon,
	})

	createFolderModal := NewCreateFolderModal()
	confirmDeleteFolder := prefabs.NewConfirm(layout.Center)
	confirmRemoveTokens := prefabs.NewConfirm(layout.Center)

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			createFolderModal.Layout(gtx, th)

			confirmDeleteFolder.Layout(gtx, th, prefabs.ConfirmText{
				Prompt: lang.Translate("Are you sure?"),
				No:     lang.Translate("NO"),
				Yes:    lang.Translate("YES"),
			})

			confirmRemoveTokens.Layout(gtx, th, prefabs.ConfirmText{
				Prompt: lang.Translate("Are you sure?"),
				No:     lang.Translate("NO"),
				Yes:    lang.Translate("YES"),
			})
		},
	})

	page := &PageSCFolders{
		animationEnter:      animationEnter,
		animationLeave:      animationLeave,
		list:                list,
		buttonOpenMenu:      buttonOpenMenu,
		folderMenuSelect:    NewFolderMenuSelect(),
		createFolderModal:   createFolderModal,
		confirmDeleteFolder: confirmDeleteFolder,
		buttonFolderGoBack:  buttonFolderGoBack,
		confirmRemoveTokens: confirmRemoveTokens,
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
	page_instance.header.ButtonRight = p.buttonOpenMenu

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

	if p.buttonOpenMenu.Clicked() {
		p.folderMenuSelect.SelectModal.Modal.SetVisible(true)
	}

	selected, key := p.folderMenuSelect.SelectModal.Selected()
	if selected {
		switch key {
		case "add_token":
			page_instance.pageRouter.SetCurrent(PAGE_ADD_SC_FORM)
			page_instance.header.AddHistory(PAGE_ADD_SC_FORM)
		case "scan_collection":
			page_instance.pageRouter.SetCurrent(PAGE_SCAN_COLLECTION)
			page_instance.header.AddHistory(PAGE_SCAN_COLLECTION)
		case "new_folder":
			p.createFolderModal.setFolder(nil)
			p.createFolderModal.modal.SetVisible(true)
		case "rename_folder":
			p.createFolderModal.setFolder(p.currentFolder)
			p.createFolderModal.modal.SetVisible(true)
		case "view_list":
			p.SetLayout(settings.FolderLayoutList)
		case "view_grid":
			p.SetLayout(settings.FolderLayoutGrid)
		case "refresh_cache":
			go func() {
				wallet := wallet_manager.OpenedWallet

				for _, item := range p.items {
					if item.token != nil {
						wallet.Memory.TokenAdd(item.token.GetHash())
						wallet.ResetBalanceResult(item.token.SCID)
					}
				}

				notification_modals.SuccessInstance.SetText("Success", lang.Translate("Cache refreshed."))
				notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			}()
		case "remove_tokens":
			p.confirmRemoveTokens.SetVisible(true)
		case "delete_folder":
			p.confirmDeleteFolder.SetVisible(true)
		}
		p.folderMenuSelect.SelectModal.Modal.SetVisible(false)
	}

	if p.confirmDeleteFolder.ClickedYes() {
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

	if p.confirmRemoveTokens.ClickedYes() {
		go func() {
			wallet := wallet_manager.OpenedWallet

			for _, item := range p.items {
				if item.token != nil { // not a folder
					wallet.DelToken(item.token.ID)
				}
			}

			notification_modals.SuccessInstance.SetText("Success", lang.Translate("Tokens removed."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			p.Load()
		}()
	}

	if p.buttonFolderGoBack.Clicked() {
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

type CreateFolderModal struct {
	folder        *wallet_manager.TokenFolder
	modal         *components.Modal
	txtFolderName *components.Input
}

func NewCreateFolderModal() *CreateFolderModal {
	modal := components.NewModal(components.ModalStyle{
		CloseOnOutsideClick: true,
		CloseOnInsideClick:  false,
		Direction:           layout.Center,
		Rounded: components.Rounded{
			SW: unit.Dp(10), SE: unit.Dp(10),
			NW: unit.Dp(10), NE: unit.Dp(10),
		},
		Animation: components.NewModalAnimationScaleBounce(),
	})

	txtFolderName := components.NewInput()
	txtFolderName.Border = widget.Border{}
	txtFolderName.Inset = layout.Inset{}
	txtFolderName.TextSize = unit.Sp(20)

	return &CreateFolderModal{
		modal:         modal,
		txtFolderName: txtFolderName,
	}
}

func (c *CreateFolderModal) addFolder(name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	wallet := wallet_manager.OpenedWallet

	tokenFolder := wallet_manager.TokenFolder{Name: name}
	currentFolder := page_instance.pageSCFolders.currentFolder
	if currentFolder != nil {
		parentId := sql.NullInt64{Int64: currentFolder.ID, Valid: true}
		tokenFolder.ParentId = parentId
	}

	return wallet.InsertFolderToken(tokenFolder)
}

func (c *CreateFolderModal) renameFolder(name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	wallet := wallet_manager.OpenedWallet
	err := wallet.UpdateFolderToken(wallet_manager.TokenFolder{
		ID:       c.folder.ID,
		Name:     name,
		ParentId: c.folder.ParentId,
	})
	if err != nil {
		return err
	}

	c.folder.Name = name
	return nil
}

func (c *CreateFolderModal) setFolder(folder *wallet_manager.TokenFolder) {
	c.folder = folder

	if c.folder != nil {
		name := c.folder.Name
		c.txtFolderName.SetValue(name)
		c.txtFolderName.Editor.SetCaret(len(name), len(name)) // move caret at the end
	} else {
		c.txtFolderName.SetValue("")
	}
}

func (c *CreateFolderModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range c.txtFolderName.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok {
			successMsg := ""
			var err error

			if c.folder == nil {
				err = c.addFolder(e.Text)
				successMsg = lang.Translate("New folder created.")
			} else {
				err = c.renameFolder(e.Text)
				successMsg = lang.Translate("Folder renamed.")
			}

			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			} else {
				c.txtFolderName.SetValue("")
				c.modal.SetVisible(false)
				page_instance.pageSCFolders.Load()
				notification_modals.SuccessInstance.SetText(lang.Translate("Success"), successMsg)
				notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			}
		}
	}

	if c.modal.Visible {
		if !c.txtFolderName.Editor.Focused() {
			c.txtFolderName.Editor.Focus()
		}
	}

	c.modal.Style.Colors = theme.Current.ModalColors
	return c.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(20), Bottom: unit.Dp(20),
			Left: unit.Dp(20), Right: unit.Dp(20),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(200)
			gtx.Constraints.Max.X = gtx.Dp(200)
			c.txtFolderName.Colors = theme.Current.InputColors
			return c.txtFolderName.Layout(gtx, th, lang.Translate("Enter folder name"))
		})
	})
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

	if item.clickable.Clicked() {
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
				lbl := material.Label(th, unit.Sp(16), item.name)
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

type FolderMenuSelect struct {
	SelectModal *prefabs.SelectModal
}

func NewFolderMenuSelect() *FolderMenuSelect {
	var items []*prefabs.SelectListItem
	addIcon, _ := widget.NewIcon(icons.ActionNoteAdd)
	items = append(items, prefabs.NewSelectListItem("add_token", prefabs.ListItemMenuItem{
		Icon:  addIcon,
		Title: "Add token", //@lang.Translate("Add token")
	}.Layout))

	scanIcon, _ := widget.NewIcon(icons.ActionSearch)
	items = append(items, prefabs.NewSelectListItem("scan_collection", prefabs.ListItemMenuItem{
		Icon:  scanIcon,
		Title: "Scan collection", //@lang.Translate("Scan collection")
	}.Layout))

	folderIcon, _ := widget.NewIcon(icons.FileCreateNewFolder)
	items = append(items, prefabs.NewSelectListItem("new_folder", prefabs.ListItemMenuItem{
		Icon:  folderIcon,
		Title: "New folder", //@lang.Translate("New folder")
	}.Layout))

	editIcon, _ := widget.NewIcon(icons.EditorBorderColor)
	items = append(items, prefabs.NewSelectListItem("rename_folder", prefabs.ListItemMenuItem{
		Icon:  editIcon,
		Title: "Rename folder", //@lang.Translate("Rename folder")
	}.Layout))

	listIcon, _ := widget.NewIcon(icons.ActionList)
	items = append(items, prefabs.NewSelectListItem("view_list", prefabs.ListItemMenuItem{
		Icon:  listIcon,
		Title: "View list", //@lang.Translate("View list")
	}.Layout))

	gridIcon, _ := widget.NewIcon(icons.ActionViewModule)
	items = append(items, prefabs.NewSelectListItem("view_grid", prefabs.ListItemMenuItem{
		Icon:  gridIcon,
		Title: "View grid", //@lang.Translate("View grid")
	}.Layout))

	refreshIcon, _ := widget.NewIcon(icons.NavigationRefresh)
	items = append(items, prefabs.NewSelectListItem("refresh_cache", prefabs.ListItemMenuItem{
		Icon:  refreshIcon,
		Title: "Refresh cache", //@lang.Translate("Refresh cache")
	}.Layout))

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	items = append(items, prefabs.NewSelectListItem("delete_folder", prefabs.ListItemMenuItem{
		Icon:  deleteIcon,
		Title: "Delete this folder", //@lang.Translate("Delete this folder")
	}.Layout))

	items = append(items, prefabs.NewSelectListItem("remove_tokens", prefabs.ListItemMenuItem{
		Icon:  deleteIcon,
		Title: "Remove tokens", //@lang.Translate("Remove tokens")
	}.Layout))

	selectModal := prefabs.NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			var filteredItems []*prefabs.SelectListItem
			for _, item := range items {
				add := true
				switch item.Key {
				case "rename_folder":
					add = page_instance.pageSCFolders.currentFolder != nil
				case "delete_folder":
					add = page_instance.pageSCFolders.currentFolder != nil
				case "view_grid":
					add = settings.App.FolderLayout == settings.FolderLayoutList
				case "view_list":
					add = settings.App.FolderLayout == settings.FolderLayoutGrid
				}

				if add {
					filteredItems = append(filteredItems, item)
				}
			}

			selectModal.Layout(gtx, th, filteredItems)
		},
	})

	return &FolderMenuSelect{
		SelectModal: selectModal,
	}
}
