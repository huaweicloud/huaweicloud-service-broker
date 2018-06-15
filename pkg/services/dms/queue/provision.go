package queue

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/golangsdk/openstack/dms/v1/groups"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/queues"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/database"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// Provision implematation
func (b *DMSBroker) Provision(instanceID string, details brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {

	// Check dms instance length in back database
	var length int
	err := database.BackDBConnection.
		Model(&database.InstanceDetails{}).
		Where("instance_id = ? and service_id = ? and plan_id = ?", instanceID, details.ServiceID, details.PlanID).
		Count(&length).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("check dms instance length in back database failed. Error: %s", err)
	}
	// ErrInstanceAlreadyExists
	if length > 0 {
		return brokerapi.ProvisionedServiceSpec{}, brokerapi.ErrInstanceAlreadyExists
	}

	// Init dms client
	dmsClient, err := b.CloudCredentials.DMSV1Client()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dms client failed. Error: %s", err)
	}

	// Find service plan
	servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("find service plan failed. Error: %s", err)
	}

	// Get parameters from service plan metadata
	metadataParameters := MetadataParameters{}
	if servicePlan.Metadata != nil {
		if len(servicePlan.Metadata.Parameters) > 0 {
			err := json.Unmarshal(servicePlan.Metadata.Parameters, &metadataParameters)
			if err != nil {
				return brokerapi.ProvisionedServiceSpec{},
					fmt.Errorf("Error unmarshalling Parameters from service plan: %s", err)
			}
		}
	}

	// Get parameters from details
	provisionParameters := ProvisionParameters{}
	if len(details.RawParameters) > 0 {
		err := json.Unmarshal(details.RawParameters, &provisionParameters)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{},
				fmt.Errorf("Error unmarshalling rawParameters from details: %s", err)
		}
	}

	// Init provisionOpts
	provisionOpts := queues.CreateOps{}
	// queue name
	provisionOpts.Name = provisionParameters.QueueName
	// queue description
	provisionOpts.Description = provisionParameters.Description
	// queue mode
	provisionOpts.QueueMode = metadataParameters.QueueMode
	// Queue Type: DMSStandardServiceName
	if servicePlan.Name == models.DMSStandardServiceName {
		// Queue Mode: NORMAL, FIFO
		if metadataParameters.RedrivePolicy != "" {
			provisionOpts.RedrivePolicy = metadataParameters.RedrivePolicy
		}
		if metadataParameters.MaxConsumeCount > 0 {
			provisionOpts.MaxConsumeCount = metadataParameters.MaxConsumeCount
		}
		// Override metadataParameters
		if provisionParameters.RedrivePolicy != "" {
			provisionOpts.RedrivePolicy = provisionParameters.RedrivePolicy
		}
		if provisionParameters.MaxConsumeCount > 0 {
			provisionOpts.MaxConsumeCount = provisionParameters.MaxConsumeCount
		}
	}
	// Queue Type: DMSKafkaServiceName
	if servicePlan.Name == models.DMSKafkaServiceName {
		// Queue Mode: KAFKA_HA, KAFKA_HT
		if metadataParameters.RetentionHours > 0 {
			provisionOpts.RetentionHours = metadataParameters.RetentionHours
		}
		// Override metadataParameters
		if provisionParameters.RetentionHours > 0 {
			provisionOpts.RetentionHours = provisionParameters.RetentionHours
		}
	}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision dms queue opts: %v", provisionOpts))

	// Invoke sdk
	dmsQueue, err := queues.Create(dmsClient, provisionOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision dms queue failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision dms queue result: %v", dmsQueue))

	// Init provisionGroupOpts
	provisionGroupOpts := groups.CreateOps{}
	provisionGroupOpts.Groups = []groups.GroupOps{{Name: provisionParameters.GroupName}}

	// Log opts
	b.Logger.Debug(fmt.Sprintf("provision dms group opts: %v", provisionGroupOpts))

	// Invoke sdk
	dmsGroups, err := groups.Create(dmsClient, dmsQueue.ID, provisionGroupOpts).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("provision dms group failed. Error: %s", err)
	}

	// Log result
	b.Logger.Debug(fmt.Sprintf("provision dms group result: %v", dmsGroups))

	// Invoke sdk get
	freshQueue, err := queues.Get(dmsClient, dmsQueue.ID, false).Extract()
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("get dms queue failed. Error: %s", err)
	}

	// Marshal queue
	targetinfo, err := json.Marshal(freshQueue)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("marshal dms queue failed. Error: %s", err)
	}

	// create InstanceDetails in back database
	idsOpts := database.InstanceDetails{
		ServiceID:      details.ServiceID,
		PlanID:         details.PlanID,
		InstanceID:     instanceID,
		TargetID:       freshQueue.ID,
		TargetName:     freshQueue.Name,
		TargetStatus:   "",
		TargetInfo:     string(targetinfo),
		AdditionalInfo: "",
	}

	// log InstanceDetails opts
	b.Logger.Debug(fmt.Sprintf("create dms queue in back database opts: %v", idsOpts))

	err = database.BackDBConnection.Create(&idsOpts).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("create dms queue in back database failed. Error: %s", err)
	}

	// Log InstanceDetails result
	b.Logger.Debug(fmt.Sprintf("create dms queue in back database succeed: %s", instanceID))

	// Return result
	if asyncAllowed && models.OperationAsyncDMS {
		// OperationDatas for OperationProvisioning
		ods := models.OperationDatas{
			OperationType:  models.OperationProvisioning,
			ServiceID:      details.ServiceID,
			PlanID:         details.PlanID,
			InstanceID:     instanceID,
			TargetID:       freshQueue.ID,
			TargetName:     freshQueue.Name,
			TargetStatus:   "",
			TargetInfo:     string(targetinfo),
			AdditionalInfo: "",
		}

		operationdata, err := ods.ToString()
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("convert dms queue operation datas failed. Error: %s", err)
		}

		// log OperationDatas
		b.Logger.Debug(fmt.Sprintf("create dms queue operation datas: %s", operationdata))

		return brokerapi.ProvisionedServiceSpec{IsAsync: true, DashboardURL: "", OperationData: operationdata}, nil
	}

	return brokerapi.ProvisionedServiceSpec{IsAsync: false, DashboardURL: "", OperationData: ""}, nil
}
