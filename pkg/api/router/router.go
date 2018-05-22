// Package router maps HTTP methods to handler functions
package router

import (
	"net/http"

	// pprof

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"example.com/test/pkg/api/handler/merchant"
	"example.com/test/pkg/api/middleware"
	"example.com/test/pkg/env"
)

// Router -
type Router struct {
	Env      *env.Env
	Logger   zerolog.Logger
	handler  http.Handler
	basePath string
	db       *sqlx.DB
}

// NewRouter -
func NewRouter(e *env.Env, l zerolog.Logger, db *sqlx.DB) (http.Handler, error) {

	r := Router{
		Env:    e,
		Logger: l,
		db:     db,
	}
	err := r.init()

	return r.handler, err
}

func (rt *Router) init() error {

	log := rt.Logger

	log.Debug().Msgf("Initializing routes")

	m := mux.NewRouter()

	mw := middleware.NewMiddleware(rt.Env, rt.Logger, rt.db)

	// Merchants
	mh := merchant.NewHandler(rt.Env, rt.Logger)
	m.Handle(mh.GetPath(), mw.Apply(mh, mh.Post, "merchants")).Methods(http.MethodPost)
	m.Handle(mh.GetPath(), mw.Apply(mh, mh.GetCollection, "merchants")).Methods(http.MethodGet)
	m.Handle(mh.GetPath()+"/{id}", mw.Apply(mh, mh.Get, "merchants")).Methods(http.MethodGet)
	m.Handle(mh.GetPath()+"/{id}", mw.Apply(mh, mh.Delete, "merchants")).Methods(http.MethodDelete)
	m.Handle(mh.GetPath()+"/{id}", mw.Apply(mh, mh.Put, "merchants")).Methods(http.MethodPut)

	rt.handler = m

	// Set the not found handler.
	// TODO: When there's some time for refactoring, replace handler interface's
	//		 'base methods' with this handler.
	m.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	return nil
}
