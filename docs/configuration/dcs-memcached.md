# Distributed Cache Service for Memcached

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| dcs-memcached                  | Distributed Cache Service for Memcached

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| SingleNode                     | Memcached Single Node
| MasterStandby                  | Memcached Master Standby

## Plan Metadata Parameters Configuration

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| engine                 | string     | Cache engine, which is Memcached.
| speccode               | string     | DCS specifications. Values: dcs.memcached.single_node and dcs.memcached.master_standby.
| charging_type          | string     | Billing mode. Values: Yearly, Monthly and Hourly.
| capacity               | int        | Cache capacity. Unit: GB. For a DCS Memcached instance in single node or master standby mode, the cache capacity can be 2 GB, 4 GB, 8 GB, 16 GB, 32 GB, or 64 GB.
| vpc_id                 | string     | Tenant's VPC ID.
| subnet_id              | string     | Subnet ID.
| security_group_id      | string     | Tenant's security group ID.
| availability_zones     | []string   | IDs of the AZs where cache nodes reside. For details about how to obtain this parameter value, see [Regions and Endpoints](https://developer.huaweicloud.com/endpoint).


These plans are only differently configured by the Parameter Name [```speccode```].
