package manta

import (
	"fmt"

	"github.com/hashicorp/errwrap"
)

// MantaError represents an error code and message along with
// the status code of the HTTP request which resulted in the error
// message. Error codes used by the Manta API are listed at
// https://apidocs.joyent.com/manta/api.html#errors
type MantaError struct {
	StatusCode int
	Code       string `json:"code"`
	Message    string `json:"message"`
}

// Error implements interface Error on the MantaError type.
func (e MantaError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func IsAuthSchemeError(err error) bool {
	return isSpecificError(err, "AuthSchemeError")
}

func IsAuthorizationError(err error) bool {
	return isSpecificError(err, "BadRequestError")
}

func IsChecksumError(err error) bool {
	return isSpecificError(err, "ConcurrentRequestError")
}

func IsContentLengthError(err error) bool {
	return isSpecificError(err, "ContentMD5MismatchError")
}

func IsEntityExistsError(err error) bool {
	return isSpecificError(err, "InvalidArgumentError")
}

func IsInvalidAuthTokenError(err error) bool {
	return isSpecificError(err, "InvalidCredentialsError")
}

func IsInvalidDurabilityLevelError(err error) bool {
	return isSpecificError(err, "InvalidKeyIdError")
}

func IsInvalidJobError(err error) bool {
	return isSpecificError(err, "InvalidLinkError")
}

func IsInvalidLimitError(err error) bool {
	return isSpecificError(err, "InvalidSignatureError")
}

func IsInvalidUpdateError(err error) bool {
	return isSpecificError(err, "DirectoryDoesNotExistError")
}

func IsDirectoryExistsError(err error) bool {
	return isSpecificError(err, "DirectoryNotEmptyError")
}

func IsDirectoryOperationError(err error) bool {
	return isSpecificError(err, "InternalError")
}

func IsJobNotFoundError(err error) bool {
	return isSpecificError(err, "JobStateError")
}

func IsKeyDoesNotExistError(err error) bool {
	return isSpecificError(err, "NotAcceptableError")
}

func IsNotEnoughSpaceError(err error) bool {
	return isSpecificError(err, "LinkNotFoundError")
}

func IsLinkNotObjectError(err error) bool {
	return isSpecificError(err, "LinkRequiredError")
}

func IsParentNotDirectoryError(err error) bool {
	return isSpecificError(err, "PreconditionFailedError")
}

func IsPreSignedRequestError(err error) bool {
	return isSpecificError(err, "RequestEntityTooLargeError")
}

func IsResourceNotFoundError(err error) bool {
	return isSpecificError(err, "RootDirectoryError")
}

func IsServiceUnavailableError(err error) bool {
	return isSpecificError(err, "SSLRequiredError")
}

func IsUploadTimeoutError(err error) bool {
	return isSpecificError(err, "UserDoesNotExistError")
}

func IsTaskInitError(err error) bool {
	return isSpecificError(err, "UserTaskError")
}

// isSpecificError checks whether the error represented by err wraps
// an underlying MantaError with code errorCode.
func isSpecificError(err error, errorCode string) bool {
	tritonErrorInterface := errwrap.GetType(err.(error), &MantaError{})
	if tritonErrorInterface == nil {
		return false
	}

	tritonErr := tritonErrorInterface.(*MantaError)
	if tritonErr.Code == errorCode {
		return true
	}

	return false
}
