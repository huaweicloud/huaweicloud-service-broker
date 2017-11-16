# brokerapi

[![Build Status](https://travis-ci.org/frodenas/brokerapi.svg?branch=master)](https://travis-ci.org/frodenas/brokerapi)

A go package for building V2 CF Service Brokers in Go. Depends on
[lager](https://github.com/pivotal-golang/lager) and
[gorilla/mux](https://github.com/gorilla/mux).

Requires go 1.4 or greater.

## usage

`brokerapi` defines a `ServiceBroker` interface with 7 methods. Simply create
a concrete type that implements these methods, and pass an instance of it to
`brokerapi.New`, along with a `lager.Logger` for logging and a
`brokerapi.BrokerCredentials` containing some HTTP basic auth credentials.

e.g.

```
package main

import (
    "github.com/frodenas/brokerapi"
    "github.com/pivotal-golang/lager"
)

type myServiceBroker struct {}

func (*myServiceBroker) Services() brokerapi.CatalogResponse {
    // Return the services's catalog offered by the broker
}

func (*myServiceBroker) Provision(instanceID string, details brokeapi.ProvisionDetails, acceptsIncomplete bool) (brokerapi.ProvisioningResponse, bool, error) {
    // Provision a new instance here
}

func (*myServiceBroker) Update(instanceID string, details brokerapi.UpdateDetails, acceptsIncomplete bool) (bool, error) {
    // Update instance here
}

func (*myServiceBroker) Deprovision(instanceID string, details brokeapi.DeprovisionDetails, acceptsIncomplete bool) (bool, error) {
    // Deprovision instance here
}

func (*myServiceBroker) Bind(instanceID string, bindingID string, details brokerapi.BindDetails) (brokerapi.BindingResponse, error) {
    // Bind to instance here
}

func (*myServiceBroker) Unbind(instanceID string, bindingID string, details brokerapi.UnbindDetails) error {
    // Unbind from instance here
}

func (*myServiceBroker) LastOperation(instanceID string) (brokeapi.LastOperationResponse, error) {
    // Return the status of the last instance operation
}

func main() {
    serviceBroker := &myServiceBroker{}
    logger := lager.NewLogger("my-service-broker")
    credentials := brokerapi.BrokerCredentials{
        Username: "username",
        Password: "password",
    }

    brokerAPI := brokerapi.New(serviceBroker, logger, credentials)
    http.Handle("/", brokerAPI)
    http.ListenAndServe(":3000", nil)
}
```

### errors

`brokerapi` defines a handful of error types in `service_broker.go` for some
common error cases that your service broker may encounter. Return these from
your `ServiceBroker` methods where appropriate, and `brokerapi` will do the
right thing, and give Cloud Foundry an appropriate status code, as per the V2
Service Broker API specification.

The error types are:

```
ErrInstanceAlreadyExists
ErrInstanceDoesNotExist
ErrInstanceLimitMet
ErrInstanceNotUpdateable
ErrInstanceNotBindable
ErrBindingAlreadyExists
ErrBindingDoesNotExist
ErrAsyncRequired
ErrAppGUIDRequired
```
