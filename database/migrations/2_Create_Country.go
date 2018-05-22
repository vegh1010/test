package main

import (
	"gopkg.in/go-pg/migrations.v5"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		upQuery := `CREATE TABLE ` + GetDatabaseName() +`.country (
				id              VARCHAR(2)        NOT NULL,
				name            TEXT              NOT NULL,
				alpha2_code     VARCHAR(2)        NOT NULL,
				alpha3_code     VARCHAR(3)        NOT NULL,
				numeric_code    VARCHAR(3)        NOT NULL,
				status          e_country_status  NOT NULL DEFAULT 'active',
				created_at      TIMESTAMP         NOT NULL DEFAULT now(),
				updated_at      TIMESTAMP         NULL,
				deleted_at      TIMESTAMP         NULL,
  				CONSTRAINT country_pk PRIMARY KEY (id)
		);`

		_, err := db.Exec(upQuery)

		return err
	}, func(db migrations.DB) error {
		downQuery := `DROP TABLE ` + GetDatabaseName() +`.country;`

		_, err := db.Exec(downQuery)

		return err
	})
}
