package triton

import (
	"fmt"
	"github.com/hashicorp/errwrap"
)

// TritonError represents an error code and message along with
// the status code of the HTTP request which resulted in the
// error message.
type TritonError struct {
	StatusCode int
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Error implements interface Error on the TritonError type.
func (e TritonError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// IsResourceNotFound returns true if the error code is ResourceNotFound,
// and returns false otherwise.
func IsResourceNotFound(err error) bool {
	tritonErrorInterface := errwrap.GetType(err.(error), &TritonError{})
	if tritonErrorInterface == nil {
		return false
	}

	tritonErr := tritonErrorInterface.(*TritonError)
	if tritonErr.Code == "ResourceNotFound" {
		return true
	}

	return false
}

