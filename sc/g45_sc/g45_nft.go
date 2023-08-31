package g45_sc

import (
	"github.com/g45t345rt/g45w/utils"
)

var G45_NFT_PRIVATE_SHA256 = "85ba16b77cbd4cfb3ba5690d74e849c2080905075bd87af29cbbb14f33bae2bd"
var G45_NFT_PUBLIC_SHA256 = ""

type G45_NFT struct {
	SCID string
	// Private        bool
	Minter         string
	MetadataFormat string
	Metadata       string
	Collection     string
	Owner          string
	Timestamp      uint64
}

func (asset *G45_NFT) Parse(scId string, values map[string]interface{}) (err error) {
	asset.SCID = scId
	asset.Timestamp = uint64(values["timestamp"].(float64))

	asset.Collection, err = utils.DecodeString(values["collection"].(string))
	if err != nil {
		return err
	}

	asset.MetadataFormat, err = utils.DecodeString(values["metadataFormat"].(string))
	if err != nil {
		return err
	}

	asset.Metadata, err = utils.DecodeString(values["metadata"].(string))
	if err != nil {
		return err
	}

	asset.Owner, err = utils.DecodeString(values["owner"].(string))
	if err != nil {
		return err
	}

	asset.Minter, err = utils.DecodeAddress(values["minter"].(string))
	if err != nil {
		return err
	}

	return nil
}
