package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// InstanceDetailsTableName defines
var InstanceDetailsTableName = "instance_details"

// InstanceDetailsTableSQL matches with Upgrades Object
var InstanceDetailsTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	service_id varchar(255) DEFAULT NULL,
	plan_id varchar(255) DEFAULT NULL,
	instance_id varchar(255) NOT NULL,
	target_id varchar(255) DEFAULT NULL,
	target_name varchar(255) DEFAULT NULL,
	target_status varchar(255) DEFAULT NULL,
	target_info text,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, InstanceDetailsTableName)

// InstanceDetails defines for back database
type InstanceDetails struct {
	gorm.Model
	ServiceID      string
	PlanID         string
	InstanceID     string
	TargetID       string
	TargetName     string
	TargetStatus   string
	TargetInfo     string `sql:"type:text"`
	AdditionalInfo string `sql:"type:text"`
}

// GetTargetInfo for InstanceDetails
func (ids InstanceDetails) GetTargetInfo(targetinfo interface{}) error {
	if ids.TargetInfo != "" {
		err := json.Unmarshal([]byte(ids.TargetInfo), &targetinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAdditionalInfo for InstanceDetails
func (ids InstanceDetails) GetAdditionalInfo(additionalinfo interface{}) error {
	if ids.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(ids.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}
