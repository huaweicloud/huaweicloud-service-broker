package dms

import (
	"github.com/pivotal-cf/brokerapi"
)

// Bind implematation
func (b *DMSBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	return brokerapi.Binding{}, nil
}
