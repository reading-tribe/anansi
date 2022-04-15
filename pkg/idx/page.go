package idx

import (
	"fmt"
	"strings"

	"github.com/reading-tribe/anansi/pkg/errorx"
	"github.com/segmentio/ksuid"
)

type PageID string

func (id PageID) String() string {
	return string(id)
}

func (id PageID) Validate() *errorx.AnansiError {
	if !strings.HasPrefix(id.String(), "page_") {
		return errorx.NewError(
			fmt.Errorf("Precondition failure while validating page id: expected prefix page_ got: %s", id.String()),
			errorx.ValidationFailure,
			"Page id",
		)
	}
	return nil
}

func NewPageID() (PageID, error) {
	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return "", randomErr
	}
	return PageID(fmt.Sprintf("page_%s", random.String())), nil
}
