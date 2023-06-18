package g45_sc

import "encoding/json"

type TokenMetadata struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Image  string `json:"image"`
}

func (m *TokenMetadata) Parse(metadata string) (err error) {
	return json.Unmarshal([]byte(metadata), &m)
}

type NFTMetadata struct {
	ID          uint64                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Attributes  map[string]interface{} `json:"attributes"`

	Image  string                 `json:"image"`
	Images map[string]interface{} `json:"images"`

	Video  string                 `json:"video"`
	Videos map[string]interface{} `json:"videos"`

	Audio  string                 `json:"audio"`
	Audios map[string]interface{} `json:"audios"`
}

func (m *NFTMetadata) Parse(metadata string) (err error) {
	return json.Unmarshal([]byte(metadata), &m)
}
