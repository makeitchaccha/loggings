package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseDriver string

const (
	DatabaseDriverPostgres DatabaseDriver = "postgres"
	DatabaseDriverSqlite   DatabaseDriver = "sqlite"
)

func (d DatabaseDriver) Valid() bool {
	switch d {
	case DatabaseDriverPostgres, DatabaseDriverSqlite:
		return true
	default:
		return false
	}
}

type Config struct {
	Database struct {
		Driver DatabaseDriver `yaml:"driver"`
		Dsn    string         `yaml:"dsn"`
	} `yaml:"database"`
}

func (c Config) DatabaseDialector() gorm.Dialector {
	switch c.Database.Driver {
	case DatabaseDriverPostgres:
		return postgres.Open(c.Database.Dsn)
	case DatabaseDriverSqlite:
		return sqlite.Open(c.Database.Dsn)
	default:
		panic("unknown database driver")
	}
}

func New(filename string) Config {
	config, err := loadConfig(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config, continue to use blank config, reason: %v\n", err)
		config = Config{}
	}

	// override config from env
	overrideStringConfig("DATABASE_DRIVER", &config.Database.Driver)
	overrideStringConfig("DATABASE_DSN", &config.Database.Dsn)

	err = validate(config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to validate config, continue to use blank config, reason: %v\n", err)
		os.Exit(1)
	}

	return config
}

func loadConfig(filename string) (Config, error) {
	// load config from file
	file, err := os.Open(filename)

	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}

	var config Config
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}

func overrideConfig[T any](env string, target *T, f func(string) T) {
	if value, ok := os.LookupEnv(env); ok {
		*target = f(value)
	}
}

func overrideStringConfig[T ~string](env string, target *T) {
	overrideConfig(env, target, func(value string) T {
		return T(value)
	})
}

func validate(config Config) error {
	if !config.Database.Driver.Valid() {
		return fmt.Errorf("unknown database driver: %s", config.Database.Driver)
	}

	if config.Database.Dsn == "" {
		return fmt.Errorf("database dsn is required")
	}

	return nil
}
