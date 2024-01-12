package theme

import (
	"image/color"

	"gioui.org/op/paint"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/components"
)

type Theme struct {
	Key            string
	Name           string
	IndicatorColor color.NRGBA

	TextColor            color.NRGBA
	TextMuteColor        color.NRGBA
	DividerColor         color.NRGBA
	BgColor              color.NRGBA
	BgGradientStartColor color.NRGBA
	BgGradientEndColor   color.NRGBA
	HideBalanceBgColor   color.NRGBA

	// Header
	HeaderBackButtonColors components.ButtonColors
	HeaderTopBgColor       color.NRGBA

	// Bottom Bar
	BottomBarBgColor          color.NRGBA
	BottomButtonColors        components.ButtonColors
	BottomButtonSelectedColor color.NRGBA

	// Node Status
	NodeStatusBgColor        color.NRGBA
	NodeStatusTextColor      color.NRGBA
	NodeStatusDotGreenColor  color.NRGBA
	NodeStatusDotYellowColor color.NRGBA
	NodeStatusDotRedColor    color.NRGBA

	// Input
	InputColors components.InputColors

	// Button
	ButtonIconPrimaryColors components.ButtonColors
	ButtonPrimaryColors     components.ButtonColors
	ButtonSecondaryColors   components.ButtonColors
	ButtonInvertColors      components.ButtonColors
	ButtonDangerColors      components.ButtonColors

	// Tab Bars
	TabBarsColors components.TabBarsColors

	// Modal
	ModalColors       components.ModalColors
	ModalButtonColors components.ButtonColors

	// Notifications
	NotificationSuccessColors components.NotificationColors
	NotificationErrorColors   components.NotificationColors
	NotificationInfoColors    components.NotificationColors

	// Progress Bar
	ProgressBarColors components.ProgressBarColors

	// List
	ListTextColor        color.NRGBA
	ListBgColor          color.NRGBA
	ListItemHoverBgColor color.NRGBA
	ListScrollBarBgColor color.NRGBA
	ListItemTagBgColor   color.NRGBA
	ListItemTagTextColor color.NRGBA
	//ListItemsColors      components.ListItemsColors

	// Switch
	SwitchColors SwitchColors

	// Images
	ArrowDownArcImage paint.ImageOp
	ArrowUpArcImage   paint.ImageOp
	CoinbaseImage     paint.ImageOp
	TokenImage        paint.ImageOp
	ManageFilesImage  paint.ImageOp
}

type SwitchColors struct {
	Enabled  color.NRGBA
	Disabled color.NRGBA
	Track    color.NRGBA
}

// default to Light theme (avoid nil pointer in FrameEvent before settings.Load() is set)
// settings.Load() will overwrite theme.Current with system pref or settings.json theme value
var Current *Theme = Light

// don't use map[string] the ordering is not guaranteed
var Themes = []*Theme{Light, Dark, Blue}

func Get(key string) *Theme {
	for _, theme := range Themes {
		if theme.Key == key {
			return theme
		}
	}

	return nil
}

func LoadImages() {
	imgArrowUpArcBlack, _ := assets.GetImage("arrow_up_arc.png")
	opImgArrowUpArcBlack := paint.NewImageOp(imgArrowUpArcBlack)
	imgArrowUpArcWhite, _ := assets.GetImage("arrow_up_arc_white.png")
	opImgArrowUpArcWhite := paint.NewImageOp(imgArrowUpArcWhite)

	Light.ArrowUpArcImage = opImgArrowUpArcBlack
	Dark.ArrowUpArcImage = opImgArrowUpArcWhite
	Blue.ArrowUpArcImage = opImgArrowUpArcWhite

	imgArrowDownArcBlack, _ := assets.GetImage("arrow_down_arc.png")
	opImgArrowDownArcBlack := paint.NewImageOp(imgArrowDownArcBlack)
	imgArrowDownArcWhite, _ := assets.GetImage("arrow_down_arc_white.png")
	opImgArrowDownArcWhite := paint.NewImageOp(imgArrowDownArcWhite)

	Light.ArrowDownArcImage = opImgArrowDownArcBlack
	Dark.ArrowDownArcImage = opImgArrowDownArcWhite
	Blue.ArrowDownArcImage = opImgArrowDownArcWhite

	imgCoinbaseBlack, _ := assets.GetImage("coinbase.png")
	opImgCoinbaseBlack := paint.NewImageOp(imgCoinbaseBlack)
	imgCoinbaseWhite, _ := assets.GetImage("coinbase_white.png")
	opImgCoinbaseWhite := paint.NewImageOp(imgCoinbaseWhite)

	Light.CoinbaseImage = opImgCoinbaseBlack
	Dark.CoinbaseImage = opImgCoinbaseWhite
	Blue.CoinbaseImage = opImgCoinbaseWhite

	imgTokenBlack, _ := assets.GetImage("token.png")
	opImgTokenBlack := paint.NewImageOp(imgTokenBlack)
	imgTokenWhite, _ := assets.GetImage("token_white.png")
	opImgTokenWhite := paint.NewImageOp(imgTokenWhite)

	imgManageFilesBlack, _ := assets.GetImage("manage_files.png")
	opImgManageFilesBlack := paint.NewImageOp(imgManageFilesBlack)
	imgManageFilesWhite, _ := assets.GetImage("manage_files_white.png")
	opImgManageFilesWhite := paint.NewImageOp(imgManageFilesWhite)

	Light.ManageFilesImage = opImgManageFilesBlack
	Dark.ManageFilesImage = opImgManageFilesWhite
	Blue.ManageFilesImage = opImgManageFilesWhite

	Light.TokenImage = opImgTokenBlack
	Dark.TokenImage = opImgTokenWhite
	Blue.TokenImage = opImgTokenWhite
}
