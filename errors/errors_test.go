package errors

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

func TestCheckIsSpecificError(t *testing.T) {
	t.Run("API error", func(t *testing.T) {
		err := &APIError{
			StatusCode: http.StatusNotFound,
			Code:       "ResourceNotFound",
			Message:    "Resource Not Found", // note dosesn't matter
		}

		if !IsSpecificError(err, "ResourceNotFound") {
			t.Fatalf("Expected `ResourceNotFound`, got %v", err.Code)
		}

		if IsSpecificError(err, "IncorrectCode") {
			t.Fatalf("Expected `IncorrectCode`, got %v", err.Code)
		}
	})

	t.Run("Non Specific Error Type", func(t *testing.T) {
		err := errors.New("This is a new error")

		if IsSpecificError(err, "ResourceNotFound") {
			t.Fatalf("Specific Error Type Found")
		}
	})
}
