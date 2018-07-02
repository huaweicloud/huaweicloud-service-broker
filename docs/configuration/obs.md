# Object Storage Service

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| obs                            | Object Storage Service

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| Standard                       | OBS Standard features low access latency and high throughput
| Infrequent Access              | OBS Infrequent Access is applicable to storing semi-frequently accessed data requiring quick response
| Archive                        | OBS Archive is applicable to archiving rarely-accessed data

## Plan Metadata Parameters Configuration

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| storage_class          | string     | OBS storage classes (Standard, Infrequent Access, and Archive) are designed to meet customers' varying requirements on storage performance and costs. Standard: features low access latency and high throughput. Infrequent Access: applicable to storing semi-frequently accessed (less than 12 times a year) data requiring quick response. Archive: applicable to archiving rarely-accessed (once a year) data.
| bucket_policy          | string     | A bucket policy defines the access control policy of resources (buckets and objects) on OBS. private: Only the bucket owner can read, write, and delete objects in the bucket. This policy is the default bucket policy. public-read: Any user can read objects in the bucket. Only the bucket owner can write and delete objects in the bucket. public-read-write: Any user can read, write, and delete objects in the bucket.

These plans are only differently configured by the Parameter Name [```storage_class```].
