package obs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// GetPlanSchemas implematation
func (b *OBSBroker) GetPlanSchemas(serviceID string, planID string, metadata *brokerapi.ServicePlanMetadata) (*brokerapi.PlanSchemas, error) {
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

	b.Logger.Debug(fmt.Sprintf("get obs metadata parameters: %v", models.ToJson(metadataParameters)))

	// Build schemas
	schemas := brokerapi.PlanSchemas{
		ServiceInstances: brokerapi.InstanceSchemas{
			ProvisioningParametersSchema: brokerapi.InputParametersSchema{
				RequiredProperties: []string{
					"bucket_name",
				},
				PropertySchemas: map[string]brokerapi.PropertySchema{
					"bucket_name": &brokerapi.StringPropertySchema{
						Description: "Enter the bucket name, which must be globally unique. Name the bucket according to the globally applied DNS naming regulation as follows: Must contain 3 to 63 characters, including lowercase letters, digits, hyphens (-), and periods (.) Cannot be an IP address. Cannot start or end with a hyphen (-) or period (.) Cannot contain two consecutive periods (.), for example, my..bucket. Cannot contain periods (.) and hyphens (-) adjacent to each other, for example, my-.bucket or my.-bucket.",
					},
					"bucket_policy": &brokerapi.StringPropertySchema{
						Description: "A bucket policy defines the access control policy of resources (buckets and objects) on OBS. private: Only the bucket owner can read, write, and delete objects in the bucket. This policy is the default bucket policy. public-read: Any user can read objects in the bucket. Only the bucket owner can write and delete objects in the bucket. public-read-write: Any user can read, write, and delete objects in the bucket.",
					},
				},
			},
			UpdatingParametersSchema: &brokerapi.InputParametersSchema{
				PropertySchemas: map[string]brokerapi.PropertySchema{
					"bucket_policy": &brokerapi.StringPropertySchema{
						Description: "A bucket policy defines the access control policy of resources (buckets and objects) on OBS. private: Only the bucket owner can read, write, and delete objects in the bucket. This policy is the default bucket policy. public-read: Any user can read objects in the bucket. Only the bucket owner can write and delete objects in the bucket. public-read-write: Any user can read, write, and delete objects in the bucket.",
					},
					"status": &brokerapi.StringPropertySchema{
						Description: "By default, the versioning function is disabled for new buckets on OBS. The status include: Enabled and Suspended.",
					},
				},
			},
		},
		ServiceBindings: nil,
	}

	b.Logger.Debug(fmt.Sprintf("get obs schemas: %v", models.ToJson(schemas)))

	return &schemas, nil
}
