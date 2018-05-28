package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Config for all
type Config struct {
	BrokerConfig     BrokerConfig     `json:"broker_config"`
	CloudCredentials CloudCredentials `json:"cloud_credentials"`
	Catalog          Catalog          `json:"catalog"`
}

// LoadConfig from file
func LoadConfig(configFile string) (config Config, err error) {
	if configFile == "" {
		return config, errors.New("Must provide a config file")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, err
	}

	if err = config.BrokerConfig.Validate(); err != nil {
		return config, fmt.Errorf("Validating BrokerConfig contents: %s", err)
	}

	if err = config.CloudCredentials.Validate(); err != nil {
		return config, fmt.Errorf("Validating CloudCredentials contents: %s", err)
	}

	if err = config.Catalog.Validate(); err != nil {
		return config, fmt.Errorf("Validating Catalog contents: %s", err)
	}

	return config, nil
}
