package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type UserID string

func (id UserID) String() string {
	return string(id)
}

func (id UserID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "user_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating user id: expected prefix user_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Invalid user id",
		)
	}
	return nil
}

func NewUserID() (UserID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return UserID(fmt.Sprintf("user_%s", random.String())), nil
}
