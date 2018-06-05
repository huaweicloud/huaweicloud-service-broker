package obs

import (
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

// Unbind implematation
func (b *OBSBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	// Log opts
	b.Logger.Debug(fmt.Sprintf("unbind obs bucket opts: instanceID: %s bindingID: %s", instanceID, bindingID))

	// TODO do something for unbind

	// Log result
	b.Logger.Debug(fmt.Sprintf("unbind obs bucket success: instanceID: %s bindingID: %s", instanceID, bindingID))

	return nil
}
