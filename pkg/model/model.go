// Package model provides methods for interacting with the database
package model

import (
	"errors"

	"example.com/test/pkg/env"
	"github.com/jmoiron/sqlx"

	"github.com/rs/zerolog"

	"github.com/davecgh/go-spew/spew"
)

// Base -
type Base struct {
	Env    *env.Env
	Logger zerolog.Logger
	DB     *sqlx.Tx
}

// Init -
func (m *Base) Init() error {

	if m.Env == nil {
		return errors.New("Env is nil, cannot initialise")
	}
	if m.DB == nil {
		return errors.New("Database is nil, cannot initialise")
	}

	return nil
}

// Count -
func (m *Base) Count(tableName string) (int64, error) {
	if tableName == "" {
		return -1, errors.New("table name is empty")
	}

	var count int64
	err := m.DB.QueryRow(
		"SELECT COUNT(id) FROM " + tableName + " WHERE deleted_at IS NULL",
	).Scan(&count)
	return count, err
}

// DebugStruct -
func (m *Base) DebugStruct(msg string, rec interface{}) {

	// Don't debug struct if env is not development.
	if m.Env.Get("APP_ENV") != "development" {
		return
	}

	// log
	log := m.Logger

	log.Debug().Msg(msg + " " + spew.Sdump(rec))
}

// NewRecord -
func (m *Base) NewRecord() (rec interface{}) {
	return nil
}

// NewRecordArray -
func (m *Base) NewRecordArray() (rec []interface{}) {
	return nil
}
