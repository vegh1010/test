package main

import (
	"gopkg.in/go-pg/migrations.v5"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		upQuery := `CREATE TABLE ` + GetDatabaseName() +`.merchant (
					id            	UUID              NOT NULL DEFAULT gen_random_uuid(),
		  			name          	TEXT              NOT NULL,
		  			short_name    	TEXT              NOT NULL,
		  			dba_name      	TEXT              NOT NULL,
		  			country_id    	VARCHAR(2)        NOT NULL,
		  			timezone_id   	TEXT              NOT NULL,
		  			status        	` + GetDatabaseName() +`.e_merchant_status NOT NULL DEFAULT 'active',
					created_at    	TIMESTAMP         NOT NULL DEFAULT now(),
					updated_at    	TIMESTAMP         NULL,
					deleted_at    	TIMESTAMP         NULL,
					CONSTRAINT 		merchant_pk PRIMARY KEY (id),
		  			CONSTRAINT 		merchant_country_fk FOREIGN KEY (country_id) REFERENCES country (id),
		  			CONSTRAINT 		merchant_timezone_fk FOREIGN KEY (timezone_id) REFERENCES timezone (id)
		);`

		_, err := db.Exec(upQuery)

		return err
	}, func(db migrations.DB) error {
		downQuery := `DROP TABLE ` + GetDatabaseName() +`.merchant;`

		_, err := db.Exec(downQuery)

		return err
	})
}
