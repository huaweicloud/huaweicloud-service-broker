# RDS PostgreSQL Service

| Service Name                   | Description
|:-------------------------------|:-----------
| rds-postgresql                 | RDS PostgreSQL Service

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| 9.5.5                          | RDS PostgreSQL 9.5.5
| 9.6.5                          | RDS PostgreSQL 9.6.5
| 10.0.3                         | RDS PostgreSQL 10.0.3

## Provision

Provision a new instance for RDS PostgreSQL Service.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| speccode                     | string     | N         | Indicates the resource specifications code. Use rds.pg.s1.xlarge as an example. rds indicates RDS, pg indicates the DB engine, and s1.xlarge indicates the performance specification. The parameter containing rr indicates the read replica specifications. The parameter not containing rr indicates the single or primary/standby DB instance specifications. If you enable HA, the suffix ```.ha``` need be added to the DB instance name. For example, the DB instance name is ```rds.db.s1.xlarge.ha```. The default value is in the config file.
| volume_type                  | string     | N         | Specifies the volume type. Valid value: It must be COMMON (SATA) or ULTRAHIGH (SSD) and is case-sensitive. The default value is in the config file.
| volume_size                  | int        | N         | Specifies the volume size. Its value must be a multiple of 10 and the value range is 100 GB to 2000 GB. The default value is in the config file.
| availability_zone            | string     | N         | Specifies the ID of the AZ. Valid value: The value cannot be empty. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint). The default value is in the config file.
| vpc_id                       | string     | N         | Specifies the VPC ID. The default value is in the config file.
| subnet_id                    | string     | N         | Specifies the UUID for nics information. The default value is in the config file.
| security_group_id            | string     | N         | Specifies the security group ID which the RDS DB instance belongs to. The default value is in the config file.
| name                         | string     | Y         | Specifies the DB instance name. The DB instance name of the same DB engine is unique for the same tenant. Valid value: The value must be 4 to 64 characters in length and start with a letter. It is case-insensitive and can contain only letters, digits, hyphens (-), and underscores (_).
| database_port                | string     | N         | Specifies the database port number.
| database_password            | string     | Y         | Specifies the password for user root of the database. Valid value: The value cannot be empty and should contain 8 to 32 characters, including uppercase and lowercase letters, digits, and the following special characters: ~!@#%^*-_=+?
| backup_strategy_starttime    | string     | N         | Indicates the backup start time that has been set. The backup task will be triggered within one hour after the backup start time.
| backup_strategy_keepdays     | int        | N         | Specifies the number of days to retain the generated backup files. Its value range is 0 to 35.
| ha_enable                    | bool       | N         | Specifies the HA configuration parameter. Valid value: The value is true or false. The value true indicates creating HA DB instances. The value false indicates creating a single DB instance.
| ha_replicationmode           | string     | N         | Specifies the replication mode for the standby DB instance. The value cannot be empty. For PostgreSQL, the value is async or sync.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| host                   | string     | The fully-qualified address of PostgreSQL instance.
| port                   | int        | The port number to connect to PostgreSQL instance.
| name                   | string     | PostgreSQL instance name.
| username               | string     | Username of a PostgreSQL instance.
| password               | string     | Password of a PostgreSQL instance.
| uri                    | string     | The uri to connect to PostgreSQL instance.
| type                   | string     | The service type. The value is rds-postgresql.

## Unbind

Remove the bind information from the provisioned instance.

## Update

Update a previously provisioned instance.

### Update Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| volume_size                  | int        | N         | Specifies the volume size. Its value must be a multiple of 10 and the value range is 100 GB to 2000 GB.
| speccode                     | string     | N         | Indicates the resource specifications code. Use rds.pg.s1.xlarge as an example. rds indicates RDS, pg indicates the DB engine, and s1.xlarge indicates the performance specification. The parameter containing rr indicates the read replica specifications. The parameter not containing rr indicates the single or primary/standby DB instance specifications.

## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service rds-postgresql 9.5.5 myrdspostgresql -c '{
    "name": "RDSPostgreSQL",
    "database_password": "Password1234!"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service myrdspostgresql
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp myrdspostgresql
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
cf unbind-service myapp myrdspostgresql
```

### Update

To update a service, use the `cf update-service` command.

```
cf update-service myrdspostgresql -c '{
    "volume_size":400
}'
```

You can also check the status of the service instance using the `cf service` command.

```
cf service myrdspostgresql
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service myrdspostgresql
```
