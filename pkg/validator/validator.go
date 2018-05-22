package validator

import (
	"time"

	"example.com/test/pkg/util"
)

// ValidateUUID4 validates that a string is in a UUID format.
func ValidateUUID4(uuid string) bool {
	return util.ValidateUUID(uuid)
}

// ValidateDotTwoPrecision -
func ValidateDotTwoPrecision(total float64) bool {
	if float64(int(total*100))/100 != total {
		return false
	}
	return true
}

// ValidateTimestampFormat checks that a timestamp is in RFC3339 format.
func ValidateTimestampFormat(timestamp string) bool {
	if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
		return false
	}
	return true
}
