package obs

import (
	"github.com/pivotal-cf/brokerapi"
)

// Unbind implematation
func (b *OBSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	return nil
}
