package triton

import (
	"errors"
	"strings"
)

// wrappedError is an implementation of error that has both the
// outer and inner errors.
type wrappedError struct {
	Outer error
	Inner error
}

// Error returns our outer error as a string by proxying straight to the outer
// object's own Error function.
func (w *wrappedError) Error() string {
	return w.Outer.Error()
}

// WrappedErrors returns an array of wrapped errors set on our wrappedError
// object.
func (w *wrappedError) WrappedErrors() []error {
	return []error{w.Outer, w.Inner}
}

// WrapError defines that outer wraps inner, returning an error type that can be
// cleanly used with the other methods in this package, such as Contains,
// GetAll, etc.
//
// This function won't modify the error message at all (the outer message will
// be used).
func WrapError(outer, inner error) error {
	return &wrappedError{
		Outer: outer,
		Inner: inner,
	}
}

// WrapErrorf wraps an error with a formatting message. This is similar to using
// `fmt.Errorf` to wrap an error. If you're using `fmt.Errorf` to wrap errors,
// you should replace it with this.
//
// format is the format of the error message. The string '{{err}}' will be
// replaced with the original error message.
func WrapErrorf(format string, err error) error {
	outerMsg := "<nil>"
	if err != nil {
		outerMsg = err.Error()
	}

	outer := errors.New(strings.Replace(
		format, "{{err}}", outerMsg, -1))

	return WrapError(outer, err)
}
