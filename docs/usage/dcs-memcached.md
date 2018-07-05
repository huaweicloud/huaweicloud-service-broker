# Distributed Cache Service for Memcached

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dcs-memcached                  | Distributed Cache Service for Memcached

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| SingleNode                     | Memcached Single Node
| MasterStandby                  | Memcached Master Standby

## Provision

Provision a new instance for Distributed Cache Service for Memcached.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| capacity                     | int        | N         | Cache capacity. Unit: GB. For a Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB. The default value is in the config file.
| vpc_id                       | string     | N         | Tenant's VPC ID. The default value is in the config file.
| subnet_id                    | string     | N         | Subnet ID. The default value is in the config file.
| security_group_id            | string     | N         | Tenant's security group ID. The default value is in the config file.
| availability_zones           | []string   | N         | IDs of the AZs where cache nodes reside. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint). The default value is in the config file.
| username                     | string     | Y         | Username of a Memcached instance.
| password                     | string     | Y         | Password of a Memcached instance. The password of a Memcached instance must meet the following complexity requirements: A string of 6–32 characters. Contains at least two of the following character types: Uppercase letters; Lowercase letters; Digits; Special characters, such as `~!@#$%^&*()-_=+\|[{}]:'",<.>/?.
| name                         | string     | Y         | Memcached instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.
| description                  | string     | N         | Brief description of the Memcached instance. A brief description supports up to 1024 characters.
| backup_strategy_savedays     | int        | N         | Retention time. Unit: day. Range: 1–7.
| backup_strategy_backup_type  | string     | N         | Backup type. Options: auto: automatic backup. manual: manual backup.
| backup_strategy_backup_at    | []int      | N         | Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.
| backup_strategy_begin_at     | string     | N         | Time at which backup starts. "00:00-01:00" indicates that backup starts at 00:00:00.
| backup_strategy_period_type  | string     | N         | Interval at which backup is performed. Currently, only weekly backup is supported.
| maintain_begin               | string     | N         | Time at which the maintenance time window starts.
| maintain_end                 | string     | N         | Time at which the maintenance time window ends.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| host                   | string     | The fully-qualified address of Memcached instance.
| port                   | int        | The port number to connect to Memcached instance.
| username               | string     | Username of a Memcached instance.
| password               | string     | Password of a Memcached instance.
| name                   | string     | Memcached instance name.
| type                   | string     | The service type. The value is dcs-memcached.

## Unbind

Remove the bind information from the provisioned instance.

## Update

Update a previously provisioned instance.

### Update Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| name                         | string     | N         | Memcached instance name. An instance name is a string of 4–64 characters that contain letters, digits, underscores (_), and hyphens (-). An instance name must start with letters.
| description                  | string     | N         | Brief description of the Memcached instance. A brief description supports up to 1024 characters.
| backup_strategy_savedays     | int        | N         | Retention time. Unit: day. Range: 1–7.
| backup_strategy_backup_type  | string     | N         | Backup type. Options: auto: automatic backup. manual: manual backup.
| backup_strategy_backup_at    | []int      | N         | Days in a week on which backup starts. Range: 1–7. Where: 1 indicates Monday; 7 indicates Sunday.
| backup_strategy_begin_at     | string     | N         | Time at which backup starts. "00:00-01:00" indicates that backup starts at 00:00:00.
| backup_strategy_period_type  | string     | N         | Interval at which backup is performed. Currently, only weekly backup is supported.
| maintain_begin               | string     | N         | Time at which the maintenance time window starts.
| maintain_end                 | string     | N         | Time at which the maintenance time window ends.
| security_group_id            | string     | N         | Subnet ID.
| new_capacity                 | int        | N         | New cache capacity. Unit: GB. For a Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.
| old_password                 | string     | N         | The previous password of Memcached instance.
| new_password                 | string     | N         | The new password of Memcached instance.


## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service dcs-memcached SingleNode mymemcached -c '{
    "username": "username",
    "password": "Password1234!",
    "name": "MemcachedSingleNode",
    "description": "Memcached Single Node Test"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service mymemcached
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp mymemcached
```

Use `cf restage` command to ensure your env variable changes take effect.

```
cf restage myapp
```

Once bound, you can view the environment variables for a given application using the `cf env` command.

```
cf env myapp
```

### Unbind

To unbind a service from an application, use the `cf unbind-service` command.

```
cf unbind-service myapp mymemcached
```

### Update

To update a service, use the `cf update-service` command.

```
cf update-service mymemcached -c '{
    "name": "MemcachedSingleNode1",
    "description": "Memcached Single Node Test1",
    "new_capacity": 8,
    "old_password": "Password1234!",
    "new_password": "Password1234$"
}'
```

You can also check the status of the service instance using the `cf service` command.

```
cf service mymemcached
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service mymemcached
```
