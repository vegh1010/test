package merchant

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vegh1010/test/pkg/model"
	"github.com/vegh1010/test/pkg/env"
	"github.com/vegh1010/test/pkg/util"
)

// Record -
type Record struct {
	ID         string         `db:"id"`
	Name       string         `db:"name"`
	ShortName  string         `db:"short_name"`
	DBAName    string         `db:"dba_name"`
	CountryID  string         `db:"country_id"`
	TimezoneID string         `db:"timezone_id"`
	Status     string         `db:"status"`
	CreatedAt  string         `db:"created_at"`
	UpdatedAt  sql.NullString `db:"updated_at"`
	DeletedAt  sql.NullString `db:"deleted_at"`
}

// StatusCommentRecord -
type StatusCommentRecord struct {
	ID         string         `db:"id"`
	MerchantID string         `db:"merchant_id"`
	NewStatus  string         `db:"new_status"`
	OldStatus  string         `db:"old_status"`
	Comment    string         `db:"comment"`
	CreatedAt  string         `db:"created_at"`
	UpdatedAt  sql.NullString `db:"updated_at"`
	DeletedAt  sql.NullString `db:"deleted_at"`
}

// External ID status values
const (
	StatusActive     = "active"
	StatusInactive   = "inactive"
	StatusTerminated = "terminated"
)

// Model -
type Model struct {
	model.Base
}

// NewModel -
func NewModel(e *env.Env, l zerolog.Logger, d *sqlx.Tx) (*Model, error) {
	m := Model{
		model.Base{
			DB:     d,
			Env:    e,
			Logger: l,
		},
	}
	err := m.Init()
	return &m, err
}

// NewRecord -
func (m *Model) NewRecord() Record {
	return Record{}
}

// NewStatusCommentRecord -
func (m *Model) NewStatusCommentRecord() StatusCommentRecord {
	return StatusCommentRecord{}
}

// GetByID -
func (m *Model) GetByID(id string) (*Record, error) {

	// record
	rec := m.NewRecord()
	rec.ID = id

	// log
	log := m.Logger

	log.Debug().Msgf("Fetching merchant record by ID %s", id)

	// db
	db := m.DB

	stmt := db.Stmtx(getByIDStmt)

	err := stmt.QueryRowx(rec.ID).StructScan(&rec)
	if err != nil {
		log.Error().Msgf("Error executing update ", err)
		return nil, err
	}

	return &rec, nil
}

// GetByParam -
func (m *Model) GetByParam(params map[string]interface{}) ([]*Record, error) {

	// records
	var recs []*Record

	// log
	log := m.Logger

	// db
	db := m.DB

	// sqlStmt
	sqlStmt := `
SELECT *
FROM merchant
WHERE deleted_at IS NULL
`

	// params
	for k := range params {
		sqlStmt = sqlStmt + fmt.Sprintf("AND %s = :%s\n", k, k)
	}

	rows, err := db.NamedQuery(sqlStmt, params)
	if err != nil {
		log.Error().Msgf("Error querying row %s", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e Record
		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		recs = append(recs, &e)
	}

	m.DebugStruct("Fetched", recs)

	return recs, rows.Err()
}

// GetOneByParam -
func (m *Model) GetOneByParam(params map[string]interface{}) (*Record, error) {

	// records
	recs, err := m.GetByParam(params)
	if err != nil {

		log.Warn().Msgf("Error querying %s", err)
		return nil, err
	}

	if len(recs) != 1 {
		return nil, sql.ErrNoRows
	}

	m.DebugStruct("Fetched", recs[0])

	return recs[0], nil
}

// Create -
func (m *Model) Create(rec *Record) error {

	// log
	log := m.Logger

	// db
	db := m.DB

	stmt := db.NamedStmt(createRecordStmt)

	// id
	rec.ID = util.GetUUID()

	// status - initially is always inactive
	rec.Status = "inactive"

	// created at
	rec.CreatedAt = util.GetTime()

	m.DebugStruct("Create ", rec)

	err := stmt.QueryRowx(rec).StructScan(rec)
	if err != nil {
		log.Error().Msgf("Error executing insert %v", err)
		return err
	}

	return nil
}

// Update -
func (m *Model) Update(rec *Record) error {

	// log
	log := m.Logger

	// db
	db := m.DB

	stmt := db.NamedStmt(updateRecordStmt)

	oldUpdatedAt := rec.UpdatedAt

	rec.UpdatedAt.String = util.GetTime()
	rec.UpdatedAt.Valid = true

	err := stmt.QueryRowx(rec).StructScan(rec)
	if err != nil {
		rec.UpdatedAt = oldUpdatedAt
		log.Error().Msgf("Error executing update %v", err)
		return err
	}

	return nil
}

// Delete -
func (m *Model) Delete(id string) error {

	// log
	log := m.Logger

	log.Debug().Msgf("Delete ID %s", id)

	// db
	db := m.DB

	rec := m.NewRecord()
	rec.ID = id

	stmt := db.NamedStmt(deleteRecordStmt)

	// deleted at
	rec.DeletedAt.String = util.GetTime()
	rec.DeletedAt.Valid = true

	err := stmt.QueryRowx(rec).StructScan(&rec)
	if err != nil {
		log.Error().Msgf("Error executing delete %s", err)
		return err
	}

	return nil
}

// Remove -
func (m *Model) Remove(id string) error {

	// log
	log := m.Logger

	log.Debug().Msgf("Remove merchant record ID %s", id)

	// db
	db := m.DB

	rec := m.NewRecord()
	rec.ID = id

	// remove merchant record
	stmt := db.Stmtx(removeRecordStmt)

	res, err := stmt.Exec(rec.ID)
	if err != nil {
		log.Error().Msgf("Error executing delete %s", err)
		return err
	}

	// rows affected
	raf, err := res.RowsAffected()
	if raf != 1 {
		return fmt.Errorf("Expecting to remove exactly one row removed, got %d  ", raf)
	}

	log.Debug().Msgf("Result RowsAffected %d", raf)

	return err
}

// ValidateResult is used for validating against merchant config and existing transactions
type ValidateResult struct {
	CountryID  sql.NullBool `db:"country_id"`
	TimezoneID sql.NullBool `db:"timezone_id"`
}

// ValidateRecord - validates properties of a record are valid for creating or updating
func (m *Model) ValidateRecord(rec *Record) (*ValidateResult, error) {

	// log
	log := m.Logger

	// db
	db := m.DB

	log.Debug().Msgf("Validating merchant record %v", rec)

	var stmt *sqlx.NamedStmt

	if rec.ID != "" {
		stmt = db.NamedStmt(validateRecordWithIDStmt)
	} else {
		stmt = db.NamedStmt(validateRecordWithoutIDStmt)
	}

	vrec := ValidateResult{}

	err := stmt.QueryRowx(rec).StructScan(&vrec)
	if err != nil {
		log.Error().Msgf("Error executing validation query %v", err)
		return nil, err
	}

	return &vrec, nil
}
