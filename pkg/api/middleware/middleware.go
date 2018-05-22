package middleware

import (
	"net/http"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/vegh1010/test/pkg/handler"
	"github.com/vegh1010/test/pkg/middleware/tx"
	"github.com/vegh1010/test/pkg/env"
)

// Middleware -
type Middleware struct {
	e  *env.Env
	l  zerolog.Logger
	db *sqlx.DB
}

// NewMiddleware returns a handler with all appropriate middleware applied to a specified handler.
func NewMiddleware(e *env.Env, l zerolog.Logger, db *sqlx.DB) *Middleware {
	return &Middleware{e: e, l: l, db: db}
}

// Apply - Applies selected middleware to handler chain
func (mw *Middleware) Apply(h handler.Handler, hf http.HandlerFunc, path string) http.Handler {
	var nh http.Handler = hf

	// tx
	nh = tx.NewTx(mw.e, mw.l, mw.db, nh)

	return nh
}
