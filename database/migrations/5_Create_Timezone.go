package main

import (
	"gopkg.in/go-pg/migrations.v5"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		upQuery := `CREATE TABLE ` + GetDatabaseName() +`.timezone (
		  			id              TEXT                NOT NULL,
		  			status          e_timezone_status   NOT NULL DEFAULT 'active',
					created_at      TIMESTAMP           NOT NULL DEFAULT now(),
					updated_at      TIMESTAMP           NULL,
					deleted_at      TIMESTAMP           NULL,
		  			CONSTRAINT 		timezone_pk PRIMARY KEY (id)
		);`

		_, err := db.Exec(upQuery)

		return err
	}, func(db migrations.DB) error {
		downQuery := `DROP TABLE ` + GetDatabaseName() +`.timezone;`

		_, err := db.Exec(downQuery)

		return err
	})
}
