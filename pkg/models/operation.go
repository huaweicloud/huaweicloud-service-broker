package models

import (
	"encoding/json"
	"fmt"
)

const (
	// OperationProvisioning for asyc provision
	OperationProvisioning string = "provisioning"

	// OperationDeprovisioning for asyc deprovision
	OperationDeprovisioning string = "deprovisioning"

	// OperationUpdating for asyc update
	OperationUpdating string = "updating"
)

const (
	// OperationAsyncDCS for default way
	OperationAsyncDCS bool = true

	// OperationAsyncDMS for default way
	OperationAsyncDMS bool = false

	// OperationAsyncOBS for default way
	OperationAsyncOBS bool = false

	// OperationAsyncRDS for default way
	OperationAsyncRDS bool = true
)

// OperationDatas defines full operation datas
type OperationDatas struct {
	OperationType  string
	ServiceID      string
	PlanID         string
	InstanceID     string
	TargetID       string
	TargetName     string
	TargetStatus   string
	TargetInfo     string
	AdditionalInfo string
}

// ToString for convert
func (ods OperationDatas) ToString() (string, error) {
	// Marshal operation datas
	operationdatas, err := json.Marshal(ods)
	if err != nil {
		return "", fmt.Errorf("marshal operation datas failed. Error: %s", err)
	}
	return string(operationdatas), nil
}
