package theme

import (
	"image/color"

	"github.com/g45t345rt/g45w/components"
)

var Dark = &Theme{
	Key:                 "dark",
	Name:                "Dark", //@lang.Translate("Dark")
	ThemeIndicatorColor: color.NRGBA{A: 255},

	TextColor:            whiteColor,
	TextMuteColor:        color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	DividerColor:         color.NRGBA{R: 255, G: 255, B: 255, A: 25},
	BgColor:              blackColor,
	BgGradientStartColor: color.NRGBA{R: 30, G: 30, B: 30, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	HideBalanceBgColor:   color.NRGBA{A: 255},

	HeaderBackButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &whiteColor,
	},
	HeaderTopBgColor: color.NRGBA{R: 30, G: 30, B: 30, A: 255},

	BottomBarBgColor:         blackColor,
	BottomBarWalletBgColor:   whiteColor,
	BottomBarWalletTextColor: blackColor,
	BottomButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &whiteColor,
	},
	BottomButtonSelectedColor: whiteColor,

	NodeStatusBgColor:        color.NRGBA{A: 255},
	NodeStatusTextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	NodeStatusDotGreenColor:  color.NRGBA{R: 0, G: 200, B: 0, A: 255},
	NodeStatusDotYellowColor: color.NRGBA{R: 255, G: 255, B: 0, A: 255},
	NodeStatusDotRedColor:    color.NRGBA{R: 200, G: 0, B: 0, A: 255},

	InputColors: components.InputColors{
		BackgroundColor: blackColor,
		TextColor:       whiteColor,
		BorderColor:     whiteColor,
		HintColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	},

	ButtonIconPrimaryColors: components.ButtonColors{
		TextColor: whiteColor,
	},
	ButtonPrimaryColors: components.ButtonColors{
		TextColor:       blackColor,
		BackgroundColor: whiteColor,
	},
	ButtonSecondaryColors: components.ButtonColors{
		TextColor:   whiteColor,
		BorderColor: whiteColor,
	},
	ButtonInvertColors: components.ButtonColors{
		TextColor:       whiteColor,
		BackgroundColor: blackColor,
	},
	ButtonDangerColors: components.ButtonColors{
		TextColor:       whiteColor,
		BackgroundColor: color.NRGBA{R: 200, G: 0, B: 0, A: 255},
	},

	TabBarsColors: components.TabBarsColors{
		InactiveColor: whiteColor,
		ActiveColor:   whiteColor,
	},

	ModalColors: components.ModalColors{
		BackgroundColor: blackColor,
		BackdropColor:   &color.NRGBA{R: 20, G: 20, B: 20, A: 230},
	},
	ModalButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 100},
		HoverTextColor: &whiteColor,
	},

	NotificationSuccessColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 0, G: 200, B: 0, A: 255},
		TextColor:       whiteColor,
	},
	NotificationErrorColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 200, G: 0, B: 0, A: 255},
		TextColor:       whiteColor,
	},
	NotificationInfoColors: components.NotificationColors{
		BackgroundColor: whiteColor,
		TextColor:       blackColor,
	},

	ProgressBarColors: components.ProgressBarColors{
		BackgroundColor: blackColor,
		IndicatorColor:  whiteColor,
	},

	ListTextColor:        whiteColor,
	ListBgColor:          color.NRGBA{R: 15, G: 15, B: 15, A: 255},
	ListItemHoverBgColor: color.NRGBA{A: 100},
	ListScrollBarBgColor: whiteColor,
	ListItemTagBgColor:   blackColor,
	ListItemTagTextColor: whiteColor,
}
