package rds

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/pivotal-cf/brokerapi"
)

// Unbind implematation
func (b *RDSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {

	// Log opts
	b.Logger.Debug(fmt.Sprintf("unbind rds instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// Check rds bind length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.BindDetails{}).
		Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return fmt.Errorf("check rds bind length in back database failed. Error: %s", err)
	}
	// ErrBindingDoesNotExist
	if length == 0 {
		return brokerapi.ErrBindingDoesNotExist
	}

	// Check rds instance length in back database
	err = database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return fmt.Errorf("check rds instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceDoesNotExist
	if length == 0 {
		return brokerapi.ErrInstanceDoesNotExist
	}

	// Get BindDetails in back database
	bds := database.BindDetails{}
	err = database.BackDBConnection.
		Where("bind_id = ? and instance_id = ? and service_id = ? and plan_id = ?", bindingID, instanceID, details.ServiceID, details.PlanID).
		First(&bds).Error
	if err != nil {
		return brokerapi.ErrBindingDoesNotExist
	}

	// Delete BindDetails in back database
	err = database.BackDBConnection.Delete(&bds).Error
	if err != nil {
		return fmt.Errorf("delete rds bind in back database failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("unbind rds instance success: instanceID: %s bindingID: %s", instanceID, bindingID))

	return nil
}
