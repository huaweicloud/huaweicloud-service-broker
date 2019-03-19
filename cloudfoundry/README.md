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
