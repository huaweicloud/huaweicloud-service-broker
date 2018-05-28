package rds

import (
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *RDSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
