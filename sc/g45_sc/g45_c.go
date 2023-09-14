package g45_sc

import (
	"encoding/hex"
	"encoding/json"
	"regexp"

	"github.com/g45t345rt/g45w/utils"
)

var G45_C_SHA256 = "8729f4b6fe18509ddea9c54addf8513a729d6b382ed9eb50df27ea9548ef680c"

type G45_C struct {
	SCID           string
	FrozenAssets   bool
	FrozenMetadata bool
	Owner          string
	OriginalOwner  string
	Collection     string
	MetadataFormat string
	Metadata       string
	Assets         map[string]uint64
	AssetCount     uint64
	Timestamp      uint64
}

func (collection *G45_C) Parse(scId string, values map[string]interface{}) (err error) {
	collection.SCID = scId
	collection.FrozenAssets = values["frozenAssets"].(float64) != 0
	collection.FrozenMetadata = values["frozenMetadata"].(float64) != 0

	collection.MetadataFormat, err = utils.DecodeString(values["metadataFormat"].(string))
	if err != nil {
		return
	}

	collection.Metadata, err = utils.DecodeString(values["metadata"].(string))
	if err != nil {
		return
	}

	collection.Timestamp = uint64(values["timestamp"].(float64))

	collection.Owner, err = utils.DecodeAddress(values["owner"].(string))
	if err != nil {
		return
	}

	collection.OriginalOwner, err = utils.DecodeAddress(values["originalOwner"].(string))
	if err != nil {
		return
	}

	assetKey, err := regexp.Compile(`assets_(.+)`)
	if err != nil {
		return
	}

	collection.Assets = make(map[string]uint64)
	for sKey, sValue := range values {
		if assetKey.Match([]byte(sKey)) {
			valueBytes, err := hex.DecodeString(sValue.(string))
			if err != nil {
				return err
			}

			var assets map[string]uint64
			err = json.Unmarshal(valueBytes, &assets)
			if err != nil {
				return err
			}

			for key, value := range assets {
				collection.Assets[key] = value
			}
		}
	}

	collection.AssetCount = uint64(len(collection.Assets))
	return
}
