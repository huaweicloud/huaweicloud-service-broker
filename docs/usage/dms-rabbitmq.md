# Distributed Message Service for RabbitMQ

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-rabbitmq                   | Distributed Message Service for RabbitMQ

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| SingleNode                     | RabbitMQ Single Node
| Cluster                        | RabbitMQ Cluster

## Provision

Provision a new instance for Distributed Message Service for RabbitMQ.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| vpc_id                       | string     | N         | Indicates the ID of a VPC. The default value is in the config file.
| subnet_id                    | string     | N         | Indicates the ID of a subnet. The default value is in the config file.
| security_group_id            | string     | N         | Indicates the ID of a security group. The default value is in the config file.
| availability_zones           | []string   | N         | Indicates the ID of an AZ. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint). The default value is in the config file.
| username                     | string     | Y         | Indicates a username. A username consists of 1 to 64 characters and supports only letters, digits, and hyphens (-).
| password                     | string     | Y         | Indicates the password of an instance. An instance password must meet the following complexity requirements: Must be 6 to 32 characters long. Must contain at least two of the following character types: Lowercase letters Uppercase letters Digits Special characters (`~!@#$%^&*()-_=+\|[{}]:'",<.>/?)
| name                         | string     | Y         | Indicates the name of an instance. An instance name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).
| description                  | string     | N         | Indicates the description of an instance. It is a character string containing not more than 1024 characters.
| maintain_begin               | string     | N         | Indicates the time at which a maintenance time window starts. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The start time must be set to 22:00:00, 02:00:00, 06:00:00, 10:00:00, 14:00:00, or 18:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_begin is left blank, parameter maintain_end is also blank. In this case, the system automatically allocates the default start time 02:00:00.
| maintain_end                 | string     | N         | Indicates the time at which a maintenance time window ends. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The end time is four hours later than the start time. For example, if the start time is 22:00:00, the end time is 02:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_end is left blank, parameter maintain_begin is also blank. In this case, the system automatically allocates the default end time 06:00:00.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| host                   | string     | The fully-qualified address of RabbitMQ instance.
| port                   | int        | The port number to connect to RabbitMQ instance.
| username               | string     | Username of a RabbitMQ instance.
| password               | string     | Password of a RabbitMQ instance.
| uri                    | string     | The uri to connect to RabbitMQ instance.
| type                   | string     | The service type. The value is dms-rabbitmq.

## Unbind

Remove the bind information from the provisioned instance.

## Update

Update a previously provisioned instance.

### Update Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| name                         | string     | N         | Indicates the name of an instance. An instance name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).
| description                  | string     | N         | Indicates the description of an instance. It is a character string containing not more than 1024 characters.
| maintain_begin               | string     | N         | Indicates the time at which a maintenance time window starts. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The start time must be set to 22:00:00, 02:00:00, 06:00:00, 10:00:00, 14:00:00, or 18:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_begin is left blank, parameter maintain_end is also blank. In this case, the system automatically allocates the default start time 02:00:00.
| maintain_end                 | string     | N         | Indicates the time at which a maintenance time window ends. Format: HH:mm:ss. The start time and end time of a maintenance time window must indicate the time segment of a supported maintenance time window. For details, see section Querying Maintenance Time Windows. The end time is four hours later than the start time. For example, if the start time is 22:00:00, the end time is 02:00:00. Parameters maintain_begin and maintain_end must be set in pairs. If parameter maintain_end is left blank, parameter maintain_begin is also blank. In this case, the system automatically allocates the default end time 06:00:00.
| security_group_id            | string     | N         | Indicates the ID of a security group.

## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service dms-rabbitmq SingleNode myrabbitmq -c '{
    "username": "username",
    "password": "Password1234!",
    "name": "RabbitMQSingleNode",
    "description": "RabbitMQ Single Node Test"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service myrabbitmq
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp myrabbitmq
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
cf unbind-service myapp myrabbitmq
```

### Update

To update a service, use the `cf update-service` command.

```
cf update-service myrabbitmq -c '{
    "name": "RabbitMQSingleNode1",
    "description": "RabbitMQ Single Node Test1"
}'
```

You can also check the status of the service instance using the `cf service` command.

```
cf service myrabbitmq
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service myrabbitmq
```
