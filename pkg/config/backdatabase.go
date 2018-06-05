package config

import (
	"errors"
)

// BackDatabase for broker
type BackDatabase struct {
	DatabaseType     string `json:"database_type"`
	DatabaseHost     string `json:"database_host"`
	DatabasePort     string `json:"database_port"`
	DatabaseName     string `json:"database_name"`
	DatabaseUsername string `json:"database_username"`
	DatabasePassword string `json:"database_password"`
}

// Validate BackDatabase
func (c BackDatabase) Validate() error {

	if c.DatabaseType == "" {
		return errors.New("Must provide a non-empty DatabaseType")
	}

	if c.DatabaseHost == "" {
		return errors.New("Must provide a non-empty DatabaseHost")
	}

	if c.DatabasePort == "" {
		return errors.New("Must provide a non-empty DatabasePort")
	}

	if c.DatabaseName == "" {
		return errors.New("Must provide a non-empty DatabaseName")
	}

	if c.DatabaseUsername == "" {
		return errors.New("Must provide a non-empty DatabaseUsername")
	}

	if c.DatabasePassword == "" {
		return errors.New("Must provide a non-empty DatabasePassword")
	}

	return nil
}
