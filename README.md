# Huawei Cloud Service Broker
[![Go Report Card](https://goreportcard.com/badge/github.com/huaweicloud/huaweicloud-service-broker?branch=master)](https://goreportcard.com/badge/github.com/huaweicloud/huaweicloud-service-broker)
[![Build Status](https://travis-ci.org/huaweicloud/huaweicloud-service-broker.svg?branch=master)](https://travis-ci.org/huaweicloud/huaweicloud-service-broker)
[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/LICENSE)

Test OpenLab CI

This is a [Service Broker](https://docs.cloudfoundry.org/services/overview.html) for Huawei Cloud.
It also can be used for Flexible Engine and Open Telekom Cloud.
Currently it includes the following services support:
* [Distributed Cache Service (DCS)](http://www.huaweicloud.com/en-us/product/dcs.html)
* [Distributed Message Service (DMS)](http://www.huaweicloud.com/en-us/product/dms.html)
* [Object Storage Service (OBS)](http://www.huaweicloud.com/en-us/product/obs.html)
* [Relational Database Service (RDS)](http://www.huaweicloud.com/en-us/product/rds.html)

Different Clouds have different services, the users will be required to make their own configurations in the following ```Installation``` so that this Service Broker can be running normally.

## Prerequisites

You'll need a few prerequisites before you are getting started.

Note: You can setup a backing MySQL Database by using the following choices.
We recommend to use the choice 1.

### Choice 1: Setup a Backing MySQL Database by RDS

* Setup MySQL instance by the Huawei Cloud Relational Database Service (RDS)
* Login in MySQL with your account and password
* Create database instance by running the following command
    ```
    CREATE DATABASE broker;
    ```

### Choice 2: Setup a Backing MySQL Database by yourself

* Setup a MySQL Server and make sure that the database can be accessed by service broker
* Login in MySQL with your account and password
* Create database instance by running the following command
    ```
    CREATE DATABASE broker;
    ```
* Create database user by running the following command
    ```
    CREATE USER 'username'@'%' IDENTIFIED BY 'password';
    ```
* Grant privileges to the user by running the following command
    ```
    GRANT ALL PRIVILEGES ON broker.* TO 'username'@'%' WITH GRANT OPTION;
    FLUSH PRIVILEGES;
    ```
* Make sure MySQL can be connected remotely
    ```
    vi /etc/mysql/mysql.conf.d/mysqld.cnf
    #bind-address           = 127.0.0.1
    ```
* Restart MySQL service

The information of database instance will be used in the following Installation.

## Getting Started on Locally

### Installation

Download the [configuration file](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/config.json)
and modify the configuration file to include your own configurations. Different Clouds have different configurations.
See [configuration.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration.md) for instructions.


Using the standard `go get` (you must have [Go](https://golang.org/) already installed in your local machine):

```
$ go get github.com/huaweicloud/huaweicloud-service-broker
$ huaweicloud-service-broker -config=config.json -port=3000
```

### Usage

Application Developers can start to test the services locally. The [lifecycle.sh](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/scripts/lifecycle.sh) can guide you to test by scripts.

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

You can find more information by openning the file [service.yaml](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service.yaml). The default [port](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service.yaml#L85) for the service is ```12345```, and then you can run the following command to create ```service.yaml```. Please make sure the service is running before going to the next step.

```
$ oc create -f serivce.yaml
```

#### 4. Creating service broker in OpenShift Cluster

Firstly you need to get the service ```ClusterIP``` and ```Port``` which are created by Step 3.

```
$ oc get svc | grep service-broker
$ vi service-broker.yaml
```

Update the key ```url``` value in [service-broker.yaml](
https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/openshift/deploy/service-broker.yaml#L7) file by the service ```ClusterIP``` and ```Port```. If it is possible, you can register ```ClusterIP``` and ```Port``` into the DNS Server so that you can use the domain name as the key ```url``` value. Then you can run the following command to create ```service-broker.yaml```.

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

## Getting Started on Cloud Foundry

### Installation

The broker can be deployed to an already existing [Cloud Foundry](https://www.cloudfoundry.org/) installation:

```
$ git clone https://github.com/huaweicloud/huaweicloud-service-broker.git
$ cd huaweicloud-service-broker
```

Modify the [configuration file](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/config.json) to include your own configurations.
Different Clouds have different configurations.
See [configuration.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration.md) for instructions.

Then you can push the broker to your [Cloud Foundry](https://www.cloudfoundry.org/) environment:

```
$ cf push huaweicloud-service-broker
```

[Register the broker](https://docs.cloudfoundry.org/services/managing-service-brokers.html#register-broker) within your Cloud Foundry installation. For example:

```
$ cf create-service-broker huaweicloud-service-broker username password https://huaweicloud-service-broker.example.com
```

Make sure that the service broker is registered successfully:

```
$ cf service-brokers
```

Display access to service:

```
$ cf service-access
```

[Make Services and Plans public](https://docs.cloudfoundry.org/services/access-control.html#enable-access).
 For example, enable rds-mysql service:

```
$ cf enable-service-access rds-mysql
```

### Usage

Application Developers can start to consume the services using the standard [CF CLI commands](https://docs.cloudfoundry.org/devguide/services/managing-services.html). See [usage.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/usage.md) for instructions.

## Contributing

In the spirit of [free software](http://www.fsf.org/licensing/essays/free-sw.html), **everyone** is encouraged to help improve this project.

Here are some ways *you* can contribute:

* by using prerelease versions or master branch
* by reporting bugs
* by suggesting new features
* by writing or editing documentation
* by writing specifications
* by writing code (**no patch is too small**: fix typos, add comments, clean up inconsistent whitespace)
* by refactoring code
* by closing [issues](https://github.com/huaweicloud/huaweicloud-service-broker/issues)
* by reviewing patches

### Submitting an Issue

We use the [GitHub issue tracker](https://github.com/huaweicloud/huaweicloud-service-broker/issues) to track bugs and features. Before submitting a bug report or feature request, check to make sure it hasn't already been submitted. You can indicate support for an existing issue by voting it up. When submitting a bug report, please include a [Gist](http://gist.github.com/) that includes a stack trace and any details that may be necessary to reproduce the bug, including your Golang version and operating system. Ideally, a bug report should include a pull request with failing specs.

### Submitting a Pull Request

1. Fork the project
2. Create a topic branch
3. Implement your feature or bug fix
4. Commit and push your changes
5. Submit a pull request

## License

huaweicloud-service-broker is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
