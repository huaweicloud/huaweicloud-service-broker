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
