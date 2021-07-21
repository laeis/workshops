package errors

import (
	"github.com/pkg/errors"
)

// Error represents an error that could be wrapping another error, it includes a code for determining what
// triggered the error.
type Error struct {
	orig error
	msg  string
}

var BadRequest = errors.New("Bad request")

var NotFound = errors.New("not found")
