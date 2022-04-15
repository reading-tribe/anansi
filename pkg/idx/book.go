package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type BookID string

func (id BookID) String() string {
	return string(id)
}

func (id BookID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "book_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating book id: expected prefix book_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Invalid book id",
		)
	}
	return nil
}

func NewBookID() (BookID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return BookID(fmt.Sprintf("book_%s", random.String())), nil
}
