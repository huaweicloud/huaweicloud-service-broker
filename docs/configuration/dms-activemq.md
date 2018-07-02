# Distributed Message Service for ActiveMQ

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-activemq                   | Distributed Message Service for ActiveMQ

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| ActiveMQ                       | Advanced Message Queuing Protocol (AMQP) queue. AMQP is an open standard application layer protocol for message-oriented middleware

## Plan Metadata Parameters Configuration

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| queue_mode             | string     | Indicates the queue type. Options: AMQP: Advanced Message Queuing Protocol (AMQP) queue. AMQP is an open standard application layer protocol for message-oriented middleware.
| endpoint_name          | string     | AMQP endpoint name. The default value: dms-amqp. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint).
| endpoint_port          | string     | AMQP endpoint port. The default value: 60020. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint).

These plans are only differently configured by the Parameter Name [```queue_mode```].
