package txcontext

import (
	"context"
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type keyType string

// Key -
const Key keyType = "TxContext"

// ErrTxContextEmpty -
var ErrTxContextEmpty = errors.New("Could not find TxContext : context empty")

// GetContext returns the current tx context
func GetContext(r *http.Request) (*sqlx.Tx, error) {
	ctx := r.Context().Value(Key)
	if ctx == nil {
		return nil, ErrTxContextEmpty
	}
	tx := ctx.(*sqlx.Tx)
	return tx, nil
}

// SetContext for current transaction
func SetContext(r *http.Request, tx *sqlx.Tx) *http.Request {

	ctx := context.WithValue(r.Context(), Key, tx)

	r = r.WithContext(ctx)

	return r
}
