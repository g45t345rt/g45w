package lang

import (
	"fmt"

	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/settings"
)

type Lang struct {
	Key     string
	Name    string
	ImgPath string
}

var SupportedLanguages = []Lang{
	{Key: "en", Name: "English", ImgPath: "lang/en.png"}, //@lang.Translate("English")
	{Key: "fr", Name: "French", ImgPath: "lang/fr.png"},  //@lang.Translate("French")
	{Key: "es", Name: "Spanish", ImgPath: "lang/es.png"}, //@lang.Translate("Spanish")
}

var langValues = make(map[string]map[string]string)

func Load() error {
	for _, lang := range SupportedLanguages {
		if lang.Key == "en" {
			continue
		}

		values, err := assets.GetLang(fmt.Sprintf("%s.json", lang.Key))
		if err != nil {
			return err
		}
		langValues[lang.Key] = values
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
