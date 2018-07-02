# Distributed Message Service for Kafka

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-kafka                      | Distributed Message Service for Kafka

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| High throughput                | High-throughput Kafka queue. All message replicas are flushed to a disk asynchronously
| High reliability               | High-availability Kafka queue. All message replicas are flushed to a disk synchronously

## Plan Metadata Parameters Configuration

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| queue_mode             | string     | Indicates the queue type. Options: KAFKA_HA: High-availability Kafka queue. All message replicas are flushed to a disk synchronously. Select the high availability mode when message reliability is important. KAFKA_HT: High-throughput Kafka queue. All message replicas are flushed to a disk asynchronously. Select the high throughput mode when message delivery performance is important.
| endpoint_name          | string     | Kafka endpoint name. The default value: dms-kafka. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint).
| endpoint_port          | string     | Kafka endpoint port. The default value: 37000. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint).
| retention_hours        | int        | This parameter is mandatory only when queue_mode is set to KAFKA_HA or KAFKA_HT. This parameter indicates the retention time of messages in Kafka queues. Value range: 1 to 72 hours.

These plans are only differently configured by the Parameter Name [```queue_mode```].
