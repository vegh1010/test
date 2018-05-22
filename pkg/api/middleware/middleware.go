package middleware

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"example.com/test/pkg/env"
	"example.com/test/pkg/handler"
	"example.com/test/pkg/middleware/tx"
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
