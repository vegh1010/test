package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Env - contains environment values
type Env struct {
	env map[string]string
}

// NewEnv - Create a new environment
func NewEnv() *Env {

	e := Env{env: map[string]string{}}

	err := e.init()
	if err != nil {
		panic(fmt.Sprintf("NewEnv failed for %v", err))
	}

	return &e
}

func (e *Env) init() error {

	h := os.Getenv("APP_HOME")

	if h == "" {
		panic("APP HOME not set")
	}

	envFile := fmt.Sprintf("%s/%s", os.Getenv("APP_HOME"), ".env")
	err := godotenv.Load(envFile)
	if err != nil {
		// not fatal
		fmt.Printf("Missing local .env file %s, not loading\n", envFile)
	}

	// possible items
	envItems := []string{

		// general
		"APP_ENV",
		"APP_HOME",
		"APP_URL",
		"APP_LOG_LEVEL",
		"APP_TIMED_REQUESTS",
		"APP_PRETTY_LOGS",
		"APP_SERVER_PORT",
		"BUILD_NUMBER",

		// database
		"APP_DATABASE_HOST",
		"APP_DATABASE_USER",
		"APP_DATABASE_PASS",
		"APP_DATABASE_NAME",
		"APP_DATABASE_PORT",
		"APP_DATABASE_MAX_IDLE_CONNS",
		"APP_DATABASE_MAX_OPEN_CONNS",
	}

	// required items
	// - Mostly everything is required to exists even if it
	//   has a bogus value at the moment
	reqItems := []string{

		// general
		"APP_ENV",
		"APP_HOME",
		"APP_URL",
		"APP_LOG_LEVEL",
		"APP_TIMED_REQUESTS",
		"APP_PRETTY_LOGS",
		"APP_SERVER_PORT",

		// database
		"APP_DATABASE_HOST",
		"APP_DATABASE_USER",
		"APP_DATABASE_PASS",
		"APP_DATABASE_NAME",
		"APP_DATABASE_PORT",
	}

	for _, envItem := range envItems {
		osEnvValue := os.Getenv(envItem)

		// dont show passwords or keys
		if strings.Contains(envItem, "PASS") == false && strings.Contains(envItem, "KEY") == false {
			fmt.Printf("ENV Key %s Val %s\n", envItem, osEnvValue)
		}
		e.env[envItem] = os.Getenv(envItem)
	}

	err = e.CheckRequired(reqItems)
	if err != nil {
		panic(fmt.Errorf("Required environment variables missing %v", err))
	}

	return nil
}

// CheckRequired - Tests that the provided list of environment variabled
// have been set
func (e *Env) CheckRequired(reqd []string) error {

	// required check
	for _, reqItem := range reqd {
		if e.env[reqItem] == "" {
			return fmt.Errorf("Could not find required env value for : %s", reqItem)
		}
	}

	return nil
}

// Get an env key
func (e *Env) Get(k string) string {
	v := e.env[k]
	return v
}

// Set an env key
func (e *Env) Set(k string, v string) {
	e.env[k] = v
}
