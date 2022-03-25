package dbmodel

import "fmt"

type Language string

const (
	LanguageFrench  Language = "French"
	LanguageEnglish Language = "English"
	LanguageGerman  Language = "German"
)

func (l Language) Validate() error {
	valid := false

	if l == LanguageEnglish {
		valid = true
	}

	if l == LanguageFrench {
		valid = true
	}

	if l == LanguageGerman {
		valid = true
	}

	if !valid {
		return fmt.Errorf("currently unsupported language: %s", l)
	}

	return nil
}
