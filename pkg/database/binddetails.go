package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// BindDetailsTableName defines
var BindDetailsTableName = "bind_details"

// BindDetailsTableSQL matches with Upgrades Object
var BindDetailsTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	service_id varchar(255) DEFAULT NULL,
	plan_id varchar(255) DEFAULT NULL,
	instance_id varchar(255) DEFAULT NULL,
	bind_id varchar(255) DEFAULT NULL,
	bind_info text,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, BindDetailsTableName)

// BindDetails defines for back database
type BindDetails struct {
	gorm.Model
	ServiceID      string
	PlanID         string
	InstanceID     string
	BindID         string
	BindInfo       string `sql:"type:text"`
	AdditionalInfo string `sql:"type:text"`
}

// GetBindInfo for BindDetails
func (ids BindDetails) GetBindInfo(bindinfo interface{}) error {
	if ids.BindInfo != "" {
		err := json.Unmarshal([]byte(ids.BindInfo), &bindinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAdditionalInfo for BindDetails
func (ids BindDetails) GetAdditionalInfo(additionalinfo interface{}) error {
	if ids.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(ids.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}
