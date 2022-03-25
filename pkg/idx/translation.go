package idx

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewTranslationID() (string, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return fmt.Sprintf("translation_%s", random.String()), nil
}
