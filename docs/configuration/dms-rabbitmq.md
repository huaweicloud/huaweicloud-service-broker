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

## Plan Metadata Parameters Configuration

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| engine                 | string     | Cache engine, which is rabbitmq.
| engine_version         | string     | Cache engine version, which is 3.7.0.
| speccode               | string     | DMS specifications. Values: dms.instance.rabbitmq.single.1u2g, dms.instance.rabbitmq.single.2u4g, dms.instance.rabbitmq.single.4u8g, dms.instance.rabbitmq.cluster.2u4g.2, dms.instance.rabbitmq.cluster.2u4g.3, dms.instance.rabbitmq.cluster.2u4g.4.
| charging_type          | string     | Billing mode. Values: Yearly, Monthly and Hourly.
| vpc_id                 | string     | Indicates the ID of a VPC.
| subnet_id              | string     | Indicates the ID of a subnet.
| security_group_id      | string     | Indicates the ID of a security group.
| availability_zones     | []string   | Indicates the ID of an AZ. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint).

These plans are only differently configured by the Parameter Name [```speccode```].
