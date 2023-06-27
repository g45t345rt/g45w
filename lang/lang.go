package lang

import (
	"fmt"

	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/settings"
)

var SELECT_WALLET = Lang{key: "select_wallet"}

var defaultLang = "en"
var languages = []string{"en", "fr"}
var keyValues = make(map[string]map[string]string)

func Load() error {
	for _, lang := range languages {
		values, err := assets.GetLang(fmt.Sprintf("%s.json", lang))
		if err != nil {
			return err
		}
		keyValues[lang] = values
	}

	return nil
}

func GetValue(lang string, key string) string {
	values, ok := keyValues[lang]
	if !ok {
		if lang != defaultLang {
			return GetValue(defaultLang, key)
		} else {
			return "UNKNOWN_LANG_MAP"
		}
	}

	value, ok := values[key]
	if !ok {
		if lang != defaultLang {
			return GetValue(defaultLang, key)
		} else {
			return "UNKNOWN_KEY"
		}
	}

	return value
}

type Lang struct {
	key string
}

func (l Lang) String() string {
	lang := settings.Instance.Language
	return GetValue(lang, l.key)
}
