package idx

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewSessionID() (string, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return fmt.Sprintf("sk_%s", random.String()), nil
}
