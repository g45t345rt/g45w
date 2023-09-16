package dex_sc

import (
	"github.com/g45t345rt/g45w/utils"
)

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
	RemoveCount       uint64 // liquidity remove count
	SwapCount         uint64
}

func (pair *Pair) CalcShare(share uint64, reverse bool) (value uint64) {
	if reverse {
		value = utils.MultDiv(pair.Liquidity2, share, pair.SharesOutstanding)
	} else {
		value = utils.MultDiv(pair.Liquidity1, share, pair.SharesOutstanding)
	}
	return
}

func (pair *Pair) CalcOwnership(share uint64) float32 {
	if pair.SharesOutstanding == 0 {
		return 0
	}

	return float32(share) / float32(pair.SharesOutstanding) * 100.0
}

func (pair *Pair) CalcSwap(amt uint64, reverse bool) (receive uint64, fee uint64, slip float64) {
	if amt == 0 {
		return
	}

	if reverse {
		receiveAmt := float64(amt) * float64(pair.Liquidity1) / float64(pair.Liquidity2+amt)
		receiveAmtMinusFee := receiveAmt * float64(10000-pair.Fee) / float64(10000)
		receive = uint64(receiveAmtMinusFee)
		fee = uint64(receiveAmt) - uint64(receiveAmtMinusFee)
		if pair.Liquidity2 != 0 {
			slip = 100.0 - (1.0 / (1.0 + float64(amt)/float64(pair.Liquidity2)) * 100.0)
		}
	} else {
		receiveAmt := float64(amt) * float64(pair.Liquidity2) / float64(pair.Liquidity1+amt)
		receiveAmtMinusFee := receiveAmt * float64(10000-pair.Fee) / float64(10000)
		receive = uint64(receiveAmtMinusFee)
		fee = uint64(receiveAmt) - uint64(receiveAmtMinusFee)
		if pair.Liquidity1 != 0 {
			slip = 100.0 - (1.0 / (1.0 + float64(amt)/float64(pair.Liquidity1)) * 100.0)
		}
	}

	return
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
