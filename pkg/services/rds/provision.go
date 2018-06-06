package rds

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *RDSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Check rds instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("check rds instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceAlreadyExists
	if length > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Init provisionOpts
	provisionOpts := instances.CreateOps{}
	/*Name:             instanceID,
	Datastore:        map[string]string{"type": rds_prop.DatastoreType, "version": rds_prop.DatastoreVersion},
	FlavorRef:        rds_prop.FlavorId,
	Volume:           map[string]interface{}{"type": rds_prop.VolumeType, "size": rds_prop.VolumeSize},
	Region:           rds_prop.Region,
	AvailabilityZone: rds_prop.AvailabilityZone,
	Vpc:              rds_prop.VpcId,
	Nics:             map[string]string{"subnetId": rds_prop.SubnetId},
	SecurityGroup:    map[string]string{"id": rds_prop.SecurityGroupId},
	DbPort:           rds_prop.Dbport,
	BackupStrategy:   map[string]interface{}{"startTime": rds_prop.BackupStrategyStarttime, "keepDays": rds_prop.BackupStrategyKeepdays},
	DbRtPd:           rds_prop.Dbpassword,*/
	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &provisionOpts)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error unmarshalling rawParameters: %s", err)
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision rds instance opts: %v", provisionOpts))

	// Invoke sdk
	rdsInstance, err := instances.Create(rdsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision rds instance failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision rds instance result: %v", rdsInstance))

	// Invoke sdk get
	freshInstance, err := instances.Get(rdsClient, rdsInstance.ID).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get rds instance failed. Error: %s", err)
	}

	// Marshal instance
	targetinfo, err := json.Marshal(freshInstance)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal rds instance failed. Error: %s", err)
	}

	// create InstanceDetails in back database
	idsOpts := database.InstanceDetails{
		ServiceID:      details.ServiceID,
		PlanID:         details.PlanID,
		InstanceID:     instanceID,
		TargetID:       freshInstance.ID,
		TargetName:     freshInstance.Name,
		TargetStatus:   freshInstance.Status,
		TargetInfo:     string(targetinfo),
		AdditionalInfo: "",
	}

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("create rds instance in back database opts: %v", idsOpts))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds instance in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create rds instance in back database succeed: %s", instanceID))

	// Return result
	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
