package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type ChildProfileID string

func (id ChildProfileID) String() string {
	return string(id)
}

func (id ChildProfileID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "childprofile_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating book id: expected prefix childprofile_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Invalid child profile id",
		)
	}
	return nil
}

func NewChildProfile() (ChildProfileID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return ChildProfileID(fmt.Sprintf("childprofile_%s", random.String())), nil
}
