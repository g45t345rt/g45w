package g45_sc

import (
	"regexp"

	"github.com/g45t345rt/g45w/utils"
)

var G45_AT_PRIVATE_SHA256 = "9e75a99d059f09ca4dec1fb8a591ecf920654c721251d2de5e4af147f863612f"
var G45_AT_PUBLIC_SHA256 = ""

type G45_AT struct {
	SCID string
	// Private          bool
	Minter           string
	OriginalMinter   string
	FrozenMetadata   bool
	FrozenMint       bool
	FrozenCollection bool
	MetadataFormat   string
	Metadata         string
	MaxSupply        uint64
	TotalSupply      uint64
	Decimals         uint64
	Collection       string
	Owners           map[string]uint64
	Timestamp        uint64
}

func (asset *G45_AT) Parse(scId string, values map[string]interface{}) (err error) {
	asset.SCID = scId
	asset.Timestamp = uint64(values["timestamp"].(float64))

	asset.Collection, err = utils.DecodeString(values["collection"].(string))
	if err != nil {
		return err
	}

	asset.FrozenMetadata = values["frozenMetadata"].(float64) != 0
	asset.FrozenMint = values["frozenMint"].(float64) != 0
	asset.FrozenCollection = values["frozenCollection"].(float64) != 0

	asset.MetadataFormat, err = utils.DecodeString(values["metadataFormat"].(string))
	if err != nil {
		return err
	}

	asset.Metadata, err = utils.DecodeString(values["metadata"].(string))
	if err != nil {
		return err
	}

	asset.MaxSupply = uint64(values["maxSupply"].(float64))
	asset.TotalSupply = uint64(values["totalSupply"].(float64))
	asset.Decimals = uint64(values["decimals"].(float64))

	asset.Minter, err = utils.DecodeAddress(values["minter"].(string))
	if err != nil {
		return err
	}

	asset.OriginalMinter, err = utils.DecodeAddress(values["originalMinter"].(string))
	if err != nil {
		return err
	}

	ownerKey, _ := regexp.Compile(`owner_(.+)`)
	asset.Owners = make(map[string]uint64)
	for key, value := range values {
		if ownerKey.Match([]byte(key)) {
			owner := ownerKey.ReplaceAllString(key, "$1")
			asset.Owners[owner] = uint64(value.(float64))
		}
	}
	return nil
}
