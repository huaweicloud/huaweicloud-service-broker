# Distributed Message Service for Standard

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-standard                   | Distributed Message Service for Standard

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| PartitionLevelFIFO             | Standard queue. Messages might be retrieved in an order different from which they were sent
| GlobalFIFO                     | FIFO delivery. Messages are retrieved in the order they were sent

## Provision

Provision a new instance for Distributed Message Service for Standard.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| redrive_policy               | string     | N         | This parameter is mandatory only when queue_mode is NORMAL or FIFO. Indicates whether to enable dead letter messages. Dead letter messages indicate messages that cannot be normally consumed. If a message fails to be consumed after the number of consumption attempts of this message reaches the maximum value, DMS stores this message into the dead letter queue. This message will be retained in the deal letter queue for 72 hours. During this period, consumers can consume the dead letter message. Dead letter messages can be consumed only by the consumer group that generates these dead letter messages. Dead letter messages of a FIFO queue are stored and consumed based on the FIFO sequence. Options: enable disable. The default value is in the config file.
| max_consume_count            | int        | N         | This parameter is mandatory only when redrive_policy is set to enable. This parameter indicates the maximum number of allowed message consumption failures. When a message fails to be consumed after the number of consumption attempts of this message reaches this value, DMS stores this message into the dead letter queue. Value range: 1-100. The default value is in the config file.
| queue_name                   | string     | Y         | Indicates the name of a queue. A queue name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).
| group_name                   | string     | Y         | Indicates the name of a consumer group. A string of 1 to 32 characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_).
| description                  | string     | N         | Indicates the description of an instance. It is a character string containing not more than 1024 characters.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| region                 | string     | The region name of a Standard instance.
| projectid              | string     | The projectid of a Standard instance.
| protocoltype           | string     | The protocal type of a Standard instance. It may HTTPS or TCP.
| url                    | string     | The connection address of a Standard instance.
| ak                     | string     | The Access Key of a Standard instance.
| sk                     | string     | The Secret Key of a Standard instance.
| queueid                | string     | The queue ID of a Standard instance.
| groupid                | string     | The consumer group ID of a Standard instance.
| type                   | string     | The service type. The value is dms-standard.

## Unbind

Remove the bind information from the provisioned instance.

## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service dms-standard PartitionLevelFIFO mystandard -c '{
    "queue_name": "standardqueue",
    "group_name": "standardgroup"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service mystandard
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp mystandard
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
cf unbind-service myapp mystandard
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service mystandard
```
