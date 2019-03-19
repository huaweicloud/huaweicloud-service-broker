## Getting Started on OpenShift

### Installation

#### 1. Making configurations for Service Broker

The broker can be deployed to an already existing [OpenShift](https://www.openshift.com/) Cluster
which has enabled [Service Catalog](https://github.com/kubernetes-incubator/service-catalog/).

We recommend to finish the following installation in the OpenShift Cluster Master Node.

```
$ git clone https://github.com/huaweicloud/huaweicloud-service-broker.git
$ cd huaweicloud-service-broker/openshift/deploy/
$ vi config.json
```
Modify the file ```config.json``` to include your own configurations. Different Clouds have different configurations.
See [configuration.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration.md) for instructions.

The [default authorization](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/config.json#L4-L5) for visiting huaweicloud service broker is as following:

```
    "broker_config": {
        "log_level": "DEBUG",
        "username": "username",
        "password": "password"
    },
```

#### 2. Creating secrets in OpenShift Cluster

Encode the authorization ```username``` and ```password``` by using base64.

```
$ echo -n "username" | base64
$ echo -n "password" | base64
```

Update the key ```username``` value and the key ```password``` value in [secret-auth.yaml](
https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/secret-auth.yaml#L7-L8) file by using the result of the above two commands, and then you can run the following command to create ```secret-auth.yaml```.

```
$ oc create -f secret-auth.yaml
```

Encode the file ```config.json``` content by using base64.

```
$ base64 -w 0 config.json
```

Update the key ```config.json``` value in [secret-config.yaml](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/secret-config.yaml#L7) file by using the result of the above command, and then you can run the following command to create ```secret-config.yaml```.

```
$ oc create -f secret-config.yaml
```

#### 3. Creating service in OpenShift Cluster

You can find more information by openning the file [service.yaml](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service.yaml). The default [name](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service.yaml#L77) for the service is ```service-broker```, and the default ```namespace``` for the service is ```default```, and the default [port](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service.yaml#L85) for the service is ```12345```, and then you can run the following command to create ```service.yaml```. Please make sure the service is running before going to the next step.

```
$ oc create -f serivce.yaml
```

#### 4. Creating service broker in OpenShift Cluster

If you do not use the default configurations in Step 3, you can update the key ```url``` value in [service-broker.yaml](
https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service-broker.yaml#L12) file by the service ```name```, ```namespace``` and ```port``` which are created by Step 3.

```
$ vi service-broker.yaml
```
Then you can run the following command to create ```service-broker.yaml```.

```
$ oc create -f service-broker.yaml
```

If the service broker is created successfully,
you can find a service broker named ```cluster-service-broker```
by running the following command.

```
$ oc get clusterservicebrokers
```

You can also find the lastest Services ```clusterserviceclasses``` and Service Plans ```clusterserviceplans``` by running the following command.

```
$ oc get clusterserviceclasses -o=custom-columns=SERVICE\ NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName
$ oc get clusterserviceplans -o=custom-columns=NAME:.metadata.name,EXTERNAL\ NAME:.spec.externalName,SERVICE\ CLASS:.spec.clusterServiceClassRef.name --sort-by=.spec.clusterServiceClassRef.name
```

Currently the following Services and Service Plans are tested in OpenShift.
<table>
  <tr align="left">
    <th width="20%">Service Name</th>
    <th width="30%">Service Description</th>
    <th width="20%">Service Plan Name</th>
    <th width="30%">Service Plan Description</th>
  </tr>
  <tr>
    <td rowspan="5">rds-mysql</td>
    <td rowspan="5">RDS MySQL Service</td>
    <td>5-7-17</td>
    <td>RDS MySQL 5.7.17</td>
  </tr>
  <tr>
    <td>5-6-35</td>
    <td>RDS MySQL 5.6.35</td>
  </tr>
  <tr>
    <td>5-6-34</td>
    <td>RDS MySQL 5.6.34</td>
  </tr>
  <tr>
    <td>5-6-33</td>
    <td>RDS MySQL 5.6.33</td>
  </tr>
  <tr>
    <td>5-6-30</td>
    <td>RDS MySQL 5.6.30</td>
  </tr>
  <tr>
    <td>rds-sqlserver</td>
    <td>RDS SQLServer Service</td>
    <td>2014-SP2-SE</td>
    <td>RDS SQLServer 2014 SP2 SE</td>
  </tr>
  <tr>
    <td rowspan="2">dcs-redis</td>
    <td rowspan="2">Distributed Cache Service for Redis</td>
    <td>SingleNode</td>
    <td>Redis Single Node</td>
  </tr>
  <tr>
    <td>MasterStandby</td>
    <td>Redis Master Standby</td>
  </tr>
</table>

### Usage

Application Developers can start to consume the services
by creating ```ServiceInstance``` and ```ServiceBinding``` resources. 
Take MySQL as an example.

#### 1. Creating ServiceInstance in OpenShift Cluster

```
$ cd ./../examples/mysql/
$ vi mysql-service-instance.yaml
```

The ```mysql-service-instance.yaml``` example is using the Service ```rds-mysql``` and Service Plan ```5-7-17```. About the key [parameters](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/examples/mysql/mysql-service-instance.yaml#L12), you can find more informations in the [rds-mysql.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/usage/rds-mysql.md#provision-parameters). Then you can run the following command to create ```mysql-service-instance.yaml```.

```
$ oc create -f mysql-service-instance.yaml
```

The following command will get more informations about the created ```mysql-service-instance```. Please make sure the ```Status``` of ```mysql-service-instance``` is OK before going to the next step.

```
$ oc describe serviceinstance mysql-service-instance
```

#### 2. Creating ServiceBinding in OpenShift Cluster

```
$ vi mysql-service-binding.yaml
$ oc create -f mysql-service-binding.yaml
```

This example will store the binding informations into a secret resource named ```mysql-service-secret```.

#### 3. Using ServiceBinding in Pod

```
$ oc create -f pod.yaml
```

The ```pod.yaml``` will use ```mysql-service-secret``` and mount it as a volume so that the nginx application can use the binding informations as an input.
