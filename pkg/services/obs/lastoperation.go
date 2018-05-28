package obs

import (
	"github.com/pivotal-cf/brokerapi"
)

// LastOperation implematation
func (b *OBSBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
