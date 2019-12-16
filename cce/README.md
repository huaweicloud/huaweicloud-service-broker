## Getting Started on CCE

### Prerequisites

- [CCE](https://support.huaweicloud.com/en-us/productdesc-cce/cce_productdesc_0001.html) 1.9+ with RBAC enabled.
- [Helm](https://github.com/helm/helm#install) 2.7.0+ has been installed.
- [Kubernetes Service Catalog](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md) has been installed.

### Installing the Chart

For convenience we can export environment variables firstly.
```
export backDatabase_databaseHost=<back database host ip>
export backDatabase_databasePort=<back database port>
export backDatabase_databaseName=<back database name>
export backDatabase_databaseUsername=<back database username>
export backDatabase_databasePassword=<back database password>
export cloudCredentials_authUrl=<auth url for cloud>
export cloudCredentials_username=<username for cloud>
export cloudCredentials_password=<password for cloud>
export cloudCredentials_domainName=<domain name for cloud>
export cloudCredentials_tenantName=<tenant name for cloud>
export cloudCredentials_region=<region name for cloud>
export cloudCredentials_accessKey=<access key for cloud>
export cloudCredentials_secretKey=<secret key for cloud>
export catalog_primaryAvailabilityZone=<primary availability zone>
export catalog_secondaryAvailabilityZone=<secondary availability zone>
export catalog_vpcID=<vpc id>
export catalog_subnetID=<subnet id>
export catalog_securityGroupID=<security group id>
```

Installation of this chart is by helm.

```
$ git clone https://github.com/huaweicloud/huaweicloud-service-broker.git
$ cd huaweicloud-service-broker
$ helm install cce/charts/ --name service-broker --namespace huaweicloud \
  --set backDatabase.databaseHost=$backDatabase_databaseHost \
  --set backDatabase.databasePort=$backDatabase_databasePort \
  --set backDatabase.databaseName=$backDatabase_databaseName \
  --set backDatabase.databaseUsername=$backDatabase_databaseUsername \
  --set backDatabase.databasePassword=$backDatabase_databasePassword \
  --set cloudCredentials.authUrl=$cloudCredentials_authUrl \
  --set cloudCredentials.username=$cloudCredentials_username \
  --set cloudCredentials.password=$cloudCredentials_password \
  --set cloudCredentials.domainName=$cloudCredentials_domainName \
  --set cloudCredentials.tenantName=$cloudCredentials_tenantName \
  --set cloudCredentials.region=$cloudCredentials_region \
  --set cloudCredentials.accessKey=$cloudCredentials_accessKey \
  --set cloudCredentials.secretKey=$cloudCredentials_secretKey \
  --set catalog.primaryAvailabilityZone=$catalog_primaryAvailabilityZone \
  --set catalog.secondaryAvailabilityZone=$catalog_secondaryAvailabilityZone \
  --set catalog.vpcID=$catalog_vpcID \
  --set catalog.subnetID=$catalog_subnetID \
  --set catalog.securityGroupID=$catalog_securityGroupID
```

please see the following configurable parameters that can be configured during installation.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| service.replicas | service replicas count | 1 |
| service.image | service image name and version | quay.io/huaweicloud/huaweicloud-service-broker:latest |
| service.imagePullPolicy | service image pull policy: IfNotPresent, Always, or Never | Always |
| service.containerPort | service container port | 3000 |
| brokerConfig.logLevel | broker config log level | "DEBUG" |
| brokerConfig.username | broker auth username | "username" |
| brokerConfig.password | broker auth password | "password" |
| backDatabase.databaseType | back database type | "mysql" |
| backDatabase.databaseHost | back database host ip | "127.0.0.1" |
| backDatabase.databasePort | back database port | 3306 |
| backDatabase.databaseName | back database name | "broker" |
| backDatabase.databaseUsername | back database username | "******" |
| backDatabase.databasePassword | back database password | "******" |
| cloudCredentials.authUrl | auth url for cloud | "https://iam.eu-west-0.prod-cloud-ocb.orange-business.com/v3" |
| cloudCredentials.username | username for cloud | "******" |
| cloudCredentials.password | password for cloud | "******" |
| cloudCredentials.domainName | domain name for cloud | "******" |
| cloudCredentials.tenantName | tenant name for cloud | "eu-west-0" |
| cloudCredentials.region | region name for cloud | "eu-west-0" |
| cloudCredentials.accessKey | access key for cloud | "******" |
| cloudCredentials.secretKey | secret key for cloud | "******" |
| cloudCredentials.rds_version | rds version | "******" |
| catalog.primaryAvailabilityZone | primary availability zone | "eu-west-0a" |
| catalog.secondaryAvailabilityZone | secondary availability zone | "eu-west-0b" |
| catalog.vpcID | vpc id | "******" |
| catalog.subnetID | subnet id | "******" |
| catalog.securityGroupID | security group id | "******" |

If the service broker is created successfully,
you can find a service broker named ```cluster-service-broker```
by running the following command.

```
$ kubectl get clusterservicebrokers
```

You can also find the lastest Services ```clusterserviceclasses``` and Service Plans ```clusterserviceplans``` by running the following command.

```
$ kubectl get clusterserviceclasses -o=custom-columns=SERVICE\ NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName
$ kubectl get clusterserviceplans -o=custom-columns=NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName,SERVICE\ CLASS:.spec.clusterServiceClassRef.name --sort-by=.spec.clusterServiceClassRef.name
```

### Uninstalling the Chart

```
$ helm delete --purge service-broker
$ kubectl delete namespace huaweicloud
```

### Usage

Application Developers can start to consume the services
by creating ```ServiceInstance``` and ```ServiceBinding``` resources. 
Take MySQL as an example.

#### 1. Creating ServiceInstance in CCE

```
$ cd cce/examples/mysql/
$ vi mysql-service-instance.yaml
```

The ```mysql-service-instance.yaml``` example is using the Service ```rds-mysql``` and Service Plan ```5-7-17```. About the key [parameters](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/cce/examples/mysql/mysql-service-instance.yaml#L12), you can find more informations in the [rds-mysql.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/usage/rds-mysql.md#provision-parameters). Then you can run the following command to create ```mysql-service-instance.yaml```.

```
$ kubectl create -f mysql-service-instance.yaml
```

The following command will get more informations about the created ```mysql-service-instance```. Please make sure the ```Status``` of ```mysql-service-instance``` is OK before going to the next step.

```
$ kubectl describe serviceinstance mysql-service-instance
```

#### 2. Creating ServiceBinding in CCE

```
$ vi mysql-service-binding.yaml
$ kubectl create -f mysql-service-binding.yaml
```

This example will store the binding informations into a secret resource named ```mysql-service-secret```.

#### 3. Using ServiceBinding in Pod

```
$ kubectl create -f pod.yaml
```

The ```pod.yaml``` will use ```mysql-service-secret``` and mount it as a volume so that the nginx application can use the binding informations as an input.
