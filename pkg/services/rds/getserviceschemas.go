package rds

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// GetPlanSchemas implematation
func (b *RDSBroker) GetPlanSchemas(serviceID string, planID string, metadata *brokerapi.ServicePlanMetadata) (*brokerapi.PlanSchemas, error) {
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

	b.Logger.Debug(fmt.Sprintf("get rds metadata parameters: %v", models.ToJson(metadataParameters)))

	// Build schemas
	schemas := brokerapi.PlanSchemas{
		ServiceInstances: brokerapi.InstanceSchemas{
			ProvisioningParametersSchema: brokerapi.InputParametersSchema{
				RequiredProperties: []string{
					"name",
					"database_password",
				},
				PropertySchemas: map[string]brokerapi.PropertySchema{
					"name": &brokerapi.StringPropertySchema{
						Description: "Specifies the DB instance name. The DB instance name of the same DB engine is unique for the same tenant. Valid value: The value must be 4 to 64 characters in length and start with a letter. It is case-insensitive and can contain only letters, digits, hyphens (-), and underscores (_).",
					},
					"database_password": &brokerapi.StringPropertySchema{
						Description: "Specifies the password for user root of the database. Valid value: The value cannot be empty and should contain 8 to 32 characters, including uppercase and lowercase letters, digits, and the following special characters: ~!@#%^*-_=+?",
					},
					"speccode": &brokerapi.StringPropertySchema{
						Description: "Indicates the resource specifications code. Use rds.mysql.m1.xlarge as an example. rds indicates RDS, mysql indicates the DB engine, and m1.xlarge indicates the performance specification (large-memory). The parameter containing rr indicates the read replica specifications. The parameter not containing rr indicates the single or primary/standby DB instance specifications. If you enable HA, the suffix .ha need be added to the DB instance name. For example, the DB instance name is rds.db.s1.xlarge.ha.",
					},
					"volume_type": &brokerapi.StringPropertySchema{
						Description: "Specifies the volume type. Valid value: It must be COMMON (SATA) or ULTRAHIGH (SSD) and is case-sensitive.",
					},
					"volume_size": &brokerapi.IntPropertySchema{
						Description: "Specifies the volume size. Its value must be a multiple of 10 and the value range is 100 GB to 2000 GB.",
					},
					"availability_zone": &brokerapi.StringPropertySchema{
						Description: "Specifies the ID of the AZ.",
					},
					"vpc_id": &brokerapi.StringPropertySchema{
						Description: "Specifies the VPC ID.",
					},
					"subnet_id": &brokerapi.StringPropertySchema{
						Description: "Specifies the UUID for nics information.",
					},
					"security_group_id": &brokerapi.StringPropertySchema{
						Description: "Specifies the security group ID which the RDS DB instance belongs to.",
					},
					"database_port": &brokerapi.StringPropertySchema{
						Description: "Specifies the database port number.",
					},
					"backup_strategy_starttime": &brokerapi.StringPropertySchema{
						Description: "Indicates the backup start time that has been set. The backup task will be triggered within one hour after the backup start time.",
					},
					"backup_strategy_keepdays": &brokerapi.IntPropertySchema{
						Description: "Specifies the number of days to retain the generated backup files. Its value range is 0 to 35.",
					},
					"ha_enable": &brokerapi.StringPropertySchema{
						Description:   "Specifies the HA configuration parameter. Valid value: The value is true or false. The value true indicates creating HA DB instances. The value false indicates creating a single DB instance.",
						AllowedValues: []string{"true", "false"},
						DefaultValue:  "false",
					},
					"ha_replicationmode": &brokerapi.StringPropertySchema{
						Description: "Specifies the replication mode for the standby DB instance. The value cannot be empty. For MySQL, the value is async or semisync.",
					},
				},
			},
			UpdatingParametersSchema: &brokerapi.InputParametersSchema{
				PropertySchemas: map[string]brokerapi.PropertySchema{
					"volume_size": &brokerapi.IntPropertySchema{
						Description: "Specifies the volume size. Its value must be a multiple of 10 and the value range is 100 GB to 2000 GB.",
					},
					"speccode": &brokerapi.StringPropertySchema{
						Description: "Indicates the resource specifications code. Use rds.mysql.m1.xlarge as an example. rds indicates RDS, mysql indicates the DB engine, and m1.xlarge indicates the performance specification (large-memory). The parameter containing rr indicates the read replica specifications. The parameter not containing rr indicates the single or primary/standby DB instance specifications. If you enable HA, the suffix .ha need be added to the DB instance name. For example, the DB instance name is rds.db.s1.xlarge.ha.",
					},
				},
			},
		},
		ServiceBindings: nil,
	}

	b.Logger.Debug(fmt.Sprintf("get rds schemas: %v", models.ToJson(schemas)))

	return &schemas, nil
}
