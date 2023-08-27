package dex_sc

import "github.com/g45t345rt/g45w/utils"

var DEX_SC_SHA256 = "51f330aeb991da9c845b77daf45304acf6e07a0e87125b8c0b8884dadbbd9dba"

// DEX Tokens
// DFRAX: f42fd725bc3659a7e6502ce416363afea0951e7f21af4f8f71b42090206e29d4
// DLINK: ab8ee3627b212a0b3803c127f3de7c44465fac21ec30692cb7988b14059990bb
// DUSDC: bc161c4f65285d5d927e9749fddbd127859748be7e161099f2f6785edc70b3dc
// DUSDT: f93b8d7fbbbf4e8f8a1e91b7ce21ac5d2b6aecc4de88cde8e929bce5f1746fbd
// DWBTC: b0bb9c1c75fc0e84dd92ce03f0619d1b61737981f0bb796911ea31529a76358c
// DWETH: fb855d8edd1d95ea94e9544224019c3fe4e636086f7266808879d6134ee2b8f1
// DgOHM: 92136ec02ca1e0db8e1767f7d5d221c7951263790fe4ee6616c4dd6c011e65ba
// DDAI: 	93707e89ba07f9aafc862ae07df1bfa70f488d5157d37439b85498fb79b6d1e6

type SC struct {
	SCID           string
	Name           string
	Decimals       uint64
	ImageUrl       string
	Symbol         string
	TotalSupply    uint64
	NativeSymbol   string
	NativeDecimals uint64
	Quorum         uint64
	NumTrustees    uint64
	Version        string
	BridgeOpen     bool
	BridgeFee      uint64
}

func (asset *SC) Parse(scId string, values map[string]interface{}) (err error) {
	asset.SCID = scId

	asset.Name, err = utils.DecodeString(values["name"].(string))
	if err != nil {
		return err
	}

	asset.Decimals = uint64(values["decimals"].(float64))

	asset.ImageUrl, err = utils.DecodeString(values["image_url"].(string))
	if err != nil {
		return err
	}

	asset.Symbol, err = utils.DecodeString(values["symbol"].(string))
	if err != nil {
		return err
	}

	asset.TotalSupply = uint64(values["totalsupply"].(float64))

	asset.NativeSymbol, err = utils.DecodeString(values["native_symbol"].(string))
	if err != nil {
		return err
	}

	asset.NativeDecimals = uint64(values["native_decimals"].(float64))
	asset.Quorum = uint64(values["quorum"].(float64))
	asset.NumTrustees = uint64(values["numTrustees"].(float64))

	asset.Version, err = utils.DecodeString(values["version"].(string))
	if err != nil {
		return err
	}

	asset.BridgeOpen = values["bridgeOpen"].(float64) != 0
	asset.BridgeFee = uint64(values["bridgeFee"].(float64))

	return nil
}
