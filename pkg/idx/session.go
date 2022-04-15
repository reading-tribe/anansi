package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type SessionID string

func (id SessionID) String() string {
	return string(id)
}

func (id SessionID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "sk_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating session id: expected prefix sk_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Session id",
		)
	}
	return nil
}

func NewSessionID() (SessionID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return SessionID(fmt.Sprintf("sk_%s", random.String())), nil
}
