package util

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vegh1010/test/pkg/wrappederror"
)

// GetUUID returns a unique identifier
func GetUUID() string {
	uuidByte, _ := uuid.NewRandom()
	uuidString := uuidByte.String()
	return uuidString
}

// ValidateUUID validates a unique identifier
func ValidateUUID(s string) bool {
	_, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return true
}

// GetTime returns the current UTC time
func GetTime() string {

	t := time.Now()

	// format UTC
	tf := t.UTC().Format(time.RFC3339)

	return tf
}

// GetFutureTime returns the current UTC time
func GetFutureTime(d time.Duration) string {

	t := time.Now().UTC()

	// add duration
	if d != 0 {
		t = t.Add(d)
	}

	// format UTC
	tf := t.Format(time.RFC3339)

	return tf
}

// ToNullString Convert an ordinary string into sql.NullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// ToNullInt64 Convert an ordinary string into sql.NullInt64
func ToNullInt64(i int64) sql.NullInt64 {

	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// ToNullBool Convert an ordinary bool into sql.NullBool
func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

// PtrToNullBool Convert an ordinary bool into sql.NullBool
func PtrToNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  *b,
		Valid: true,
	}
}

// StringInSlice -
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GenerateToken -
func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// MaskLastFourCharactersClear masks a string except for the last 4 characters.
//
// NOTE: This seems to be what ServiceCloud uses for bank account information.
func MaskLastFourCharactersClear(str string) string {
	// If the string is less than 5 characters long, just return it.
	if len(str) < 5 {
		return str
	}
	strPrefixLen := len(str) - 4
	return strings.Repeat("*", strPrefixLen) + str[strPrefixLen:]
}

// RollbackTxWithError is a helper function that either returns a wrappederror.Err
// or also wraps that error in another error if the tx rollback fails.
//
// Even if e is nil, RollbackTxWithError will never return a nil value.
func RollbackTxWithError(e error, msg string, tx *sqlx.Tx) error {
	werrors := wrappederror.New()

	// Add the first error.
	werrors.Add(e, msg)

	err := tx.Rollback()
	if err != nil {
		werrors.Add(err, "Error rolling back tx")
	}

	return werrors
}
