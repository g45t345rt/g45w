package unknown_sc

import "github.com/g45t345rt/g45w/utils"

type SC struct {
	SCID     string
	Name     string
	Decimals uint64
	ImageUrl string
	Symbol   string
}

func (asset *SC) Parse(scId string, values map[string]interface{}) (err error) {
	asset.SCID = scId

	name, ok := values["name"]
	if ok {
		asset.Name, _ = utils.DecodeString(name.(string))
	}

	decimals, ok := values["decimals"]
	if ok {
		asset.Decimals = uint64(decimals.(float64))
	}

	imageUrl, ok := values["image_url"]
	if ok {
		asset.ImageUrl, _ = utils.DecodeString(imageUrl.(string))
	}

	image, ok := values["image"]
	if ok {
		asset.ImageUrl, _ = utils.DecodeString(image.(string))
	}

	symbol, ok := values["symbol"]
	if ok {
		asset.Symbol, _ = utils.DecodeString(symbol.(string))
	}

	return nil
}
