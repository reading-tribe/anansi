package idx

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func NewDiversityAndInclusionCatalogueID() (string, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return fmt.Sprintf("dei_%s", random.String()), nil
}
