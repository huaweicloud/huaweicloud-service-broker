package dms

import (
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *DMSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
