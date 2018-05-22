package main

import (
	"fmt"
	"gopkg.in/pg.v5"
	"os"
)

//initialize database connection
func DBConnect() (*pg.DB) {
	fmt.Println("Database connection established")

	envs := GetEnv()
	for key, value := range envs {
		fmt.Println(key, ":", value)
	}

	dbusername := envs["APP_DATABASE_USER"]
	dbpassword := envs["APP_DATABASE_PASS"]
	dbhost := envs["APP_DATABASE_HOST"]
	dbport := envs["APP_DATABASE_PORT"]
	addr := fmt.Sprint(dbhost, ":", dbport)

	db := pg.Connect(&pg.Options{
		User:     dbusername,
		Password: dbpassword,
		Addr:     addr,
		Database: "postgres",
	})
	return db
}

//get database name based on env APP_DATABASE_NAME
func GetDatabaseName() (string) {
	configs := GetEnv()

	return configs["APP_DATABASE_NAME"]
}

//get all env needed for postgres
func GetEnv() (map[string]string) {
	envs := []string{
		"APP_DATABASE_HOST",
		"APP_DATABASE_USER",
		"APP_DATABASE_PASS",
		"APP_DATABASE_NAME",
		"APP_DATABASE_PORT",
	}
	var configs = make(map[string]string)

	for _, envField := range envs {
		field := os.Getenv(envField)
		if field == "" {
			panic("Environment " + envField + " Not Found")
		}
		configs[field] = field
	}

	return configs
}
