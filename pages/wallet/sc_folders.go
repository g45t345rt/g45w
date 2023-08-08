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
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
	"github.com/g45t345rt/g45w/containers/notification_modals"
	"github.com/g45t345rt/g45w/lang"
	"github.com/g45t345rt/g45w/prefabs"
	"github.com/g45t345rt/g45w/router"
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
	deleteFolderConfirm *components.Confirm
	buttonFolderGoBack  *components.Button

	currentFolder *wallet_manager.TokenFolder // nil is root
	folderCount   int
	tokenCount    int
	folderPath    string
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
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		Animation: components.NewButtonAnimationScale(.98),
	})

	backIcon, _ := widget.NewIcon(icons.ContentBackspace)
	buttonFolderGoBack := components.NewButton(components.ButtonStyle{
		Icon:      backIcon,
		TextColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	})

	createFolderModal := NewCreateFolderModal()
	deleteFolderConfirm := components.NewConfirm(layout.Center)

	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			createFolderModal.Layout(gtx, th)

			deleteFolderConfirm.Prompt = lang.Translate("Are you sure?")
			deleteFolderConfirm.NoText = lang.Translate("NO")
			deleteFolderConfirm.YesText = lang.Translate("YES")
			deleteFolderConfirm.Layout(gtx, th)
		},
	})

	return &PageSCFolders{
		animationEnter:      animationEnter,
		animationLeave:      animationLeave,
		list:                list,
		buttonOpenMenu:      buttonOpenMenu,
		folderMenuSelect:    NewFolderMenuSelect(),
		createFolderModal:   createFolderModal,
		deleteFolderConfirm: deleteFolderConfirm,
		buttonFolderGoBack:  buttonFolderGoBack,
	}
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
		case "new_folder":
			p.createFolderModal.setFolder(nil)
			p.createFolderModal.modal.SetVisible(true)
		case "rename_folder":
			p.createFolderModal.setFolder(p.currentFolder)
			p.createFolderModal.modal.SetVisible(true)
		case "delete_folder":
			p.deleteFolderConfirm.SetVisible(true)
		}
		p.folderMenuSelect.SelectModal.Modal.SetVisible(false)
	}

	if p.deleteFolderConfirm.ClickedYes() {
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

				start := len(widgets)
				for i := 0; i < len(p.items); i += 3 {
					widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
						itemIndex := (index - start) * 3

						dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if itemIndex < len(p.items) {
									return p.items[itemIndex].Layout(gtx, th)
								}
								return layout.Dimensions{}
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if itemIndex+1 < len(p.items) {
									return p.items[itemIndex+1].Layout(gtx, th)
								}
								return layout.Dimensions{}
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if itemIndex+2 < len(p.items) {
									return p.items[itemIndex+2].Layout(gtx, th)
								}
								return layout.Dimensions{}
							}),
						)
						return dims
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
		BgColor:             color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Rounded: components.Rounded{
			SW: unit.Dp(10), SE: unit.Dp(10),
			NW: unit.Dp(10), NE: unit.Dp(10),
		},
		Animation: components.NewModalAnimationScaleBounce(),
		Backdrop:  components.NewModalBackground(),
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
		c.txtFolderName.SetValue(c.folder.Name)
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

	return c.modal.Layout(gtx, nil, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(20), Bottom: unit.Dp(20),
			Left: unit.Dp(20), Right: unit.Dp(20),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Dp(200)
			gtx.Constraints.Max.X = gtx.Dp(200)
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
	img, _ := assets.GetImage("token.png")
	status := utils.ReduceTxId(token.SCID)

	tokenImage := &components.Image{
		Src:     paint.NewImageOp(img),
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
			page_instance.pageSCToken.token = item.token
			page_instance.pageRouter.SetCurrent(PAGE_SC_TOKEN)
			page_instance.header.AddHistory(PAGE_SC_TOKEN)
			app_instance.Window.Invalidate()
		}
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Constraints.Max.X
				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.UniformRRect(image.Rectangle{
					Max: gtx.Constraints.Max,
				}, gtx.Dp(10)).Op(gtx.Ops))

				if item.folderIcon != nil {
					return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return item.folderIcon.Layout(gtx, color.NRGBA{A: 255})
					})
				}

				if item.tokenImage != nil {
					return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
			lbl := material.Label(th, unit.Sp(12), item.status)
			lbl.Alignment = text.Middle
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
	)
}

type FolderMenuItem struct {
	Key   string
	Icon  *widget.Icon
	Title string
}

func (t FolderMenuItem) Layout(gtx layout.Context, index int, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.X = gtx.Dp(45)
			gtx.Constraints.Max.Y = gtx.Dp(30)
			return t.Icon.Layout(gtx, color.NRGBA{A: 255})
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(20), lang.Translate(t.Title))
			return lbl.Layout(gtx)
		}),
	)
}

type FolderMenuSelect struct {
	SelectModal *prefabs.SelectModal
}

func NewFolderMenuSelect() *FolderMenuSelect {
	var items []*prefabs.SelectListItem
	addIcon, _ := widget.NewIcon(icons.ActionNoteAdd)
	items = append(items, prefabs.NewSelectListItem("add_token", FolderMenuItem{
		Icon:  addIcon,
		Title: "Add token", //@lang.Translate("Add token")
	}.Layout))

	folderIcon, _ := widget.NewIcon(icons.FileCreateNewFolder)
	items = append(items, prefabs.NewSelectListItem("new_folder", FolderMenuItem{
		Icon:  folderIcon,
		Title: "New folder", //@lang.Translate("New folder")
	}.Layout))

	editIcon, _ := widget.NewIcon(icons.EditorBorderColor)
	items = append(items, prefabs.NewSelectListItem("rename_folder", FolderMenuItem{
		Icon:  editIcon,
		Title: "Rename folder", //@lang.Translate("Rename folder")
	}.Layout))

	listIcon, _ := widget.NewIcon(icons.ActionList)
	items = append(items, prefabs.NewSelectListItem("view_list", FolderMenuItem{
		Icon:  listIcon,
		Title: "View list", //@lang.Translate("View list")
	}.Layout))

	viewFolderIcon, _ := widget.NewIcon(icons.FileFolder)
	items = append(items, prefabs.NewSelectListItem("view_folder", FolderMenuItem{
		Icon:  viewFolderIcon,
		Title: "View folder", //@lang.Translate("View folder")
	}.Layout))

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	items = append(items, prefabs.NewSelectListItem("delete_folder", FolderMenuItem{
		Icon:  deleteIcon,
		Title: "Delete this folder", //@lang.Translate("Delete this folder")
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
