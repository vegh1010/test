package tx

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog"

	"example.com/test/pkg/env"
	"example.com/test/pkg/txcontext"
)

// tx -
type tx struct {
	Env    *env.Env
	Logger zerolog.Logger
	DB     *sqlx.DB
}

// NewTx -
func NewTx(e *env.Env, l zerolog.Logger, db *sqlx.DB, h http.Handler) http.Handler {

	a := &tx{
		Env:    e,
		Logger: l,
		DB:     db,
	}

	mw := a.Middleware(h)

	return mw
}

// Middleware -
func (t tx) Middleware(h http.Handler) http.Handler {

	log := t.Logger

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := t.DB.Beginx()

		if err != nil {
			log.Error().Msgf("Could not Beginx in tx for %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal error"))
			return
		}
		r = txcontext.SetContext(r, tx)

		h.ServeHTTP(w, r)

	})
}
