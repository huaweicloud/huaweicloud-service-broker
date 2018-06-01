package dms

import (
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

// Unbind implematation
func (b *DMSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	// Log opts
	b.Logger.Debug(fmt.Sprintf("unbind dms instance opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// TODO do something for unbind

	// Log result
	b.Logger.Debug(fmt.Sprintf("unbind dms instance success: instanceID: %s bindingID: %s", instanceID, bindingID))

	return nil
}
