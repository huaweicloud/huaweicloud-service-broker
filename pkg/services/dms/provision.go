package dms

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

	// Init provisionOpts
	provisionOpts := queues.CreateOps{}
	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &provisionOpts)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error unmarshalling rawParameters: %s", err)
		}
	}

	// Find service plan
	servicePlan, err := b.Catalog.FindServicePlan(details.ServiceID, details.PlanID)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("find service plan failed. Error: %s", err)
	}

	// Setting provisionOpts
	if provisionOpts.Name == "" {
		provisionOpts.Name = instanceID
	}

	// TODO need to confirm different queue mode
	if servicePlan.Name == models.DMSStandardServiceName {
		// Queue Type: DMSStandardServiceName
		// Queue Mode: NORMAL, FIFO
		provisionOpts.QueueMode = "NORMAL"
	} else if servicePlan.Name == models.DMSKafkaServiceName {
		// Queue Type: DMSKafkaServiceName
		// Queue Mode: KAFKA_HA, KAFKA_HT
		provisionOpts.QueueMode = "KAFKA_HT"
	} else if servicePlan.Name == models.DMSActiveMQServiceName {
		// Queue Type: DMSActiveMQServiceName
		// Queue Mode: AMQP
		provisionOpts.QueueMode = "AMQP"
	} else if servicePlan.Name == models.DMSRabbitMQServiceName {
		// TODO need to invoke another instance interface
	} else {
		return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("unknown service name: %s", servicePlan.Name)
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
	if len(details.RawParameters) >= 0 {
		err := json.Unmarshal(details.RawParameters, &provisionGroupOpts)
		if err != nil {
			return brokerapi.ProvisionedServiceSpec{}, fmt.Errorf("Error unmarshalling rawParameters: %s", err)
		}
	}

	// Setting provisionGroupOpts
	if len(provisionGroupOpts.Groups) == 0 {
		provisionGroupOpts.Groups = []groups.GroupOps{{Name: instanceID}}
	}

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
