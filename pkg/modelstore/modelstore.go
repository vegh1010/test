// Package modelstore provides methods for interacting with the DB
package modelstore

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"example.com/test/pkg/env"

	// models

	"example.com/test/pkg/model/merchant"
)

// ModelStore - contains a map of model structs
type ModelStore struct {
	models map[string]interface{}
	Env    *env.Env
	Logger zerolog.Logger
	DB     *sqlx.Tx
}

// NewModelStore -
func NewModelStore(e *env.Env, l zerolog.Logger, d *sqlx.Tx) (*ModelStore, error) {
	m := ModelStore{
		Env:    e,
		Logger: l,
		DB:     d,
	}
	err := m.init()
	return &m, err
}

// modelstore init initialises all models creating a named map
func (m *ModelStore) init() error {

	// log
	log := m.Logger

	log.Debug().Msg("Initializing models")

	m.models = make(map[string]interface{})

	var err error

	// models
	m.models["merchant"], err = merchant.NewModel(m.Env, m.Logger, m.DB)

	log.Debug().Msg("Done Initializing models")

	return err
}

// GetMerchantModel -
func (m *ModelStore) GetMerchantModel() (*merchant.Model, error) {

	model := m.models["merchant"]
	if model == nil {
		return nil, errors.New("Merchant model does not exist")
	}

	return model.(*merchant.Model), nil
}
