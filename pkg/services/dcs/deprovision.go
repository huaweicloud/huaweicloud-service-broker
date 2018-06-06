package dcs

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dcs/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/pivotal-cf/brokerapi"
)

// Deprovision implematation
func (b *DCSBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	// Check dcs instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("check dcs instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceDoesNotExist
	if length > 0 {
		return brokerapi.DeprovisionServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// get InstanceDetails in back database
	ids := database.InstanceDetails{}
	err = database.BackDBConnection.
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		First(&ids).Error
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, brokerapi.ErrInstanceDoesNotExist
	}

	// Log InstanceDetails
	b.Logger.Debug(fmt.Sprintf("dcs instance in back database: %v", ids))

	// Init dcs client
	dcsClient, err := b.CloudCredentials.DCSV1Client()
	if err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("create dcs client failed. Error: %s", err)
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("deprovision dcs instance opts: %s", instanceID))

	// Invoke sdk
	result := instances.Delete(dcsClient, ids.TargetID)
	if result.Err != nil {
		return brokerapi.DeprovisionServiceSpec{}, fmt.Errorf("deprovision dcs instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("deprovision dcs instance success: %s", instanceID))

	// Return result
	return brokerapi.DeprovisionServiceSpec{IsAsync: false, OperationData: ""}, nil
}
