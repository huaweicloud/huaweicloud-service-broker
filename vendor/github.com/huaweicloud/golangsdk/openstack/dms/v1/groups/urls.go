package groups

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/golangsdk"
)

// endpoint/queues/{queue_id}/groups
const dmsPath = "dms/"
const resourcePathQueues = "queues"
const resourcePathGroups = "groups"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient, queueID string) string {
	return strings.Replace(client.ServiceURL(resourcePathQueues, queueID, resourcePathGroups), dmsPath, "", -1)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, queueID string, groupID string) string {
	return strings.Replace(client.ServiceURL(resourcePathQueues, queueID, resourcePathGroups, groupID), dmsPath, "", -1)
}

// listURL will build the list url of list function
func listURL(client *golangsdk.ServiceClient, queueID string, includeDeadLetter bool) string {
	return strings.Replace(client.ServiceURL(resourcePathQueues, queueID, fmt.Sprintf("%s?include_deadletter=%t", resourcePathGroups, includeDeadLetter)), dmsPath, "", -1)
}
