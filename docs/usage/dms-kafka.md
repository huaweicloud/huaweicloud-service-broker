# Distributed Message Service for Kafka

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-kafka                      | Distributed Message Service for Kafka

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| HighThroughput                 | High-throughput Kafka queue. All message replicas are flushed to a disk asynchronously
| HighReliability                | High-availability Kafka queue. All message replicas are flushed to a disk synchronously

## Provision

Provision a new instance for Distributed Message Service for Kafka.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| retention_hours              | int        | N         | This parameter is mandatory only when queue_mode is set to KAFKA_HA or KAFKA_HT. This parameter indicates the retention time of messages in Kafka queues. Value range: 1 to 72 hours. The default value is in the config file.
| queue_name                   | string     | Y         | Indicates the name of a queue. A queue name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).
| group_name                   | string     | Y         | Indicates the name of a consumer group. A string of 1 to 32 characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_).
| description                  | string     | N         | Indicates the description of an instance. It is a character string containing not more than 1024 characters.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| region                 | string     | The region name of a Kafka instance.
| projectid              | string     | The projectid of a Kafka instance.
| protocoltype           | string     | The protocal type of a Kafka instance. It may HTTPS or TCP.
| url                    | string     | The connection address of a Kafka instance.
| ak                     | string     | The Access Key of a Kafka instance.
| sk                     | string     | The Secret Key of a Kafka instance.
| queueid                | string     | The queue ID of a Kafka instance.
| groupid                | string     | The consumer group ID of a Kafka instance.
| type                   | string     | The service type. The value is dms-kafka.

## Unbind

Remove the bind information from the provisioned instance.

## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service dms-kafka HighThroughput mykafka -c '{
    "queue_name": "kafkaqueue",
    "group_name": "kafkagroup"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service mykafka
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp mykafka
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
cf unbind-service myapp mykafka
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service mykafka
```
