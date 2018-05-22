#!/usr/bin/env bash
echo "Launching Migration in Office Dev mode"

export APP_DATABASE_HOST=127.0.0.1
export APP_DATABASE_USER=dev
export APP_DATABASE_PASS=dev
export APP_DATABASE_NAME=test
export APP_DATABASE_PORT=5432

# initialize postgres database gopg_migrations table which manages migration changes
go run *.go init

# run migration - it will run the schema based on number sequence
go run *.go

# to rollback migration
# go run *.go down
# this will rollback one migration number
