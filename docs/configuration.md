# Configuration

A sample configuration can be found at [config.json](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/config.json).

## General Configuration (broker_config)

| Option     | Required | Type   | Description
|:-----------|:--------:|:------ |:-----------
| log_level  | Y        | String | Broker Log Level (DEBUG, INFO, ERROR, FATAL)
| username   | Y        | String | Broker Auth Username
| password   | Y        | String | Broker Auth Password

## Back Database Configuration (back_database)

| Option              | Required | Type   | Description
|:--------------------|:--------:|:------ |:-----------
| database_type       | Y        | String | Back database type. Currently only mysql is supported
| database_host       | Y        | String | Back database host
| database_port       | Y        | String | Back database host port
| database_name       | Y        | String | Back database instance name
| database_username   | Y        | String | Back database instance user name
| database_password   | Y        | String | Back database instance user password

## Cloud Credentials Configuration (cloud_credentials)

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

If the ca file is not needed for authentication, the value of ca can be set empty string.

### Catalog

| Option   | Required | Type      | Description
|:---------|:--------:|:--------- |:-----------
| services | N        | []Service | A list of Services in Service Broker

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
| plans                         | N        | []ServicePlan | A list of Service Plans in Service Broker
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
| metadata.parameters  | N        | String        | Parameters of the plan to be set for each Service Plan
| free                 | N        | Boolean       | This field allows the plan to be limited by the non_basic_services_allowed field

metadata.parameters are key value field for each Service Plan.
The Application Developers can find how to configure the each service plan by the following instructions.

* Distributed Cache Service for Redis: [dcs-redis](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dcs-redis.md)
* Distributed Cache Service for Memcached: [dcs-memcached](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dcs-memcached.md)
* Distributed Cache Service for IMDG: [dcs-imdg](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dcs-imdg.md)
* Distributed Message Service for Standard: [dms-standard](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dms-standard.md)
* Distributed Message Service for ActiveMQ: [dms-activemq](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dms-activemq.md)
* Distributed Message Service for Kafka: [dms-kafka](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dms-kafka.md)
* Distributed Message Service for RabbitMQ: [dms-rabbitmq](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/dms-rabbitmq.md)
* Object Storage Service: [obs](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/obs.md)
* RDS MySQL Service: [rds-mysql](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/rds-mysql.md)
* RDS SQLServer Service: [rds-sqlserver](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/rds-sqlserver.md)
* RDS PostgreSQL Service: [rds-postgresql](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/rds-postgresql.md)
* RDS MySQL Service: [rds-hwsql](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration/rds-hwsql.md)
