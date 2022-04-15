package repository

import (
	"github.com/reading-tribe/anansi/pkg/dbmodel"
)

type LanguageRepository interface {
	ListLanguages() ([]dbmodel.Language, error)
}

type languageRepository struct{}

func NewLanguageRepository() LanguageRepository {
	return languageRepository{}
}

func (l languageRepository) ListLanguages() ([]dbmodel.Language, error) {
	return []dbmodel.Language{
		dbmodel.LanguageEnglish,
		dbmodel.LanguageFrench,
		dbmodel.LanguageGerman,
	}, nil
}
