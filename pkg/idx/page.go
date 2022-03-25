package idx

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewPageID() (string, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return fmt.Sprintf("page_%s", random.String()), nil
}
