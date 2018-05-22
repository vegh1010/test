package main

import (
	"gopkg.in/go-pg/migrations.v5"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		upQuery := `CREATE TYPE ` + GetDatabaseName() +`.e_merchant_status AS ENUM (
		  		'active',
		  		'inactive',
		  		'terminated'
		);`

		_, err := db.Exec(upQuery)

		return err
	}, func(db migrations.DB) error {
		downQuery := `DROP TYPE ` + GetDatabaseName() +`.e_merchant_status;`

		_, err := db.Exec(downQuery)

		return err
	})
}
