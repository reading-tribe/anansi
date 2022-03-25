package idx

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewBookID() (string, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return fmt.Sprintf("book_%s", random.String()), nil
}
