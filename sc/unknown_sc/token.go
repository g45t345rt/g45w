package unknown_sc

import "github.com/g45t345rt/g45w/utils"

type Token struct {
	SCID     string
	Name     string
	Decimals uint64
	ImageUrl string
	Symbol   string
}

func (token *Token) Parse(scId string, values map[string]interface{}) {
	token.SCID = scId

	name, ok := values["name"]
	if ok {
		token.Name, _ = utils.DecodeString(name.(string))
	}

	decimals, ok := values["decimals"]
	if ok {
		token.Decimals = uint64(decimals.(float64))
	}

	imageUrl, ok := values["image_url"]
	if ok {
		token.ImageUrl, _ = utils.DecodeString(imageUrl.(string))
	}

	image, ok := values["image"]
	if ok {
		token.ImageUrl, _ = utils.DecodeString(image.(string))
	}

	symbol, ok := values["symbol"]
	if ok {
		token.Symbol, _ = utils.DecodeString(symbol.(string))
	}

	return
}
