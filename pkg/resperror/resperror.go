package resperror

import "fmt"

// General consts.
const (
	ErrValidation           = "Validation Error"
	ErrSystem               = "System Error"
	ErrUnknownValidationErr = "Unknown validation error"
	ErrNotFoundTitle        = "Not Found"
	ErrNotFoundDetail       = "Resource Not Found"
	ErrJSONSyntax           = "JSON Syntax Error"
)

// General error detail postfixes/prefixed.
const (
	// Postfix
	ErrIsRequired             = " is required"
	ErrIsAnInvalidUUID4       = " is an invalid UUID4"
	ErrInvalidTimestampFormat = " timestamp has an invalid format - Must be formatted in RFC3339 - i.e. 2006-01-02T15:04:05Z"
	ErrInvalidFloatFormat     = " has an invalid format - Must be 2 decimal places - i.e. 99.99"
	ErrIsInvalidInteger       = " must be a valid integer"

	// Prefix
	ErrIsInvalid = "Invalid value for "
)

// Error codes.
const (
	// ErrCodeSystem - For a system error code.
	ErrCodeSystem = 1

	// ErrCodeNotFound - For a resource not found error code.
	ErrCodeNotFound = 2

	// ErrorCodeValidation - For an unknown validation code.
	ErrCodeValidation = 100

	// General validation error codes.
	ErrCodeRequired             = 101
	ErrCodeInvalidFormat        = 102
	ErrCodeBadFormat            = 103
	ErrCodeBadTimestampFormat   = 104
	ErrCodeBadUUIDFormat        = 105
	ErrCodeBadFloatFormat       = 106
	ErrCodeInvalidIntegerFormat = 107

	// Merchant codes.
	ErrCodeInvalidCountry                     = 301
	ErrCodeInvalidTimezone                    = 302
	ErrCodeDuplicateClientRef                 = 303
	ErrCodeTerminatedMerchantCannotBeModified = 304
)

// IsValidationErr -
func IsValidationErr(code int) bool {
	// Currently anything above 100 is a validation error.
	return code >= 100
}

// TODO: Move error message details into consts.

// Data -
type Data struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Implement the error interface for Data.
func (d *Data) Error() string {
	return fmt.Sprintf("Code: %d, Title: %s, Detail: %s", d.Code, d.Title, d.Detail)
}

// Response -
type Response struct {
	Error *Data `json:"error"`
}

// SystemErr is a helper function for constructing a general system
// error from a detail string
func SystemErr(detail string) *Data {
	return &Data{
		Code:   ErrCodeSystem,
		Title:  ErrSystem,
		Detail: detail,
	}
}

// ValidationErr is a helper function for constructing a general validation
// error from a detail string.
func ValidationErr(detail string) *Data {
	return &Data{
		Code:   ErrCodeBadFormat,
		Title:  ErrValidation,
		Detail: detail,
	}
}

// ValidationRequired is a helper function for constructing a validation
// error for a required field.
func ValidationRequired(field string) *Data {
	return &Data{
		Code:   ErrCodeRequired,
		Title:  ErrValidation,
		Detail: field + ErrIsRequired,
	}
}

// ValidationInvalid is a helper function for constructing a validation
// error for an invalid field that doesn't have specific details.
func ValidationInvalid(field string) *Data {
	return &Data{
		Code:   ErrCodeInvalidFormat,
		Title:  ErrValidation,
		Detail: ErrIsInvalid + field,
	}
}

// ValidationInvalidIntegerFormat is a helper function for constructing a validation
// error for an invalid integer field.
func ValidationInvalidIntegerFormat(field string) *Data {
	return &Data{
		Code:   ErrCodeInvalidIntegerFormat,
		Title:  ErrValidation,
		Detail: field + ErrIsInvalidInteger,
	}
}

// TODO: Create ValidationInvalidWithDetails helper function?

// ValidationInvalidUUID4 is a helper function for constructing a validation
// error for a field that has an invalid UUID4 format.
func ValidationInvalidUUID4(field string) *Data {
	return &Data{
		Code:   ErrCodeBadUUIDFormat,
		Title:  ErrValidation,
		Detail: field + ErrIsAnInvalidUUID4,
	}
}

// ValidationTimestampFormat is a helper function for constructing a validation
// error for a field that has an invalid timestamp format.
func ValidationTimestampFormat(field string) *Data {
	return &Data{
		Code:   ErrCodeBadTimestampFormat,
		Title:  ErrValidation,
		Detail: field + ErrInvalidTimestampFormat,
	}
}

// ValidationFloatFormat is a helper function for constructing a validation
// error for a field that has an invalid floating point format.
func ValidationFloatFormat(field string) *Data {
	return &Data{
		Code:   ErrCodeBadFloatFormat,
		Title:  ErrValidation,
		Detail: field + ErrInvalidFloatFormat,
	}
}

// ValidationJSONSyntax is a helper function for constructing a validation
// error when json syntax is invalid.
func ValidationJSONSyntax(offset int64) *Data {
	return &Data{
		Code:   ErrCodeValidation,
		Title:  ErrJSONSyntax,
		Detail: fmt.Sprintf("%s at offset %d from start", ErrJSONSyntax, offset),
	}
}

// ErrorNotFound -
var ErrorNotFound = &Data{
	Code:   ErrCodeNotFound,
	Title:  ErrNotFoundTitle,
	Detail: ErrNotFoundDetail,
}

// ErrorUnknownValidation -
var ErrorUnknownValidation = &Data{
	Code:   ErrCodeValidation,
	Title:  ErrValidation,
	Detail: ErrUnknownValidationErr,
}

// ErrorInvalidCountry - Merchant
var ErrorInvalidCountry = &Data{
	Code:   ErrCodeInvalidCountry,
	Title:  ErrValidation,
	Detail: "Field country value is not present in list of available countries",
}

// ErrorInvalidTimezone - Merchant
var ErrorInvalidTimezone = &Data{
	Code:   ErrCodeInvalidTimezone,
	Title:  ErrValidation,
	Detail: "Field timezone value is not present in list of available timezones",
}

// ErrTerminatedMerchantCannotBeModified - Merchant
var ErrTerminatedMerchantCannotBeModified = &Data{
	Code:   ErrCodeTerminatedMerchantCannotBeModified,
	Title:  ErrValidation,
	Detail: "Terminated merchants cannot be modified",
}

// ErrorMap for looking error codes
var ErrorMap = map[int]*Data{}
