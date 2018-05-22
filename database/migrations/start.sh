#!/usr/bin/env bash
echo "Launching Migration in Office Dev mode"

#export PGSQL_PORT=5432
#export PGSQL_PASS=dev
#export PGSQL_HOST=192.168.120.3
#export PGSQL_DB_PREFIX=dev
#export PGSQL_USER=dev
export PORT=5010

# initialize postgres database gopg_migrations table which manages migration changes
go run *.go init

# run migration - it will run the schema based on number sequence
go run *.go

# to rollback migration
# go run *.go down
# this will rollback one migration number
