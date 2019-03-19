# Huawei Cloud Service Broker
[![Go Report Card](https://goreportcard.com/badge/github.com/huaweicloud/huaweicloud-service-broker?branch=master)](https://goreportcard.com/badge/github.com/huaweicloud/huaweicloud-service-broker)
[![Build Status](https://travis-ci.org/huaweicloud/huaweicloud-service-broker.svg?branch=master)](https://travis-ci.org/huaweicloud/huaweicloud-service-broker)
[![LICENSE](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://github.com/huaweicloud/huaweicloud-service-broker/blob/master/LICENSE)

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

## Getting Started

* [Getting Started on Locally](locally/README.md)
* [Getting Started on CCE](cce/README.md)
* [Getting Started on OpenShift](openshift/README.md)
* [Getting Started on Cloud Foundry](cloudfoundry/README.md)

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
