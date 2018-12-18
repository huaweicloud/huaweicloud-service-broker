package instance

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

	b.Logger.Debug(fmt.Sprintf("get dms instance metadata parameters: %v", models.ToJson(metadataParameters)))

	// Build schemas
	schemas := brokerapi.PlanSchemas{
		ServiceInstances: brokerapi.InstanceSchemas{
			ProvisioningParametersSchema: brokerapi.InputParametersSchema{
				RequiredProperties: []string{
					"username",
					"password",
					"name",
				},
				PropertySchemas: map[string]brokerapi.PropertySchema{
					"username": &brokerapi.StringPropertySchema{
						Description: "Indicates a username. A username consists of 1 to 64 characters and supports only letters, digits, and hyphens (-).",
					},
					"password": &brokerapi.StringPropertySchema{
						Description: "Indicates the password of an instance. An instance password must meet the following complexity requirements: Must be 6 to 32 characters long. Must contain at least two of the following character types: Lowercase letters Uppercase letters Digits Special characters (`~!@#$%^&*()-_=+|[{}]:'\",<.>/?)",
					},
					"name": &brokerapi.StringPropertySchema{
						Description: "Indicates the name of an instance. An instance name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).",
					},
					"description": &brokerapi.StringPropertySchema{
						Description: "Indicates the description of an instance. It is a character string containing not more than 1024 characters.",
					},
					"vpc_id": &brokerapi.StringPropertySchema{
						Description: "Indicates the ID of a VPC.",
					},
					"subnet_id": &brokerapi.StringPropertySchema{
						Description: "Indicates the ID of a subnet.",
					},
					"security_group_id": &brokerapi.StringPropertySchema{
						Description: "Indicates the ID of a security group.",
					},
					"availability_zones": &brokerapi.ArrayPropertySchema{
						Description: "Indicates the ID of an AZ.",
						ItemsSchema: &brokerapi.StringPropertySchema{
							Description: "AZ ID.",
						},
					},
					"maintain_begin": &brokerapi.StringPropertySchema{
						Description: "Indicates the time at which a maintenance time window starts. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The start time must be set to 22:00:00, 02:00:00, 06:00:00, 10:00:00, 14:00:00, or 18:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_begin is left blank, parameter maintain_end is also blank. In this case, the system automatically allocates the default start time 02:00:00.",
					},
					"maintain_end": &brokerapi.StringPropertySchema{
						Description: "Indicates the time at which a maintenance time window ends. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The end time is four hours later than the start time. For example, if the start time is 22:00:00, the end time is 02:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_end is left blank, parameter maintain_begin is also blank. In this case, the system automatically allocates the default end time 06:00:00.",
					},
				},
			},
			UpdatingParametersSchema: &brokerapi.InputParametersSchema{
				PropertySchemas: map[string]brokerapi.PropertySchema{
					"name": &brokerapi.StringPropertySchema{
						Description: "Indicates the name of an instance. An instance name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).",
					},
					"description": &brokerapi.StringPropertySchema{
						Description: "Indicates the description of an instance. It is a character string containing not more than 1024 characters.",
					},
					"maintain_begin": &brokerapi.StringPropertySchema{
						Description: "Indicates the time at which a maintenance time window starts. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The start time must be set to 22:00:00, 02:00:00, 06:00:00, 10:00:00, 14:00:00, or 18:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_begin is left blank, parameter maintain_end is also blank. In this case, the system automatically allocates the default start time 02:00:00.",
					},
					"maintain_end": &brokerapi.StringPropertySchema{
						Description: "Indicates the time at which a maintenance time window ends. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The end time is four hours later than the start time. For example, if the start time is 22:00:00, the end time is 02:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_end is left blank, parameter maintain_begin is also blank. In this case, the system automatically allocates the default end time 06:00:00.",
					},
					"security_group_id": &brokerapi.StringPropertySchema{
						Description: "Indicates the ID of a security group.",
					},
				},
			},
		},
		ServiceBindings: nil,
	}

	b.Logger.Debug(fmt.Sprintf("get dms instance schemas: %v", models.ToJson(schemas)))

	return &schemas, nil
}
