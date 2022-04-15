package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type DiversityAndInclusionID string

func (id DiversityAndInclusionID) String() string {
	return string(id)
}

func (id DiversityAndInclusionID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "dei_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating diversity and inclusion id: expected prefix dei_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Invalid diversity and inclusion id",
		)
	}
	return nil
}

func NewDiversityAndInclusionCatalogueID() (DiversityAndInclusionID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return DiversityAndInclusionID(fmt.Sprintf("dei_%s", random.String())), nil
}
