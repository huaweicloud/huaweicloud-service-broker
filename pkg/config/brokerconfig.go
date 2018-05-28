package config

import "errors"

// BrokerConfig for broker
type BrokerConfig struct {
	LogLevel string `json:"log_level"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate BrokerConfig
func (c BrokerConfig) Validate() error {
	if c.LogLevel == "" {
		return errors.New("Must provide a non-empty LogLevel")
	}

	if c.Username == "" {
		return errors.New("Must provide a non-empty Username")
	}

	if c.Password == "" {
		return errors.New("Must provide a non-empty Password")
	}

	return nil
}
