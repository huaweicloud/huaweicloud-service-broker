package queue

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// GetPlanSchemas implematation
func (b *DMSBroker) GetPlanSchemas(serviceID string, planID string, metadata *brokerapi.ServicePlanMetadata) (*brokerapi.PlanSchemas, error) {
	// Get parameters from service plan metadata
	metadataParameters := MetadataParameters{}
	if metadata != nil {
		if len(metadata.Parameters) > 0 {
			err := json.Unmarshal(metadata.Parameters, &metadataParameters)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling Parameters from service plan: %s", err)
			}
		}
	}

	b.Logger.Debug(fmt.Sprintf("get dms queue metadata parameters: %v", models.ToJson(metadataParameters)))

	service, err := b.Catalog.FindService(serviceID)
	if err != nil {
		return nil, fmt.Errorf("find dms queue service failed. Error: %s", err)
	}

	// Build schemas
	schemas := brokerapi.PlanSchemas{}
	if service.Name == models.DMSStandardServiceName {
		schemas = brokerapi.PlanSchemas{
			ServiceInstances: brokerapi.InstanceSchemas{
				ProvisioningParametersSchema: brokerapi.InputParametersSchema{
					RequiredProperties: []string{
						"queue_name",
						"group_name",
					},
					PropertySchemas: json.RawMessage(`{
						"queue_name": {
							"type": "string",
							"description": "Indicates the name of a queue. A queue name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-)."
						},
						"group_name": {
							"type": "string",
							"description": "Indicates the name of a consumer group. A string of 1 to 32 characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_)."
						},
						"description": {
							"type": "string",
							"description": "Indicates the description of an instance. It is a character string containing not more than 1024 characters."
						},
						"redrive_policy": {
							"type": "string",
							"description": "This parameter is mandatory only when queue_mode is NORMAL or FIFO. Indicates whether to enable dead letter messages. Dead letter messages indicate messages that cannot be normally consumed. If a message fails to be consumed after the number of consumption attempts of this message reaches the maximum value, DMS stores this message into the dead letter queue. This message will be retained in the deal letter queue for 72 hours. During this period, consumers can consume the dead letter message. Dead letter messages can be consumed only by the consumer group that generates these dead letter messages. Dead letter messages of a FIFO queue are stored and consumed based on the FIFO sequence. Options: enable disable.",
							"default": "` + metadataParameters.RedrivePolicy + `"
						},
						"max_consume_count": {
							"type": "integer",
							"description": "This parameter is mandatory only when redrive_policy is set to enable. This parameter indicates the maximum number of allowed message consumption failures. When a message fails to be consumed after the number of consumption attempts of this message reaches this value, DMS stores this message into the dead letter queue. Value range: 1-100."
						}
					}`),
				},
				UpdatingParametersSchema: nil,
			},
			ServiceBindings: nil,
		}
	} else if service.Name == models.DMSActiveMQServiceName {
		schemas = brokerapi.PlanSchemas{
			ServiceInstances: brokerapi.InstanceSchemas{
				ProvisioningParametersSchema: brokerapi.InputParametersSchema{
					RequiredProperties: []string{
						"queue_name",
						"group_name",
					},
					PropertySchemas: json.RawMessage(`{
						"queue_name": {
							"type": "string",
							"description": "Indicates the name of a queue. A queue name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-)."
						},
						"group_name": {
							"type": "string",
							"description": "Indicates the name of a consumer group. A string of 1 to 32 characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_)."
						},
						"description": {
							"type": "string",
							"description": "Indicates the description of an instance. It is a character string containing not more than 1024 characters."
						}
					}`),
				},
				UpdatingParametersSchema: nil,
			},
			ServiceBindings: nil,
		}
	} else if service.Name == models.DMSKafkaServiceName {
		schemas = brokerapi.PlanSchemas{
			ServiceInstances: brokerapi.InstanceSchemas{
				ProvisioningParametersSchema: brokerapi.InputParametersSchema{
					RequiredProperties: []string{
						"queue_name",
						"group_name",
					},
					PropertySchemas: json.RawMessage(`{
						"queue_name": {
							"type": "string",
							"description": "Indicates the name of a queue. A queue name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-)."
						},
						"group_name": {
							"type": "string",
							"description": "Indicates the name of a consumer group. A string of 1 to 32 characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_)."
						},
						"description": {
							"type": "string",
							"description": "Indicates the description of an instance. It is a character string containing not more than 1024 characters."
						},
						"retention_hours": {
							"type": "integer",
							"description": "This parameter is mandatory only when queue_mode is set to KAFKA_HA or KAFKA_HT. This parameter indicates the retention time of messages in Kafka queues. Value range: 1 to 72 hours.",
							"default": ` + fmt.Sprintf("%d", metadataParameters.RetentionHours) + `
						}
					}`),
				},
				UpdatingParametersSchema: nil,
			},
			ServiceBindings: nil,
		}
	}

	b.Logger.Debug(fmt.Sprintf("get dms queue schemas: %v", models.ToJson(schemas)))

	return &schemas, nil
}
