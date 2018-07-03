# Distributed Message Service for ActiveMQ

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-activemq                   | Distributed Message Service for ActiveMQ

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| ActiveMQ                       | Advanced Message Queuing Protocol (AMQP) queue. AMQP is an open standard application layer protocol for message-oriented middleware

## Provision

Provision a new instance for Distributed Message Service for ActiveMQ.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| queue_name                   | string     | Y         | Indicates the name of a queue. A queue name starts with a letter, consists of 4 to 64 characters, and supports only letters, digits, and hyphens (-).
| group_name                   | string     | Y         | Indicates the name of a consumer group. A string of 1 to 32 characters that contain a-z, A-Z, 0-9, hyphens (-), and underscores (_).
| description                  | string     | N         | Indicates the description of an instance. It is a character string containing not more than 1024 characters.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| region                 | string     | The region name of a ActiveMQ instance.
| projectid              | string     | The projectid of a ActiveMQ instance.
| protocoltype           | string     | The protocal type of a ActiveMQ instance. It may HTTPS or TCP.
| url                    | string     | The connection address of a ActiveMQ instance.
| ak                     | string     | The Access Key of a ActiveMQ instance.
| sk                     | string     | The Secret Key of a ActiveMQ instance.
| queueid                | string     | The queue ID of a ActiveMQ instance.
| groupid                | string     | The consumer group ID of a ActiveMQ instance.
| type                   | string     | The service type. The value is dms-activemq.

## Unbind

Remove the bind information from the provisioned instance.

## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service dms-activemq ActiveMQ myactivemq -c '{
    "queue_name": "activemqqueue",
    "group_name": "activemqgroup"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service myactivemq
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp myactivemq
```

Once bound, you can view the environment variables for a given application using the `cf env` command.

```
cf env myapp
```

### Unbind

To unbind a service from an application, use the `cf unbind-service` command.

```
cf unbind-service myapp myactivemq
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service myactivemq
```
