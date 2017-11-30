package rdsbroker

import (
	"errors"
	"fmt"
)

type Config struct {
	IdentityEndpoint             string  `json:"identity_endpoint"`
	Ca                           string  `json:"ca"`
	Username                     string  `json:"username"`
	Password                     string  `json:"password"`
	DomainName                   string  `json:"domain_name"`
	ProjectName                  string  `json:"project_name"`
	ProjectID                    string  `json:"project_id"`
	Region                       string  `json:"region"`
	//
	DBPrefix                     string  `json:"db_prefix"`
	AllowUserProvisionParameters bool    `json:"allow_user_provision_parameters"`
	AllowUserUpdateParameters    bool    `json:"allow_user_update_parameters"`
	AllowUserBindParameters      bool    `json:"allow_user_bind_parameters"`
	Catalog                      Catalog `json:"catalog"`
}

func (c Config) Validate() error {
	if c.Region == "" {
		return errors.New("Must provide a non-empty Region")
	}

	if c.DBPrefix == "" {
		return errors.New("Must provide a non-empty DBPrefix")
	}

	if err := c.Catalog.Validate(); err != nil {
		return fmt.Errorf("Validating Catalog configuration: %s", err)
	}

	return nil
}
