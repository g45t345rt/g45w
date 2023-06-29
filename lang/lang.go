package lang

import (
	"fmt"

	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/settings"
)

var SupportedLanguages = []string{"en", "fr", "es"}

var langValues = make(map[string]map[string]string)

func Load() error {
	for _, lang := range SupportedLanguages {
		if lang == "en" {
			continue
		}

		values, err := assets.GetLang(fmt.Sprintf("%s.json", lang))
		if err != nil {
			return err
		}
		langValues[lang] = values
	}

	return nil
}

func Translate(eng string) string {
	lang := settings.App.Language
	values, ok := langValues[lang]
	if !ok {
		return eng
	}

	value, ok := values[eng]
	if !ok {
		return eng
	}

	return value
}
