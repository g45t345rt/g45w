package theme

import (
	"image/color"

	"github.com/g45t345rt/g45w/components"
)

type Theme struct {
	Key                 string
	Name                string
	ThemeIndicatorColor color.NRGBA

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
	BottomBarWalletBgColor    color.NRGBA
	BottomBarWalletTextColor  color.NRGBA
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
}

var Current Theme

// don't use map[string] the ordering is not guaranteed
var Themes = []Theme{Light, Dark, Blue}

func Get(key string) *Theme {
	for _, theme := range Themes {
		if theme.Key == key {
			return &theme
		}
	}

	return nil
}

var whiteColor = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
var blackColor = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
