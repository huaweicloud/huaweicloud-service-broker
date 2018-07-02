# Distributed Message Service for Standard

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dms-standard                   | Distributed Message Service for Standard

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| Partition-level FIFO           | Standard queue. Messages might be retrieved in an order different from which they were sent
| Global FIFO                    | FIFO delivery. Messages are retrieved in the order they were sent

## Plan Metadata Parameters Configuration

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| queue_mode             | string     | Indicates the queue type. Options: NORMAL: Standard queue. Best-effort ordering. Messages might be retrieved in an order different from which they were sent. Select standard queues when throughput is important. FIFO: First-ln-First-out (FIFO) queue. FIFO delivery. Messages are retrieved in the order they were sent. Select FIFO queues when the order of messages is important. Default value: NORMAL.
| redrive_policy         | string     | This parameter is mandatory only when queue_mode is NORMAL or FIFO. Indicates whether to enable dead letter messages. Dead letter messages indicate messages that cannot be normally consumed. If a message fails to be consumed after the number of consumption attempts of this message reaches the maximum value, DMS stores this message into the dead letter queue. This message will be retained in the deal letter queue for 72 hours. During this period, consumers can consume the dead letter message. Dead letter messages can be consumed only by the consumer group that generates these dead letter messages. Dead letter messages of a FIFO queue are stored and consumed based on the FIFO sequence. Options: enable disable Default value: disable.
| max_consume_count      | int        | This parameter is mandatory only when redrive_policy is set to enable. This parameter indicates the maximum number of allowed message consumption failures. When a message fails to be consumed after the number of consumption attempts of this message reaches this value, DMS stores this message into the dead letter queue. Value range: 1-100.

These plans are only differently configured by the Parameter Name [```queue_mode```].
