package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	postgres "github.com/lib/pq"
	"github.com/rs/zerolog"
	"example.com/test/pkg/env"
)

// NewDB gets a new db and returns an error if need be
func NewDB(l zerolog.Logger, e *env.Env) *sqlx.DB {
	var conn *sqlx.DB
	var err error

	connectString := connectString(e)

	// timed requests
	if e.Get("APP_TIMED_REQUESTS") != "" && e.Get("APP_TIMED_REQUESTS") != "0" {
		sql.Register("wrapped-postgres", WrapDriver(&postgres.Driver{}, l))
		db, err := sql.Open("wrapped-postgres", connectString)
		if err != nil {
			panic(fmt.Sprintf("MustGetNewDB error: %v", err))
		}
		conn = sqlx.NewDb(db, "postgres")
	} else {
		conn, err = sqlx.Connect("postgres", connectString)
		if err != nil {
			panic(fmt.Sprintf("MustGetNewDB error: %v", err))
		}
	}

	poolConfig(e, conn)

	return conn
}

func connectString(e *env.Env) string {

	host := e.Get("APP_DATABASE_HOST")
	user := e.Get("APP_DATABASE_USER")
	pass := e.Get("APP_DATABASE_PASS")
	port := e.Get("APP_DATABASE_PORT")
	dbname := e.Get("APP_DATABASE_NAME")

	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, pass, dbname, host, port)

}

func poolConfig(e *env.Env, d *sqlx.DB) error {

	maxIdleCons := 50
	maxOpenCons := 100
	var err error

	maxIdleConsStr := e.Get("APP_DATABASE_MAX_IDLE_CONNS")

	if maxIdleConsStr != "" {
		maxIdleCons, err = strconv.Atoi(maxIdleConsStr)
		if err != nil {
			return err
		}
	}

	maxOpenConsStr := e.Get("APP_DATABASE_MAX_OPEN_CONNS")

	if maxOpenConsStr != "" {
		maxOpenCons, err = strconv.Atoi(maxOpenConsStr)
		if err != nil {
			return err
		}
	}

	d.SetMaxIdleConns(maxIdleCons)
	d.SetMaxOpenConns(maxOpenCons)

	return nil
}
