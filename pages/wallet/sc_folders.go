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
	"github.com/g45t345rt/g45w/wallet_manager"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type PageSCFolders struct {
	isActive bool

	animationEnter *animation.Animation
	animationLeave *animation.Animation

	tokenItems []*TokenFolderItem

	list                *widget.List
	buttonOpenMenu      *components.Button
	tokenMenuSelect     *TokenMenuSelect
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
		tokenMenuSelect:     NewTokenMenuSelect(),
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
	page_instance.header.SetTitle(lang.Translate("Tokens"))
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

	folderId := sql.NullInt32{}
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

	p.tokenItems = make([]*TokenFolderItem, 0)

	for _, folder := range folders {
		p.tokenItems = append(p.tokenItems, NewTokenFolderItem(folder))
	}

	p.folderCount = len(folders)

	tokens, err := wallet.GetTokens(wallet_manager.GetTokensParams{
		FolderId: folderId,
	})
	if err != nil {
		return err
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

	if p.buttonOpenMenu.Clickable.Clicked() {
		p.tokenMenuSelect.SelectModal.Modal.SetVisible(true)
	}

	selected := p.tokenMenuSelect.SelectModal.Selected()
	if selected {
		switch p.tokenMenuSelect.SelectModal.SelectedKey {
		case "add_token":
			page_instance.pageRouter.SetCurrent(PAGE_ADD_SC_FORM)
			page_instance.header.AddHistory(PAGE_ADD_SC_FORM)
		case "new_folder":
			p.createFolderModal.modal.SetVisible(true)
		case "delete_folder":
			p.deleteFolderConfirm.SetVisible(true)
		}
		p.tokenMenuSelect.SelectModal.Modal.SetVisible(false)
	}

	if p.deleteFolderConfirm.ClickedYes() {
		err := p.deleteCurrentFolder()
		if err != nil {
			notification_modals.ErrorInstance.SetText("Error", err.Error())
			notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
		} else {
			notification_modals.SuccessInstance.SetText("Success", lang.Translate("Folder deleted and all subfolders."))
			notification_modals.SuccessInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			p.changeFolder(p.currentFolder.ParentId)
		}
	}

	if p.buttonFolderGoBack.Clickable.Clicked() {
		p.changeFolder(p.currentFolder.ParentId)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Left: unit.Dp(30), Right: unit.Dp(30),
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
			widgets := []layout.ListElement{
				func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Spacer{Height: unit.Dp(10)}.Layout(gtx)
				},
			}

			if len(p.tokenItems) == 0 {
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

				itemIndex := 0
				for i := 0; i < len(p.tokenItems); i += 3 {
					widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
						dims := layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if itemIndex < len(p.tokenItems) {
									return p.tokenItems[itemIndex].Layout(gtx, th)
								}
								return layout.Dimensions{}
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if itemIndex+1 < len(p.tokenItems) {
									return p.tokenItems[itemIndex+1].Layout(gtx, th)
								}
								return layout.Dimensions{}
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(20)}.Layout),
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								if itemIndex+2 < len(p.tokenItems) {
									return p.tokenItems[itemIndex+2].Layout(gtx, th)
								}
								return layout.Dimensions{}
							}),
						)

						itemIndex += 3
						return dims
					})
				}
			}

			widgets = append(widgets, func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(20)}.Layout(gtx)
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

func (p *PageSCFolders) changeFolder(id sql.NullInt32) error {
	if id.Valid {
		tokenFolder, _ := wallet_manager.OpenedWallet.GetTokenFolder(id.Int32)
		p.currentFolder = tokenFolder
	} else {
		p.currentFolder = nil
	}

	return p.Load()
}

type CreateFolderModal struct {
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

func (c *CreateFolderModal) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, e := range c.txtFolderName.Editor.Events() {
		e, ok := e.(widget.SubmitEvent)
		if ok {
			wallet := wallet_manager.OpenedWallet

			tokenFolder := wallet_manager.TokenFolder{Name: e.Text}
			currentFolder := page_instance.pageSCFolders.currentFolder
			if currentFolder != nil {
				parentId := sql.NullInt32{Int32: currentFolder.ID, Valid: true}
				tokenFolder.ParentId = parentId
			}

			err := wallet.StoreFolderToken(&tokenFolder)
			if err != nil {
				notification_modals.ErrorInstance.SetText(lang.Translate("Error"), err.Error())
				notification_modals.ErrorInstance.SetVisible(true, notification_modals.CLOSE_AFTER_DEFAULT)
			} else {
				c.txtFolderName.SetValue("")
				c.modal.SetVisible(false)
				page_instance.pageSCFolders.Load()
				notification_modals.SuccessInstance.SetText(lang.Translate("Success"), lang.Translate("New folder created."))
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
	folderIcon *widget.Icon
	clickable  *widget.Clickable
	folder     wallet_manager.TokenFolder
}

func NewTokenFolderItem(folder wallet_manager.TokenFolder) *TokenFolderItem {
	folderIcon, _ := widget.NewIcon(icons.FileFolder)
	return &TokenFolderItem{
		folder:     folder,
		folderIcon: folderIcon,
		clickable:  new(widget.Clickable),
	}
}

func (item *TokenFolderItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if item.clickable.Hovered() {
		pointer.CursorPointer.Add(gtx.Ops)
	}

	if item.clickable.Clicked() {
		id := sql.NullInt32{Int32: item.folder.ID, Valid: true}
		page_instance.pageSCFolders.changeFolder(id)
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return item.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.Y = gtx.Constraints.Max.X
				paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 255, A: 255}, clip.UniformRRect(image.Rectangle{
					Max: gtx.Constraints.Max,
				}, gtx.Dp(10)).Op(gtx.Ops))

				return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return item.folderIcon.Layout(gtx, color.NRGBA{A: 255})
				})
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(5)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(16), item.folder.Name)
			lbl.Alignment = text.Middle
			lbl.Font.Weight = font.Bold
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(1)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(th, unit.Sp(12), "? tokens")
			lbl.Alignment = text.Middle
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
	)
}

type TokenMenuItem struct {
	Icon  *widget.Icon
	Title string
}

func (t TokenMenuItem) Layout(gtx layout.Context, index int, th *material.Theme) layout.Dimensions {
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

type TokenMenuSelect struct {
	SelectModal *prefabs.SelectModal
}

func NewTokenMenuSelect() *TokenMenuSelect {
	items := []*prefabs.SelectListItem{}

	contractIcon, _ := widget.NewIcon(icons.ActionNoteAdd)
	items = append(items, prefabs.NewSelectListItem("add_token", TokenMenuItem{
		Icon:  contractIcon,
		Title: "Add token", //@lang.Translate("Add token")
	}.Layout))

	folderIcon, _ := widget.NewIcon(icons.FileCreateNewFolder)
	items = append(items, prefabs.NewSelectListItem("new_folder", TokenMenuItem{
		Icon:  folderIcon,
		Title: "New folder", //@lang.Translate("New folder")
	}.Layout))

	listIcon, _ := widget.NewIcon(icons.ActionList)
	items = append(items, prefabs.NewSelectListItem("view_list", TokenMenuItem{
		Icon:  listIcon,
		Title: "View list", //@lang.Translate("View list")
	}.Layout))

	viewFolderIcon, _ := widget.NewIcon(icons.FileFolder)
	items = append(items, prefabs.NewSelectListItem("view_folder", TokenMenuItem{
		Icon:  viewFolderIcon,
		Title: "View folder", //@lang.Translate("View folder")
	}.Layout))

	deleteIcon, _ := widget.NewIcon(icons.ActionDelete)
	items = append(items, prefabs.NewSelectListItem("delete_folder", TokenMenuItem{
		Icon:  deleteIcon,
		Title: "Delete this folder", //@lang.Translate("Delete this folder")
	}.Layout))

	selectModal := prefabs.NewSelectModal()
	app_instance.Router.AddLayout(router.KeyLayout{
		DrawIndex: 1,
		Layout: func(gtx layout.Context, th *material.Theme) {
			selectModal.Layout(gtx, th, items)
		},
	})

	return &TokenMenuSelect{
		SelectModal: selectModal,
	}
}
