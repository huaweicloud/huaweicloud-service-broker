package dcs

import (
	"github.com/pivotal-cf/brokerapi"
)

// Unbind implematation
func (b *DCSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}
