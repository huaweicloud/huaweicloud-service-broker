package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/pivotal-cf/brokerapi"

	"github.com/huaweicloud/huaweicloud-service-broker/pkg/broker"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/config"
	"github.com/huaweicloud/huaweicloud-service-broker/pkg/logger"
)

var (
	configFilePath string
	port           string
)

func init() {
	flag.StringVar(&configFilePath, "config", "", "Location of Service Broker config file")
	flag.StringVar(&port, "port", "3000", "Service Broker listen port")
}

func main() {

	flag.Parse()

	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading Service Broker config file: %s", err)
	}

	logger := logger.BuildLogger(config.BrokerConfig.LogLevel)

	serviceBroker, err := broker.New(logger, config)
	if err != nil {
		log.Fatalf("Error new Service Broker: %s", err)
	}

	credentials := brokerapi.BrokerCredentials{
		Username: config.BrokerConfig.Username,
		Password: config.BrokerConfig.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
	http.Handle("/", brokerAPI)

	fmt.Println("### Service Broker started on port ###", port)
	http.ListenAndServe(":"+port, nil)
}
