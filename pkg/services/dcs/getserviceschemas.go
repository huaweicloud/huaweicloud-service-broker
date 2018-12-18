package dcs

import (
	"encoding/json"
	"fmt"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/models"
	"github.com/pivotal-cf/brokerapi"
)

// GetPlanSchemas implematation
func (b *DCSBroker) GetPlanSchemas(serviceID string, planID string, metadata *brokerapi.ServicePlanMetadata) (*brokerapi.PlanSchemas, error) {
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

	b.Logger.Debug(fmt.Sprintf("get dcs metadata parameters: %v", models.ToJson(metadataParameters)))

	service, err := b.Catalog.FindService(serviceID)
	if err != nil {
		return nil, fmt.Errorf("find dcs queue service failed. Error: %s", err)
	}

	// Build schemas
	schemas := brokerapi.PlanSchemas{}
	if service.Name == models.DCSRedisServiceName {
		schemas = brokerapi.PlanSchemas{
			ServiceInstances: brokerapi.InstanceSchemas{
				ProvisioningParametersSchema: brokerapi.InputParametersSchema{
					RequiredProperties: []string{
						"password",
						"name",
					},
					PropertySchemas: map[string]brokerapi.PropertySchema{
						"password": &brokerapi.StringPropertySchema{
							Description: "Password of a Redis instance. The password of a Redis instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as `~!@#$%^&*()-_=+|[{}]:'\",<.>/?.",
						},
						"name": &brokerapi.StringPropertySchema{
							Description: "Redis instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.",
						},
						"description": &brokerapi.StringPropertySchema{
							Description: "Brief description of the Redis instance. A brief description supports up to 1024 characters.",
						},
						"capacity": &brokerapi.IntPropertySchema{
							Description: "Cache capacity. Unit: GB. For a Redis instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. For a Redis instance in cluster mode, the cache capacity can be 64, 128, 256, 512, or 1024 GB.",
						},
						"vpc_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's VPC ID.",
						},
						"subnet_id": &brokerapi.StringPropertySchema{
							Description: "Subnet ID.",
						},
						"security_group_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's security group ID.",
						},
						"availability_zones": &brokerapi.ArrayPropertySchema{
							Description: "IDs of the AZs where cache nodes reside. Note: In the current version of Orange Cloud, only one AZ ID can be set in the request.",
							ItemsSchema: &brokerapi.StringPropertySchema{
								Description: "AZ ID.",
							},
						},
						"backup_strategy_savedays": &brokerapi.IntPropertySchema{
							Description: "Retention time. Unit: day. Range: 1–7.",
						},
						"backup_strategy_backup_type": &brokerapi.StringPropertySchema{
							Description: "Backup type. Options: auto: automatic backup. manual: manual backup.",
						},
						"backup_strategy_backup_at": &brokerapi.ArrayPropertySchema{
							Description: "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.",
							ItemsSchema: &brokerapi.IntPropertySchema{
								Description: "Day in a week on which backup starts.",
							},
						},
						"backup_strategy_begin_at": &brokerapi.StringPropertySchema{
							Description: "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00.",
						},
						"backup_strategy_period_type": &brokerapi.StringPropertySchema{
							Description: "Interval at which backup is performed. Currently, only weekly backup is supported.",
						},
						"maintain_begin": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window starts.",
						},
						"maintain_end": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window ends.",
						},
					},
				},
				UpdatingParametersSchema: &brokerapi.InputParametersSchema{
					PropertySchemas: map[string]brokerapi.PropertySchema{
						"name": &brokerapi.StringPropertySchema{
							Description: "Redis instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.",
						},
						"description": &brokerapi.StringPropertySchema{
							Description: "Brief description of the Redis instance. A brief description supports up to 1024 characters.",
						},
						"backup_strategy_savedays": &brokerapi.IntPropertySchema{
							Description: "Retention time. Unit: day. Range: 1–7.",
						},
						"backup_strategy_backup_type": &brokerapi.StringPropertySchema{
							Description: "Backup type. Options: auto: automatic backup. manual: manual backup.",
						},
						"backup_strategy_backup_at": &brokerapi.ArrayPropertySchema{
							Description: "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.",
							ItemsSchema: &brokerapi.IntPropertySchema{
								Description: "Day in a week on which backup starts.",
							},
						},
						"backup_strategy_begin_at": &brokerapi.StringPropertySchema{
							Description: "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00.",
						},
						"backup_strategy_period_type": &brokerapi.StringPropertySchema{
							Description: "Interval at which backup is performed. Currently, only weekly backup is supported.",
						},
						"maintain_begin": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window starts.",
						},
						"maintain_end": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window ends.",
						},
						"security_group_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's security group ID.",
						},
						"new_capacity": &brokerapi.IntPropertySchema{
							Description: "New cache capacity. Unit: GB. For a Redis instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. For a Redis instance in cluster mode, it does not support extend the cache capacity.",
						},
						"old_password": &brokerapi.StringPropertySchema{
							Description: "The previous password of Redis instance.",
						},
						"new_password": &brokerapi.StringPropertySchema{
							Description: "The new password of Redis instance.",
						},
					},
				},
			},
			ServiceBindings: nil,
		}
	} else if service.Name == models.DCSMemcachedServiceName {
		schemas = brokerapi.PlanSchemas{
			ServiceInstances: brokerapi.InstanceSchemas{
				ProvisioningParametersSchema: brokerapi.InputParametersSchema{
					RequiredProperties: []string{
						"username",
						"password",
						"name",
					},
					PropertySchemas: map[string]brokerapi.PropertySchema{
						"username": &brokerapi.StringPropertySchema{
							Description: "Username of a Memcached instance.",
						},
						"password": &brokerapi.StringPropertySchema{
							Description: "Password of a Memcached instance. The password of a Memcached instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as `~!@#$%^&*()-_=+|[{}]:'\",<.>/?.",
						},
						"name": &brokerapi.StringPropertySchema{
							Description: "Memcached instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.",
						},
						"description": &brokerapi.StringPropertySchema{
							Description: "Brief description of the Memcached instance. A brief description supports up to 1024 characters.",
						},
						"capacity": &brokerapi.IntPropertySchema{
							Description: "Cache capacity. Unit: GB. For a Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.",
						},
						"vpc_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's VPC ID.",
						},
						"subnet_id": &brokerapi.StringPropertySchema{
							Description: "Subnet ID.",
						},
						"security_group_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's security group ID.",
						},
						"availability_zones": &brokerapi.ArrayPropertySchema{
							Description: "IDs of the AZs where cache nodes reside.",
							ItemsSchema: &brokerapi.StringPropertySchema{
								Description: "AZ ID.",
							},
						},
						"backup_strategy_savedays": &brokerapi.IntPropertySchema{
							Description: "Retention time. Unit: day. Range: 1–7.",
						},
						"backup_strategy_backup_type": &brokerapi.StringPropertySchema{
							Description: "Backup type. Options: auto: automatic backup. manual: manual backup.",
						},
						"backup_strategy_backup_at": &brokerapi.ArrayPropertySchema{
							Description: "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.",
							ItemsSchema: &brokerapi.IntPropertySchema{
								Description: "Day in a week on which backup starts.",
							},
						},
						"backup_strategy_begin_at": &brokerapi.StringPropertySchema{
							Description: "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00.",
						},
						"backup_strategy_period_type": &brokerapi.StringPropertySchema{
							Description: "Interval at which backup is performed. Currently, only weekly backup is supported.",
						},
						"maintain_begin": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window starts.",
						},
						"maintain_end": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window ends.",
						},
					},
				},
				UpdatingParametersSchema: &brokerapi.InputParametersSchema{
					PropertySchemas: map[string]brokerapi.PropertySchema{
						"name": &brokerapi.StringPropertySchema{
							Description: "Memcached instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.",
						},
						"description": &brokerapi.StringPropertySchema{
							Description: "Brief description of the Memcached instance. A brief description supports up to 1024 characters.",
						},
						"backup_strategy_savedays": &brokerapi.IntPropertySchema{
							Description: "Retention time. Unit: day. Range: 1–7.",
						},
						"backup_strategy_backup_type": &brokerapi.StringPropertySchema{
							Description: "Backup type. Options: auto: automatic backup. manual: manual backup.",
						},
						"backup_strategy_backup_at": &brokerapi.ArrayPropertySchema{
							Description: "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.",
							ItemsSchema: &brokerapi.IntPropertySchema{
								Description: "Day in a week on which backup starts.",
							},
						},
						"backup_strategy_begin_at": &brokerapi.StringPropertySchema{
							Description: "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00.",
						},
						"backup_strategy_period_type": &brokerapi.StringPropertySchema{
							Description: "Interval at which backup is performed. Currently, only weekly backup is supported.",
						},
						"maintain_begin": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window starts.",
						},
						"maintain_end": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window ends.",
						},
						"security_group_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's security group ID.",
						},
						"new_capacity": &brokerapi.IntPropertySchema{
							Description: "New cache capacity. Unit: GB. For a Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. For a Redis instance in cluster mode, it does not support extend the cache capacity.",
						},
						"old_password": &brokerapi.StringPropertySchema{
							Description: "The previous password of Memcached instance.",
						},
						"new_password": &brokerapi.StringPropertySchema{
							Description: "The new password of Memcached instance.",
						},
					},
				},
			},
			ServiceBindings: nil,
		}
	} else if service.Name == models.DCSIMDGServiceName {
		schemas = brokerapi.PlanSchemas{
			ServiceInstances: brokerapi.InstanceSchemas{
				ProvisioningParametersSchema: brokerapi.InputParametersSchema{
					RequiredProperties: []string{
						"username",
						"password",
						"name",
					},
					PropertySchemas: map[string]brokerapi.PropertySchema{
						"username": &brokerapi.StringPropertySchema{
							Description: "Username of a IMDG instance.",
						},
						"password": &brokerapi.StringPropertySchema{
							Description: "Password of a IMDG instance. The password of a IMDG instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as `~!@#$%^&*()-_=+|[{}]:'\",<.>/?.",
						},
						"name": &brokerapi.StringPropertySchema{
							Description: "IMDG instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.",
						},
						"description": &brokerapi.StringPropertySchema{
							Description: "Brief description of the IMDG instance. A brief description supports up to 1024 characters.",
						},
						"capacity": &brokerapi.IntPropertySchema{
							Description: "Cache capacity. Unit: GB. For a IMDG instance in single node, the cache capacity can be 2 GB, 4 GB, 8 GB. For a IMDG instance in cluster mode, the cache capacity can be 64 GB.",
						},
						"vpc_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's VPC ID.",
						},
						"subnet_id": &brokerapi.StringPropertySchema{
							Description: "Subnet ID.",
						},
						"security_group_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's security group ID.",
						},
						"availability_zones": &brokerapi.ArrayPropertySchema{
							Description: "IDs of the AZs where cache nodes reside.",
							ItemsSchema: &brokerapi.StringPropertySchema{
								Description: "AZ ID.",
							},
						},
						"backup_strategy_savedays": &brokerapi.IntPropertySchema{
							Description: "Retention time. Unit: day. Range: 1–7.",
						},
						"backup_strategy_backup_type": &brokerapi.StringPropertySchema{
							Description: "Backup type. Options: auto: automatic backup. manual: manual backup.",
						},
						"backup_strategy_backup_at": &brokerapi.ArrayPropertySchema{
							Description: "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.",
							ItemsSchema: &brokerapi.IntPropertySchema{
								Description: "Day in a week on which backup starts.",
							},
						},
						"backup_strategy_begin_at": &brokerapi.StringPropertySchema{
							Description: "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00.",
						},
						"backup_strategy_period_type": &brokerapi.StringPropertySchema{
							Description: "Interval at which backup is performed. Currently, only weekly backup is supported.",
						},
						"maintain_begin": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window starts.",
						},
						"maintain_end": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window ends.",
						},
					},
				},
				UpdatingParametersSchema: &brokerapi.InputParametersSchema{
					PropertySchemas: map[string]brokerapi.PropertySchema{
						"name": &brokerapi.StringPropertySchema{
							Description: "IMDG instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.",
						},
						"description": &brokerapi.StringPropertySchema{
							Description: "Brief description of the IMDG instance. A brief description supports up to 1024 characters.",
						},
						"backup_strategy_savedays": &brokerapi.IntPropertySchema{
							Description: "Retention time. Unit: day. Range: 1–7.",
						},
						"backup_strategy_backup_type": &brokerapi.StringPropertySchema{
							Description: "Backup type. Options: auto: automatic backup. manual: manual backup.",
						},
						"backup_strategy_backup_at": &brokerapi.ArrayPropertySchema{
							Description: "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.",
							ItemsSchema: &brokerapi.IntPropertySchema{
								Description: "Day in a week on which backup starts.",
							},
						},
						"backup_strategy_begin_at": &brokerapi.StringPropertySchema{
							Description: "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00.",
						},
						"backup_strategy_period_type": &brokerapi.StringPropertySchema{
							Description: "Interval at which backup is performed. Currently, only weekly backup is supported.",
						},
						"maintain_begin": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window starts.",
						},
						"maintain_end": &brokerapi.StringPropertySchema{
							Description: "Time at which the maintenance time window ends.",
						},
						"security_group_id": &brokerapi.StringPropertySchema{
							Description: "Tenant's security group ID.",
						},
					},
				},
			},
			ServiceBindings: nil,
		}
	}

	b.Logger.Debug(fmt.Sprintf("get dcs schemas: %v", models.ToJson(schemas)))

	return &schemas, nil
}
