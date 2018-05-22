package main

import (
	"flag"
	"fmt"
	"gopkg.in/go-pg/migrations.v5"
)

const verbose = true

func main() {
	migrationDB := DBConnect()
	fmt.Println(flag.Args())
	database := GetDatabaseName()
	migrationName := database + `.gopg_migrations`
	migrations.SetTableName(migrationName)

	oldVersion, newVersion, err := migrations.Run(migrationDB, flag.Args()...)
	if err != nil {
		panic(err)
	}
	if verbose {
		if newVersion != oldVersion {
			fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		} else {
			fmt.Printf("version is %d\n", oldVersion)
		}
	}
}
