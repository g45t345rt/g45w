package sc

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/sc/g45_sc"
)

type SCType string

var (
	G45_NFT_TYPE SCType = "G45_NFT"
	G45_AT_TYPE  SCType = "G45_AT"
	G45_FAT_TYPE SCType = "G45_FAT"
	G45_C_TYPE   SCType = "G45_C"

	DEX_SC_TYPE SCType = "DEX_SC"

	UNKNOWN_TYPE SCType = "UNKNOWN"
)

func CheckType(code string) SCType {
	hash := sha256.New()
	hash.Write([]byte(code))
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)

	switch hashString {
	case g45_sc.G45_NFT_PUBLIC_SHA256, g45_sc.G45_NFT_PRIVATE_SHA256:
		return G45_NFT_TYPE
	case g45_sc.G45_FAT_PUBLIC_SHA256, g45_sc.G45_FAT_PRIVATE_SHA256:
		return G45_FAT_TYPE
	case g45_sc.G45_AT_PUBLIC_SHA256, g45_sc.G45_AT_PRIVATE_SHA256:
		return G45_AT_TYPE
	case g45_sc.G45_C_SHA256:
		return G45_C_TYPE
	case dex_sc.DEX_SC_SHA256:
		return DEX_SC_TYPE
	}

	return UNKNOWN_TYPE
}
