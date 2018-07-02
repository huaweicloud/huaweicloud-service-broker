package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// OperationDetailsTableName defines
var OperationDetailsTableName = "operation_details"

// OperationDetailsTableSQL matches with OperationDetails Object
var OperationDetailsTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	operation_type varchar(255) DEFAULT NULL,
	service_id varchar(255) DEFAULT NULL,
	plan_id varchar(255) DEFAULT NULL,
	instance_id varchar(255) DEFAULT NULL,
	target_id varchar(255) DEFAULT NULL,
	target_name varchar(255) DEFAULT NULL,
	target_status varchar(255) DEFAULT NULL,
	target_info text,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, OperationDetailsTableName)

// OperationDetails defines for back database
type OperationDetails struct {
	gorm.Model
	OperationType  string
	ServiceID      string
	PlanID         string
	InstanceID     string
	TargetID       string
	TargetName     string
	TargetStatus   string
	TargetInfo     string `sql:"type:text"`
	AdditionalInfo string `sql:"type:text"`
}

// GetTargetInfo for OperationDetails
func (ods OperationDetails) GetTargetInfo(targetinfo interface{}) error {
	if ods.TargetInfo != "" {
		err := json.Unmarshal([]byte(ods.TargetInfo), &targetinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAdditionalInfo for OperationDetails
func (ods OperationDetails) GetAdditionalInfo(additionalinfo interface{}) error {
	if ods.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(ods.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (ods OperationDetails) ToString() (string, error) {
	// Marshal operation datas
	operationdatas, err := json.Marshal(ods)
	if err != nil {
		return "", fmt.Errorf("marshal operation details failed. Error: %s", err)
	}
	return string(operationdatas), nil
}
