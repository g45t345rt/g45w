package theme

import (
	"image/color"

	"gioui.org/op/paint"
	"github.com/g45t345rt/g45w/assets"
)

var whiteColor = color.NRGBA{R: 250, G: 250, B: 250, A: 255}
var blackColor = color.NRGBA{R: 10, G: 10, B: 10, A: 255}
var blueColor = color.NRGBA{R: 2, G: 62, B: 138, A: 255}

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

	Light.TokenImage = opImgTokenBlack
	Dark.TokenImage = opImgTokenWhite
	Blue.TokenImage = opImgTokenWhite
}
