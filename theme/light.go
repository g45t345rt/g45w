package theme

import (
	"image/color"

	"github.com/g45t345rt/g45w/components"
)

var Light = Theme{
	Key:                 "light",
	Name:                "Light", //@lang.Translate("Light")
	ThemeIndicatorColor: color.NRGBA{R: 255, G: 255, B: 255, A: 50},

	TextColor:            blackColor,
	TextMuteColor:        color.NRGBA{A: 150},
	DividerColor:         color.NRGBA{A: 50},
	BgColor:              whiteColor,
	BgGradientStartColor: color.NRGBA{R: 250, G: 250, B: 250, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 210, G: 210, B: 210, A: 255},
	HideBalanceBgColor:   color.NRGBA{R: 200, G: 200, B: 200, A: 255},

	HeaderBackButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &color.NRGBA{A: 255},
	},
	HeaderTopBgColor: color.NRGBA{R: 250, G: 250, B: 250, A: 255},

	BottomBarBgColor:         whiteColor,
	BottomBarWalletBgColor:   blackColor,
	BottomBarWalletTextColor: whiteColor,
	BottomButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &blackColor,
	},
	BottomButtonSelectedColor: blackColor,

	NodeStatusBgColor:        blackColor,
	NodeStatusTextColor:      whiteColor,
	NodeStatusDotGreenColor:  color.NRGBA{R: 0, G: 225, B: 0, A: 255},
	NodeStatusDotYellowColor: color.NRGBA{R: 255, G: 255, B: 0, A: 255},
	NodeStatusDotRedColor:    color.NRGBA{R: 225, G: 0, B: 0, A: 255},

	InputColors: components.InputColors{
		BackgroundColor: whiteColor,
		TextColor:       blackColor,
		BorderColor:     blackColor,
	},

	ButtonIconPrimaryColors: components.ButtonColors{
		TextColor: blackColor,
	},
	ButtonPrimaryColors: components.ButtonColors{
		TextColor:       whiteColor,
		BackgroundColor: blackColor,
	},
	ButtonSecondaryColors: components.ButtonColors{
		TextColor:   blackColor,
		BorderColor: blackColor,
	},
	ButtonInvertColors: components.ButtonColors{
		TextColor:       blackColor,
		BackgroundColor: whiteColor,
	},
	ButtonDangerColors: components.ButtonColors{
		TextColor:       whiteColor,
		BackgroundColor: color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	},

	TabBarsColors: components.TabBarsColors{
		InactiveColor: blackColor,
		ActiveColor:   blackColor,
	},

	ModalColors: components.ModalColors{
		BackgroundColor: whiteColor,
		BackdropColor:   &color.NRGBA{A: 100},
	},
	ModalButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{A: 100},
		HoverTextColor: &blackColor,
	},

	NotificationSuccessColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 0, G: 225, B: 0, A: 255},
		TextColor:       whiteColor,
	},
	NotificationErrorColors: components.NotificationColors{
		BackgroundColor: color.NRGBA{R: 225, G: 0, B: 0, A: 255},
		TextColor:       whiteColor,
	},
	NotificationInfoColors: components.NotificationColors{
		BackgroundColor: whiteColor,
		TextColor:       blackColor,
	},

	ProgressBarColors: components.ProgressBarColors{
		BackgroundColor: whiteColor,
		IndicatorColor:  blackColor,
	},

	ListTextColor:        blackColor,
	ListBgColor:          whiteColor,
	ListItemHoverBgColor: color.NRGBA{A: 100},
	ListScrollBarBgColor: blackColor,
	ListItemTagBgColor:   color.NRGBA{R: 225, G: 225, B: 225, A: 255},
	ListItemTagTextColor: blackColor,
}
