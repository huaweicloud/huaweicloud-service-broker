package dcs

import (
	"encoding/json"
	"fmt"
	"strings"

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
					PropertySchemas: json.RawMessage(`{
						"password": {
							"type": "string",
							"description": "Password of a Redis instance. The password of a Redis instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as ` + fmt.Sprintf("%s", "`") + `~!@#$%^&*()-_=+|[{}]:'\",<.>/?)."
						},
						"name": {
							"type": "string",
							"description": "Redis instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters."
						},
						"description": {
							"type": "string",
							"description": "Brief description of the Redis instance. A brief description supports up to 1024 characters."
						},
						"capacity": {
							"type": "integer",
							"description": "Cache capacity. Unit: GB. For a Redis instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. For a Redis instance in cluster mode, the cache capacity can be 64, 128, 256, 512, or 1024 GB.",
							"default": ` + fmt.Sprintf("%d", metadataParameters.Capacity) + `
						},
						"vpc_id": {
							"type": "string",
							"description": "Tenant's VPC ID.",
							"default": "` + metadataParameters.VPCID + `"
						},
						"subnet_id": {
							"type": "string",
							"description": "Subnet ID.",
							"default": "` + metadataParameters.SubnetID + `"
						},
						"security_group_id": {
							"type": "string",
							"description": "Tenant's security group ID.",
							"default": "` + metadataParameters.SecurityGroupID + `"
						},
						"availability_zones": {
							"type": "string",
							"description": "IDs of the AZs where cache nodes reside. Note: In the current version of Orange Cloud, only one AZ ID can be set in the request.",
							"default": "` + strings.Replace(models.ToJson(metadataParameters.AvailabilityZones), `"`, `\"`, -1) + `"
						},
						"backup_strategy_savedays": {
							"type": "integer",
							"description": "Retention time. Unit: day. Range: 1–7."
						},
						"backup_strategy_backup_type": {
							"type": "string",
							"description": "Backup type. Options: auto: automatic backup. manual: manual backup."
						},
						"backup_strategy_backup_at": {
							"type": "string",
							"description": "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday."
						},
						"backup_strategy_begin_at": {
							"type": "string",
							"description": "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00."
						},
						"backup_strategy_period_type": {
							"type": "string",
							"description": "Interval at which backup is performed. Currently, only weekly backup is supported."
						},
						"maintain_begin": {
							"type": "string",
							"description": "Time at which the maintenance time window starts."
						},
						"maintain_end": {
							"type": "string",
							"description": "Time at which the maintenance time window ends."
						}
					}`),
				},
				UpdatingParametersSchema: &brokerapi.InputParametersSchema{
					PropertySchemas: json.RawMessage(`{
						"name": {
							"type": "string",
							"description": "Redis instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters."
						},
						"description": {
							"type": "string",
							"description": "Brief description of the Redis instance. A brief description supports up to 1024 characters."
						},
						"backup_strategy_savedays": {
							"type": "integer",
							"description": "Retention time. Unit: day. Range: 1–7."
						},
						"backup_strategy_backup_type": {
							"type": "string",
							"description": "Backup type. Options: auto: automatic backup. manual: manual backup."
						},
						"backup_strategy_backup_at": {
							"type": "string",
							"description": "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday."
						},
						"backup_strategy_begin_at": {
							"type": "string",
							"description": "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00."
						},
						"backup_strategy_period_type": {
							"type": "string",
							"description": "Interval at which backup is performed. Currently, only weekly backup is supported."
						},
						"maintain_begin": {
							"type": "string",
							"description": "Time at which the maintenance time window starts."
						},
						"maintain_end": {
							"type": "string",
							"description": "Time at which the maintenance time window ends."
						},
						"security_group_id": {
							"type": "string",
							"description": "Tenant's security group ID."
						},
						"new_capacity": {
							"type": "integer",
							"description": "New cache capacity. Unit: GB. For a Redis instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. For a Redis instance in cluster mode, it does not support extend the cache capacity."
						},
						"old_password": {
							"type": "string",
							"description": "The previous password of Redis instance."
						},
						"new_password": {
							"type": "string",
							"description": "The new password of Redis instance."
						}
					}`),
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
					PropertySchemas: json.RawMessage(`{
						"username": {
							"type": "string",
							"description": "Username of a Memcached instance."
						},
						"password": {
							"type": "string",
							"description": "Password of a Memcached instance. The password of a Memcached instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as ` + fmt.Sprintf("%s", "`") + `~!@#$%^&*()-_=+|[{}]:'\",<.>/?)."
						},
						"name": {
							"type": "string",
							"description": "Memcached instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters."
						},
						"description": {
							"type": "string",
							"description": "Brief description of the Memcached instance. A brief description supports up to 1024 characters."
						},
						"capacity": {
							"type": "integer",
							"description": "Cache capacity. Unit: GB. For a Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.",
							"default": ` + fmt.Sprintf("%d", metadataParameters.Capacity) + `
						},
						"vpc_id": {
							"type": "string",
							"description": "Tenant's VPC ID.",
							"default": "` + metadataParameters.VPCID + `"
						},
						"subnet_id": {
							"type": "string",
							"description": "Subnet ID.",
							"default": "` + metadataParameters.SubnetID + `"
						},
						"security_group_id": {
							"type": "string",
							"description": "Tenant's security group ID.",
							"default": "` + metadataParameters.SecurityGroupID + `"
						},
						"availability_zones": {
							"type": "string",
							"description": "IDs of the AZs where cache nodes reside.",
							"default": "` + strings.Replace(models.ToJson(metadataParameters.AvailabilityZones), `"`, `\"`, -1) + `"
						},
						"backup_strategy_savedays": {
							"type": "integer",
							"description": "Retention time. Unit: day. Range: 1–7."
						},
						"backup_strategy_backup_type": {
							"type": "string",
							"description": "Backup type. Options: auto: automatic backup. manual: manual backup."
						},
						"backup_strategy_backup_at": {
							"type": "string",
							"description": "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday."
						},
						"backup_strategy_begin_at": {
							"type": "string",
							"description": "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00."
						},
						"backup_strategy_period_type": {
							"type": "string",
							"description": "Interval at which backup is performed. Currently, only weekly backup is supported."
						},
						"maintain_begin": {
							"type": "string",
							"description": "Time at which the maintenance time window starts."
						},
						"maintain_end": {
							"type": "string",
							"description": "Time at which the maintenance time window ends."
						}
					}`),
				},
				UpdatingParametersSchema: &brokerapi.InputParametersSchema{
					PropertySchemas: json.RawMessage(`{
						"name": {
							"type": "string",
							"description": "Memcached instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters."
						},
						"description": {
							"type": "string",
							"description": "Brief description of the Memcached instance. A brief description supports up to 1024 characters."
						},
						"backup_strategy_savedays": {
							"type": "integer",
							"description": "Retention time. Unit: day. Range: 1–7."
						},
						"backup_strategy_backup_type": {
							"type": "string",
							"description": "Backup type. Options: auto: automatic backup. manual: manual backup."
						},
						"backup_strategy_backup_at": {
							"type": "string",
							"description": "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday."
						},
						"backup_strategy_begin_at": {
							"type": "string",
							"description": "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00."
						},
						"backup_strategy_period_type": {
							"type": "string",
							"description": "Interval at which backup is performed. Currently, only weekly backup is supported."
						},
						"maintain_begin": {
							"type": "string",
							"description": "Time at which the maintenance time window starts."
						},
						"maintain_end": {
							"type": "string",
							"description": "Time at which the maintenance time window ends."
						},
						"security_group_id": {
							"type": "string",
							"description": "Tenant's security group ID."
						},
						"new_capacity": {
							"type": "integer",
							"description": "New cache capacity. Unit: GB. For a Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. For a Redis instance in cluster mode, it does not support extend the cache capacity."
						},
						"old_password": {
							"type": "string",
							"description": "The previous password of Memcached instance."
						},
						"new_password": {
							"type": "string",
							"description": "The new password of Memcached instance."
						}
					}`),
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
					PropertySchemas: json.RawMessage(`{
						"username": {
							"type": "string",
							"description": "Username of a IMDG instance."
						},
						"password": {
							"type": "string",
							"description": "Password of a IMDG instance. The password of a IMDG instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as ` + fmt.Sprintf("%s", "`") + `~!@#$%^&*()-_=+|[{}]:'\",<.>/?)."
						},
						"name": {
							"type": "string",
							"description": "IMDG instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters."
						},
						"description": {
							"type": "string",
							"description": "Brief description of the IMDG instance. A brief description supports up to 1024 characters."
						},
						"capacity": {
							"type": "integer",
							"description": "Cache capacity. Unit: GB. For a IMDG instance in single node, the cache capacity can be 2 GB, 4 GB, 8 GB. For a IMDG instance in cluster mode, the cache capacity can be 64 GB.",
							"default": ` + fmt.Sprintf("%d", metadataParameters.Capacity) + `
						},
						"vpc_id": {
							"type": "string",
							"description": "Tenant's VPC ID.",
							"default": "` + metadataParameters.VPCID + `"
						},
						"subnet_id": {
							"type": "string",
							"description": "Subnet ID.",
							"default": "` + metadataParameters.SubnetID + `"
						},
						"security_group_id": {
							"type": "string",
							"description": "Tenant's security group ID.",
							"default": "` + metadataParameters.SecurityGroupID + `"
						},
						"availability_zones": {
							"type": "string",
							"description": "IDs of the AZs where cache nodes reside.",
							"default": "` + strings.Replace(models.ToJson(metadataParameters.AvailabilityZones), `"`, `\"`, -1) + `"
						},
						"backup_strategy_savedays": {
							"type": "integer",
							"description": "Retention time. Unit: day. Range: 1–7."
						},
						"backup_strategy_backup_type": {
							"type": "string",
							"description": "Backup type. Options: auto: automatic backup. manual: manual backup."
						},
						"backup_strategy_backup_at": {
							"type": "string",
							"description": "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday."
						},
						"backup_strategy_begin_at": {
							"type": "string",
							"description": "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00."
						},
						"backup_strategy_period_type": {
							"type": "string",
							"description": "Interval at which backup is performed. Currently, only weekly backup is supported."
						},
						"maintain_begin": {
							"type": "string",
							"description": "Time at which the maintenance time window starts."
						},
						"maintain_end": {
							"type": "string",
							"description": "Time at which the maintenance time window ends."
						}
					}`),
				},
				UpdatingParametersSchema: &brokerapi.InputParametersSchema{
					PropertySchemas: json.RawMessage(`{
						"name": {
							"type": "string",
							"description": "IMDG instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters."
						},
						"description": {
							"type": "string",
							"description": "Brief description of the IMDG instance. A brief description supports up to 1024 characters."
						},
						"backup_strategy_savedays": {
							"type": "integer",
							"description": "Retention time. Unit: day. Range: 1–7."
						},
						"backup_strategy_backup_type": {
							"type": "string",
							"description": "Backup type. Options: auto: automatic backup. manual: manual backup."
						},
						"backup_strategy_backup_at": {
							"type": "string",
							"description": "Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday."
						},
						"backup_strategy_begin_at": {
							"type": "string",
							"description": "Time at which backup starts. \"00:00-01:00\" indicates that backup starts at 00:00:00."
						},
						"backup_strategy_period_type": {
							"type": "string",
							"description": "Interval at which backup is performed. Currently, only weekly backup is supported."
						},
						"maintain_begin": {
							"type": "string",
							"description": "Time at which the maintenance time window starts."
						},
						"maintain_end": {
							"type": "string",
							"description": "Time at which the maintenance time window ends."
						},
						"security_group_id": {
							"type": "string",
							"description": "Tenant's security group ID."
						}
					}`),
				},
			},
			ServiceBindings: nil,
		}
	}

	b.Logger.Debug(fmt.Sprintf("get dcs schemas: %v", models.ToJson(schemas)))

	return &schemas, nil
}
