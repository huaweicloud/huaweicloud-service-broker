package rds

import (
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/rds/v1/instances"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *RDSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Init rds client
	rdsClient, err := b.CloudCredentials.RDSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create rds client failed. Error: %s", err)
	}

	// Init provisionOpts
	provisionOpts := instances.CreateOps{
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

	// Return result
	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
