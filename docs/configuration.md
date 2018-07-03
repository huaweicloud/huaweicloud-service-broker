# Configuration

A sample configuration can be found at [config.json](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/config.json).

## General Configuration (broker_config)

| Option     | Required | Type   | Description
|:-----------|:--------:|:------ |:-----------
| log_level  | Y        | string | Broker Log Level (DEBUG, INFO, ERROR, FATAL)
| username   | Y        | string | Broker Auth Username
| password   | Y        | string | Broker Auth Password

## Back Database Configuration (back_database)

| Option              | Required | Type   | Description
|:--------------------|:--------:|:------ |:-----------
| database_type       | Y        | string | Back database type. Currently only mysql is supported
| database_host       | Y        | string | Back database host
| database_port       | Y        | string | Back database host port
| database_name       | Y        | string | Back database instance name
| database_username   | Y        | string | Back database instance user name
| database_password   | Y        | string | Back database instance user password

## Cloud Credentials Configuration (cloud_credentials)

| Option                         | Required | Type    | Description
|:-------------------------------|:--------:|:------- |:-----------
| access_key                     | Y        | string  | IAM Auth access key
| secret_key                     | Y        | string  | IAM Auth secret key
| cacert_file                    | N        | string  | IAM Auth ca file
| cert                           | N        | string  | IAM Auth client cert file
| key                            | N        | string  | IAM Auth client key file
| domain_id                      | N        | string  | IAM Auth domain id
| domain_name                    | Y        | string  | IAM Auth domain name
| endpoint_type                  | N        | string  | IAM endpoint type
| auth_url                       | Y        | string  | IAM Auth URL
| insecure                       | N        | bool    | IAM insecure setting
| password                       | Y        | string  | IAM Auth password
| region                         | Y        | string  | IAM Auth region
| tenant_id                      | N        | string  | IAM Auth tenant id
| tenant_name                    | Y        | string  | IAM Auth tenant name
| token                          | N        | string  | IAM Auth token
| user_name                      | Y        | string  | IAM Auth username
| user_id                        | N        | string  | IAM Auth userid

### Catalog (catalog)

| Option   | Required | Type      | Description
|:---------|:--------:|:--------- |:-----------
| services | N        | []Service | A list of Services in Service Broker

### Service (services)

| Option                        | Required | Type          | Description
|:------------------------------|:--------:|:------------- |:-----------
| id                            | Y        | string        | An identifier used to correlate this service in future requests to the catalog
| name                          | Y        | string        | The CLI-friendly name of the service that will appear in the catalog. All lowercase, no spaces
| description                   | Y        | string        | A short description of the service that will appear in the catalog
| bindable                      | N        | bool          | Whether the service can be bound to applications
| tags                          | N        | []string      | A list of service tags
| metadata.displayName          | N        | string        | The name of the service to be displayed in graphical clients
| metadata.imageUrl             | N        | string        | The URL to an image
| metadata.longDescription      | N        | string        | Long description
| metadata.providerDisplayName  | N        | string        | The name of the upstream entity providing the actual service
| metadata.documentationUrl     | N        | string        | Link to documentation page for service
| metadata.supportUrl           | N        | string        | Link to support for the service
| requires                      | N        | []string      | A list of permissions that the user would have to give the service, if they provision it (only `syslog_drain` is supported)
| plan_updateable               | N        | bool          | Whether the service supports upgrade/downgrade for some plans
| plans                         | N        | []ServicePlan | A list of Service Plans in Service Broker
| dashboard_client.id           | N        | string        | The id of the Oauth2 client that the service intends to use
| dashboard_client.secret       | N        | string        | A secret for the dashboard client
| dashboard_client.redirect_uri | N        | string        | A domain for the service dashboard that will be whitelisted by the UAA to enable SSO

### Service Plan (plans)

| Option               | Required | Type                   | Description
|:---------------------|:--------:|:---------------------- |:-----------
| id                   | Y        | string                 | An identifier used to correlate this plan in future requests to the catalog
| name                 | Y        | string                 | The CLI-friendly name of the plan that will appear in the catalog. All lowercase, no spaces
| description          | Y        | string                 | A short description of the plan that will appear in the catalog
| metadata.bullets     | N        | []string               | Features of this plan, to be displayed in a bulleted-list
| metadata.costs       | N        | Cost Object            | An array-of-objects that describes the costs of a service, in what currency, and the unit of measure
| metadata.displayName | N        | string                 | Name of the plan to be display in graphical clients
| metadata.parameters  | N        | map[string]inferface{} | Parameters of the plan to be set for each Service Plan
| free                 | N        | bool                   | This field allows the plan to be limited by the non_basic_services_allowed field

metadata.parameters are key-value fields for each Service Plan.
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
