package theme

import (
	"image/color"

	"github.com/g45t345rt/g45w/components"
)

var Blue = &Theme{
	Key:            "blue",
	Name:           "Blue", //@lang.Translate("Blue")
	IndicatorColor: color.NRGBA{B: 150, A: 255},

	TextColor:            whiteColor,
	TextMuteColor:        color.NRGBA{R: 255, G: 255, B: 255, A: 150},
	DividerColor:         color.NRGBA{R: 255, G: 255, B: 255, A: 50},
	BgColor:              blueColor,
	BgGradientStartColor: color.NRGBA{R: 27, G: 107, B: 211, A: 255},
	BgGradientEndColor:   color.NRGBA{R: 16, G: 87, B: 181, A: 255},
	HideBalanceBgColor:   blueColor,

	HeaderBackButtonColors: components.ButtonColors{
		TextColor:      color.NRGBA{R: 255, G: 255, B: 255, A: 200},
		HoverTextColor: &whiteColor,
	},
	HeaderTopBgColor: color.NRGBA{R: 27, G: 107, B: 211, A: 255},

	BottomBarBgColor: blueColor,
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
		BackgroundColor: blueColor,
		TextColor:       whiteColor,
		BorderColor:     whiteColor,
		HintColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 150},
	},

	ButtonIconPrimaryColors: components.ButtonColors{
		TextColor: whiteColor,
	},
	ButtonPrimaryColors: components.ButtonColors{
		TextColor:       blueColor,
		BackgroundColor: whiteColor,
	},
	ButtonSecondaryColors: components.ButtonColors{
		TextColor:   whiteColor,
		BorderColor: whiteColor,
	},
	ButtonInvertColors: components.ButtonColors{
		TextColor:       whiteColor,
		BackgroundColor: blueColor,
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
		BackgroundColor: blueColor,
		BackdropColor:   &color.NRGBA{R: 20, G: 20, B: 50, A: 230},
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
		TextColor:       blueColor,
	},

	ProgressBarColors: components.ProgressBarColors{
		BackgroundColor: blueColor,
		IndicatorColor:  whiteColor,
	},

	ListTextColor:        whiteColor,
	ListBgColor:          color.NRGBA{R: 16, G: 87, B: 181, A: 255},
	ListItemHoverBgColor: color.NRGBA{R: 16, G: 87, B: 181, A: 255},
	ListScrollBarBgColor: whiteColor,
	ListItemTagBgColor:   blueColor,
	ListItemTagTextColor: whiteColor,

	SwitchColors: SwitchColors{
		Enabled:  whiteColor,
		Disabled: blueColor,
		Track:    color.NRGBA{A: 100},
	},
}
