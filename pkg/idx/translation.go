package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type TranslationID string

func (id TranslationID) String() string {
	return string(id)
}

func (id TranslationID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "translation_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating translation id: expected prefix translation_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Page id",
		)
	}
	return nil
}

func NewTranslationID() (TranslationID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return TranslationID(fmt.Sprintf("translation_%s", random.String())), nil
}
