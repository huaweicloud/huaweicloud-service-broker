# Huawei Cloud Service Broker
[![Go Report Card](https://goreportcard.com/badge/github.com/huaweicloud/huaweicloud-service-broker?branch=master)](https://goreportcard.com/badge/github.com/huaweicloud/huaweicloud-service-broker)
[![Build Status](https://travis-ci.org/huaweicloud/huaweicloud-service-broker.svg?branch=master)](https://travis-ci.org/huaweicloud/huaweicloud-service-broker)
[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/LICENSE)

This is a [Service Broker](https://docs.cloudfoundry.org/services/overview.html) for Huawei Cloud.
Currently it includes the following services support:
* [Distributed Cache Service (DCS)](http://www.huaweicloud.com/en-us/product/dcs.html)
* [Distributed Message Service (DMS)](http://www.huaweicloud.com/en-us/product/dms.html)
* [Object Storage Service (OBS)](http://www.huaweicloud.com/en-us/product/obs.html)
* [Relational Database Service (RDS)](http://www.huaweicloud.com/en-us/product/rds.html)

## Prerequisites

You'll need a few prerequisites before you are getting started.

### Setup a Backing MySQL Database

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
and modify the configuration file to include your own configurations. See [configuration.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration.md) for instructions.


Using the standard `go install` (you must have [Go](https://golang.org/) already installed in your local machine):

```
$ go install github.com/huaweicloud/huaweicloud-service-broker
$ huaweicloud-service-broker -config=config.json -port=3000
```

### Usage

Application Developers can start to test the services locally. The [lifecycle.sh](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/scripts/lifecycle.sh) can guide you to test by scripts.

## Getting Started on Cloud Foundry

### Installation

The broker can be deployed to an already existing [Cloud Foundry](https://www.cloudfoundry.org/) installation:

```
$ git clone https://github.com/huaweicloud/huaweicloud-service-broker.git
$ cd huaweicloud-service-broker
```

Modify the [configuration file](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/config.json) to include your own configurations. See [configuration.md](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/docs/configuration.md) for instructions.

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
