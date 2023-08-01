package lang

import (
	"fmt"
	"time"

	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/settings"
	"github.com/xeonx/timeago"
)

type Lang struct {
	Key     string
	Name    string
	ImgPath string
}

var SupportedLanguages = []Lang{
	{Key: "en", Name: "English", ImgPath: "lang/en.png"},    //@lang.Translate("English")
	{Key: "fr", Name: "French", ImgPath: "lang/fr.png"},     //@lang.Translate("French")
	{Key: "es", Name: "Spanish", ImgPath: "lang/es.png"},    //@lang.Translate("Spanish")
	{Key: "it", Name: "Italian", ImgPath: "lang/it.png"},    //@lang.Translate("Italian")
	{Key: "nl", Name: "Dutch", ImgPath: "lang/nl.png"},      //@lang.Translate("Dutch")
	{Key: "ru", Name: "Russian", ImgPath: "lang/ru.png"},    //@lang.Translate("Russian")
	{Key: "pt", Name: "Portuguese", ImgPath: "lang/pt.png"}, //@lang.Translate("Portuguese")
	{Key: "ro", Name: "Romanian", ImgPath: "lang/ro.png"},   //@lang.Translate("Romanian")
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

	// TODO: add more predefined configuration language
	max := 100 * 365 * 24 * time.Hour
	timeago.English.Max = max
	timeago.Chinese.Max = max
	timeago.German.Max = max
	timeago.Korean.Max = max
	timeago.French.Max = max
	timeago.Portuguese.Max = max
	timeago.Spanish.Max = max
	timeago.Turkish.Max = max

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
