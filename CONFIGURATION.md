# Configuration

A sample configuration can be found at [config-sample.json](https://github.com/cloudfoundry-community/pe-rds-broker/blob/master/config-sample.json).

## General Configuration

| Option     | Required | Type   | Description
|:-----------|:--------:|:------ |:-----------
| log_level  | Y        | String | Broker Log Level (DEBUG, INFO, ERROR, FATAL)
| username   | Y        | String | Broker Auth Username
| password   | Y        | String | Broker Auth Password
| rds_config | Y        | Hash   | [RDS Broker configuration](https://github.com/cloudfoundry-community/pe-rds-broker/blob/master/CONFIGURATION.md#rds-broker-configuration)

## RDS Broker Configuration

| Option                         | Required | Type    | Description
|:-------------------------------|:--------:|:------- |:-----------
| region                         | Y        | String  | RDS Region
| db_prefix                      | Y        | String  | Prefix to add to RDS DB Identifiers
| allow_user_provision_parameters| N        | Boolean | Allow users to send arbitrary parameters on provision calls (defaults to `false`)
| allow_user_update_parameters   | N        | Boolean | Allow users to send arbitrary parameters on update calls (defaults to `false`)
| allow_user_bind_parameters     | N        | Boolean | Allow users to send arbitrary parameters on bind calls (defaults to `false`)
| catalog                        | Y        | Hash    | [RDS Broker catalog](https://github.com/cloudfoundry-community/pe-rds-broker/blob/master/CONFIGURATION.md#rds-broker-catalog)

## RDS Broker catalog

Please refer to the [Catalog Documentation](https://docs.cloudfoundry.org/services/api.html#catalog-mgmt) for more details about these properties.

### Catalog

| Option   | Required | Type      | Description
|:---------|:--------:|:--------- |:-----------
| services | N        | []Service | A list of [Services](https://github.com/cloudfoundry-community/pe-rds-broker/blob/master/CONFIGURATION.md#service)

### Service

| Option                        | Required | Type          | Description
|:------------------------------|:--------:|:------------- |:-----------
| id                            | Y        | String        | An identifier used to correlate this service in future requests to the catalog
| name                          | Y        | String        | The CLI-friendly name of the service that will appear in the catalog. All lowercase, no spaces
| description                   | Y        | String        | A short description of the service that will appear in the catalog
| bindable                      | N        | Boolean       | Whether the service can be bound to applications
| tags                          | N        | []String      | A list of service tags
| metadata.displayName          | N        | String        | The name of the service to be displayed in graphical clients
| metadata.imageUrl             | N        | String        | The URL to an image
| metadata.longDescription      | N        | String        | Long description
| metadata.providerDisplayName  | N        | String        | The name of the upstream entity providing the actual service
| metadata.documentationUrl     | N        | String        | Link to documentation page for service
| metadata.supportUrl           | N        | String        | Link to support for the service
| requires                      | N        | []String      | A list of permissions that the user would have to give the service, if they provision it (only `syslog_drain` is supported)
| plan_updateable               | N        | Boolean       | Whether the service supports upgrade/downgrade for some plans
| plans                         | N        | []ServicePlan | A list of [Plans](https://github.com/cloudfoundry-community/pe-rds-broker/blob/master/CONFIGURATION.md#service-plan) for this service
| dashboard_client.id           | N        | String        | The id of the Oauth2 client that the service intends to use
| dashboard_client.secret       | N        | String        | A secret for the dashboard client
| dashboard_client.redirect_uri | N        | String        | A domain for the service dashboard that will be whitelisted by the UAA to enable SSO

### Service Plan

| Option               | Required | Type          | Description
|:---------------------|:--------:|:------------- |:-----------
| id                   | Y        | String        | An identifier used to correlate this plan in future requests to the catalog
| name                 | Y        | String        | The CLI-friendly name of the plan that will appear in the catalog. All lowercase, no spaces
| description          | Y        | String        | A short description of the plan that will appear in the catalog
| metadata.bullets     | N        | []String      | Features of this plan, to be displayed in a bulleted-list
| metadata.costs       | N        | Cost Object   | An array-of-objects that describes the costs of a service, in what currency, and the unit of measure
| metadata.displayName | N        | String        | Name of the plan to be display in graphical clients
| free                 | N        | Boolean       | This field allows the plan to be limited by the non_basic_services_allowed field in a Cloud Foundry Quota
| rds_properties       | Y        | RDSProperties | [RDS Properties](https://github.com/cloudfoundry-community/pe-rds-broker/blob/master/CONFIGURATION.md#rds-properties)

## RDS Properties

Please refer to the [Amazon Relational Database Service Documentation](https://aws.amazon.com/documentation/rds/) for more details about these properties.

| Option                          | Required | Type      | Description
|:--------------------------------|:--------:|:--------- |:-----------
| allocated_storage               | Y        | Integer   | The amount of storage (in gigabytes) to be initially allocated for the database instances (between `5` and `6144`). Not applicable when using `aurora`
| auto_minor_version_upgrade      | N        | Boolean   | Enable or disable automatic upgrades to new minor versions as they are released (defaults to `false`)
| availability_zone               | N        | String    | The Availability Zone that database instances will be created in
| backup_retention_period         | N        | Integer   | The number of days that Amazon RDS should retain automatic backups of DB instances (between `0` and `35`)
| character_set_name              | N        | String    | For supported engines, indicates that DB instances should be associated with the specified CharacterSet. Not applicable when using `aurora`
| copy_tags_to_snapshot           | N        | Boolean   | Enable or disable copying all tags from DB instances to snapshots
| db_instance_class               | Y        | String    | The name of the DB Instance Class
| db_parameter_group_name         | N        | String    | The DB parameter group name that defines the configuration settings you want applied to DB instances
| db_cluster_parameter_group_name | N        | String    | The DB cluster parameter group name that defines the configuration settings you want applied to DB clusters (only for `aurora`)
| db_security_groups              | N        | []String  | The security group(s) names that have rules authorizing connections from applications that need to access the data stored in the DB instance. Not applicable when using `aurora`
| db_subnet_group_name            | N        | String    | The DB subnet group name that defines which subnets and IP ranges the DB instance can use in the VPC
| engine                          | Y        | String    | The name of the Database Engine (only `aurora`, `mariadb`, `mysql` and `postgres` are supported)
| engine_version                  | Y        | String    | The version number of the Database Engine
| iops                            | N        | Integer   | The amount of Provisioned IOPS to be initially allocated for DB instances when using `io1` storage type. Not applicable when using `aurora`
| kms_key_id                      | N        | String    | The KMS key identifier for encrypted DB instances. Not applicable when using `aurora`
| license_model                   | N        | String    | License model information for DB instances (`license-included`, `bring-your-own-license`, `general-public-license`). Not applicable when using `aurora`
| multi_az                        | N        | Boolean   | Enable or disable Multi-AZ deployment for high availability DB Instances. Not applicable when using `aurora`
| option_group_name               | N        | String    | The DB option group name that enables any optional functionality you want the DB instances to support. Not applicable when using `aurora`
| port                            | N        | Integer   | The TCP/IP port DB instances will use for application connections
| preferred_backup_window         | N        | String    | The daily time range during which automated backups are created if automated backups are enabled
| preferred_maintenance_window    | N        | String    | The weekly time range during which system maintenance can occur
| publicly_accessible             | N        | Boolean   | Specify if DB instances will be publicly accessible
| skip_final_snapshot             | N        | Boolean   | Determines whether a final DB snapshot is created before the DB instances are deleted
| storage_encrypted               | N        | Boolean   | Specifies whether DB instances are encrypted. Not applicable when using `aurora`
| storage_type                    | N        | String    | The storage type to be associated with DB instances (`standard`, `gp2`, `io1`)
| vpc_security_group_ids          | N        | []String  | VPC security group(s) IDs that have rules authorizing connections from applications that need to access the data stored in DB instances
