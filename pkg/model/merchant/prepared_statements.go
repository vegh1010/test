package merchant

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var getByIDStmt *sqlx.Stmt
var getByIDSQL = `
SELECT *
FROM merchant
WHERE id = $1
AND deleted_at IS NULL
`

var createRecordStmt *sqlx.NamedStmt
var createRecordSQL = `
INSERT INTO merchant (
	id,
	name,
	short_name,
	dba_name,
	country_id,
	timezone_id,
	status,
	created_at
) VALUES (
	:id,
	:name,
	:short_name,
	:dba_name,
	:country_id,
	:timezone_id,
	:status,
	:created_at
)
RETURNING
	id,
	name,
	short_name,
	dba_name,
	country_id,
	timezone_id,
	status,
	created_at,
	updated_at,
	deleted_at
`

var updateRecordStmt *sqlx.NamedStmt
var updateRecordSQL = `
UPDATE merchant SET
	name 	       = :name,
	short_name     = :short_name,
	dba_name       = :dba_name,
	country_id     = :country_id,
	timezone_id    = :timezone_id,
	status         = :status,
	updated_at     = :updated_at
WHERE id = :id
AND deleted_at IS NULL
RETURNING
	id,
	name,
	short_name,
	dba_name,
	country_id,
	timezone_id,
	status,
	created_at,
	updated_at,
	deleted_at
`

var deleteRecordStmt *sqlx.NamedStmt
var deleteRecordSQL = `
UPDATE merchant SET
	deleted_at = :deleted_at
WHERE id = :id
AND deleted_at IS NULL
RETURNING
	id,
	name,
	short_name,
	dba_name,
	country_id,
	timezone_id,
	status,
	created_at,
	updated_at,
	deleted_at
`

var removeRecordStmt *sqlx.Stmt
var removeRecordSQL = `
DELETE FROM merchant
WHERE id = $1
`

var validateRecordWithIDStmt *sqlx.NamedStmt
var validateRecordWithIDSQL = `
SELECT
(
	SELECT 1
	FROM   country
	WHERE  id = :country_id
	AND    status = 'active'
	AND    deleted_at IS NULL
) country_id, (
	SELECT 1
	FROM   timezone
	WHERE  id = :timezone_id
	AND    status = 'active'
	AND    deleted_at IS NULL
) timezone_id
`

var validateRecordWithoutIDStmt *sqlx.NamedStmt
var validateRecordWithoutIDSQL = `
SELECT
(
	SELECT 1
	FROM   country
	WHERE  id = :country_id
	AND    status = 'active'
	AND    deleted_at IS NULL
) country_id, (
	SELECT 1
	FROM   timezone
	WHERE  id = :timezone_id
	AND    status = 'active'
	AND    deleted_at IS NULL
) timezone_id
`

// PrepareStatements prepares sql statements
func PrepareStatements(db *sqlx.DB) {
	var err error

	getByIDStmt, err = db.Preparex(getByIDSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare getByIDSQL %v", err)
	}

	createRecordStmt, err = db.PrepareNamed(createRecordSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare createRecordSQL %v", err)
	}

	updateRecordStmt, err = db.PrepareNamed(updateRecordSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare updateRecordSQL %v", err)
	}

	deleteRecordStmt, err = db.PrepareNamed(deleteRecordSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare deleteRecordSQL %v", err)
	}

	removeRecordStmt, err = db.Preparex(removeRecordSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare removeRecordSQL %v", err)
	}

	validateRecordWithIDStmt, err = db.PrepareNamed(validateRecordWithIDSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare validateRecordWithIDSQL %v", err)
	}

	validateRecordWithoutIDStmt, err = db.PrepareNamed(validateRecordWithoutIDSQL)
	if err != nil {
		log.Fatal().Msgf("Failed to prepare validateRecordWithoutIDSQL %v", err)
	}

}
