package main

import (
	"fmt"
	"net/http"
	"runtime"
	"github.com/vegh1010/test/pkg/db"
	"github.com/vegh1010/test/pkg/model/modelinit"
	"github.com/vegh1010/test/pkg/api/router"
	"github.com/vegh1010/test/pkg/env"
	"github.com/vegh1010/test/pkg/logger"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// environment
	e := env.NewEnv()

	// logger
	l := logger.NewLogger(e)

	// database
	db := db.NewDB(l, e)

	// prepare model statements.
	l.Info().Msg("Preparing model statements")
	modelinit.PrepareStatements(db)

	// router
	r, err := router.NewRouter(e, l, db)
	if err != nil {
		panic(fmt.Sprintf("Router error: %v", err))
	}

	// server
	sp := e.Get("APP_SERVER_PORT")
	l.Info().Msgf("Listing on http://0.0.0.0:%s", sp)

	l.Error().Msgf(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf(":%s", sp), r)))
}
