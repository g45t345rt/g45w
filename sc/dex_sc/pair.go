package dex_sc

import "github.com/g45t345rt/g45w/utils"

type Pair struct {
	SCID              string
	NumTrustees       uint64
	Asset1            string
	Asset2            string
	Symbol            string
	Quorum            uint64
	Fee               uint64
	Liquidity1        uint64
	Liquidity2        uint64
	SharesOutstanding uint64
	AddCount          uint64 // liquidity add count
	RemoveCount       uint64
	SwapCount         uint64
}

func (pair *Pair) Parse(scId string, values map[string]interface{}) (err error) {
	pair.SCID = scId
	pair.NumTrustees = uint64(values["numTrustees"].(float64))

	pair.Asset1 = values["asset1"].(string)
	pair.Asset2 = values["asset2"].(string)

	pair.Symbol, err = utils.DecodeString(values["symbol"].(string))
	if err != nil {
		return
	}

	pair.Quorum = uint64(values["quorum"].(float64))
	pair.Fee = uint64(values["fee"].(float64))
	pair.Liquidity1 = uint64(values["val1"].(float64))
	pair.Liquidity2 = uint64(values["val2"].(float64))
	pair.SharesOutstanding = uint64(values["sharesOutstanding"].(float64))
	pair.AddCount = uint64(values["adds"].(float64))
	pair.RemoveCount = uint64(values["rems"].(float64))
	pair.SwapCount = uint64(values["swaps"].(float64))

	return
}
