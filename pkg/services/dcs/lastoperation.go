package dcs

import (
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *DCSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
