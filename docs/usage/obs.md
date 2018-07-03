# Object Storage Service

## Services

| Service Name                   | Description
|:-------------------------------|:-----------
| obs                            | Object Storage Service

## Plans

| Plan Name                      | Description
|:-------------------------------|:-----------
| Standard                       | OBS Standard features low access latency and high throughput
| InfrequentAccess               | OBS Infrequent Access is applicable to storing semi-frequently accessed data requiring quick response
| Archive                        | OBS Archive is applicable to archiving rarely-accessed data

## Provision

Provision a new instance for Object Storage Service.

### Provision Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| bucket_name                  | string     | Y         | Enter the bucket name, which must be globally unique. Name the bucket according to the globally applied DNS naming regulation as follows: Must contain 3 to 63 characters, including lowercase letters, digits, hyphens (-), and periods (.) Cannot be an IP address. Cannot start or end with a hyphen (-) or period (.) Cannot contain two consecutive periods (.), for example, my..bucket. Cannot contain periods (.) and hyphens (-) adjacent to each other, for example, my-.bucket or my.-bucket.
| bucket_policy                | string     | N         | A bucket policy defines the access control policy of resources (buckets and objects) on OBS. private: Only the bucket owner can read, write, and delete objects in the bucket. This policy is the default bucket policy. public-read: Any user can read objects in the bucket. Only the bucket owner can write and delete objects in the bucket. public-read-write: Any user can read, write, and delete objects in the bucket. The default value is in the config file.

## Bind

Create a new credentials on the provisioned instance.
Bind returns the following connection details and credentials.

### Bind Credentials

| Parameter Name         | Type       | Description
|:-----------------------|:-----------|:-----------
| region                 | string     | The region name of a Object Storage Service.
| url                    | string     | The url to connect to Object Storage Service instance.
| bucketname             | string     | Name of Object Storage Service instance.
| ak                     | string     | The Access Key of Object Storage Service instance.
| sk                     | string     | The Secret Key of Object Storage Service instance.
| type                   | string     | The service type. The value is obs.

## Unbind

Remove the bind information from the provisioned instance.

## Update

Update a previously provisioned instance.

### Update Parameters

| Parameter Name               | Type       | Required  | Description
|:-----------------------------|:-----------|:----------|:-----------
| bucket_policy                | string     | N         | A bucket policy defines the access control policy of resources (buckets and objects) on OBS. private: Only the bucket owner can read, write, and delete objects in the bucket. This policy is the default bucket policy. public-read: Any user can read objects in the bucket. Only the bucket owner can write and delete objects in the bucket. public-read-write: Any user can read, write, and delete objects in the bucket.
| status                       | string     | N         | By default, the versioning function is disabled for new buckets on OBS. The status include: Enabled and Suspended.

## Deprovision

Delete the provisioned instance.

## Example on Cloud Foundry

### Provision

The following command will create a service.

```
cf create-service obs Standard myobs -c '{
    "bucket_name": "obsstandard",
    "bucket_policy": "public-read-write"
}'
```

You can check the status of the service instance using the `cf service` command.

```
cf service myobs
```

### Bind

Once the service has been successfully provisioned, you can bind to it by using
`cf bind-service` or by including it in a Cloud Foundry manifest.

```
cf bind-service myapp myobs
```

Once bound, you can view the environment variables for a given application using the `cf env` command.

```
cf env myapp
```

### Unbind

To unbind a service from an application, use the `cf unbind-service` command.

```
cf unbind-service myapp myobs
```

### Update

To update a service, use the `cf update-service` command.

```
cf update-service myobs -c '{
    "bucket_policy": "private",
    "status": "Enabled"
}'
```

You can also check the status of the service instance using the `cf service` command.

```
cf service myobs
```

### Deprovision

To deprovision the service, use the `cf delete-service` command.

```
cf delete-service myobs
```
