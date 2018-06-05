# Configuration

A sample configuration can be found at [config.json](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/config.json).

## General Configuration

| Option     | Required | Type   | Description
|:-----------|:--------:|:------ |:-----------
| log_level  | Y        | String | Broker Log Level (DEBUG, INFO, ERROR, FATAL)
| username   | Y        | String | Broker Auth Username
| password   | Y        | String | Broker Auth Password

## Cloud Credentials Configuration

| Option                         | Required | Type    | Description
|:-------------------------------|:--------:|:------- |:-----------
| access_key                     | Y        | String  | Keystone Auth access key
| secret_key                     | Y        | String  | Keystone Auth secret key
| cacert_file                    | Y        | String  | Keystone Auth ca file
| cert                           | Y        | String  | Keystone Auth client cert file
| key                            | Y        | String  | Keystone Auth client key file
| domain_id                      | Y        | String  | Keystone Auth domain id
| domain_name                    | Y        | String  | Keystone Auth domain name
| endpoint_type                  | Y        | String  | Keystone endpoint type
| auth_url                       | Y        | String  | Keystone Auth URL
| insecure                       | Y        | Boolean | Keystone insecure setting
| password                       | Y        | String  | Keystone Auth password
| region                         | Y        | String  | Keystone Auth region
| tenant_id                      | Y        | String  | Keystone Auth tenant id
| tenant_name                    | Y        | String  | Keystone Auth tenant name
| token                          | Y        | String  | Keystone Auth token
| user_name                      | Y        | String  | Keystone Auth username
| UserID                         | Y        | String  | Keystone Auth userid

1 If the broker api will be deployed in the cloudfoundry, the value of ca could be the absolute path of this project. For example: ca =  ./ca.crt. The ca.crt can be deployed
in the rds-broker itself.
2 If the ca file is not needed for Keystone authentication, the value of ca can be set empty string.


## RDS Broker catalog

Please refer to the [Catalog Documentation](https://docs.cloudfoundry.org/services/api.html#catalog-mgmt) for more details about these properties.

### Catalog

| Option   | Required | Type      | Description
|:---------|:--------:|:--------- |:-----------
| services | N        | []Service | A list of [Services](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/CONFIGURATION.md#service)

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
| plans                         | N        | []ServicePlan | A list of [Plans](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/CONFIGURATION.md#service-plan) for this service
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
| rds_properties       | Y        | RDSProperties | [RDS Properties](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/CONFIGURATION.md#rds-properties)

## RDS Properties

Please refer to the [Huawei Relational Database Service Documentation](http://support.huaweicloud.com/en-us/rds_gls/index.html) for more details about these properties.

| Option                          | Required | Type      | Description
|:--------------------------------|:--------:|:--------- |:-----------
| datastore_type                  | Y        | String    | Specifies the DB engine. Currently, MySQL, PostgreSQL, and Microsoft SQL Server are supported. The value is PostgreSQL.
| datastore_version               | Y        | String    | Specifies the DB instance version.
| flavor_name                     | N        | String    | Specifies the specification ID compliant with the UUID format.
| flavor_id                       | Y        | String    | Specifies the specification name.
| volume_type                     | Y        | String    | Specifies the volume type. Valid value: It must be COMMON (SATA) or ULTRAHIGH (SSD) and is case-sensitive.
| volume_size                     | Y        | Integer   | Specifies the volume size. Its value must be a multiple of 10 and the value range is 100 GB to 2000 GB.
| region                          | Y        | String    | Specifies the region ID. Valid value: The value cannot be empty. For details about how to obtain this parameter value, see Regions and Endpoints.
| availability_zone               | Y        | String    | Specifies the ID of the AZ. Valid value: The value cannot be empty. For details about how to obtain this parameter value, see Regions and Endpoints.
| vpc_id                          | Y        | String    | Specifies the VPC ID.
| subnet_id                       | Y        | String    | Specifies the UUID for nics information.
| security_group_id               | Y        | String    | Specifies the security group ID which the RDS DB instance belongs to.
| db_port                         | Y        | String    | Specifies the database port number.
| backup_strategy_starttime       | Y        | String    | Indicates the backup start time that has been set. The backup task will be triggered within one hour after the backup start time.
| backup_strategy_keepdays        | Y        | Integer   | Specifies the number of days to retain the generated backup files. Its value range is 0 to 35.
| db_password                     | Y        | String    | Specifies the password for user root of the database. (Valid value: The value cannot be empty and should contain 8 to 32 characters, including uppercase and lowercase letters, digits, and the following special characters: ~!@#%^*-_=+?)
| db_username                     | Y        | String    | Specifies the username for user root of the database. The default value of username is root.
| db_name                         | Y        | String    | Specifies the name for user root of the database. The default value of database name is postgres.


## RDS Broker NOTE:

This rds broker will provide the database instance from Relational Database Service to the applications in Cloud Foundry.

The Services action of this rds broker API will return a service catalog about Relational Database Service with different specifications in plan.

The Provision action of this rds broker API will create a database instance in Relational Database Service.

The Deprovision action of this rds broker API will delete the created database instance in Relational Database Service.

The LastOperation action of this rds broker API will return the status of the created database instance in Relational Database Service.

The Bind action of this rds broker API will return the credentials information about the created database instance in Relational Database Service.
These credentials information can be used by the applications in CloudFoundry.

The Unbind action of this rds broker API will return unbind the credentials information about the created database instance in Relational Database Service.

The Update action of this rds broker API will update some information for Relational Database Service database instance.
