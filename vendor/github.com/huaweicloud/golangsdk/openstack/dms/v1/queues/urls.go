package queues

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/golangsdk"
)

// endpoint/queues
const dmsPath = "dms/"
const resourcePath = "queues"

// createURL will build the rest query url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return strings.Replace(client.ServiceURL(resourcePath), dmsPath, "", -1)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return strings.Replace(client.ServiceURL(resourcePath, id), dmsPath, "", -1)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string, includeDeadLetter bool) string {
	return strings.Replace(client.ServiceURL(resourcePath, fmt.Sprintf("%s?include_deadletter=%t", id, includeDeadLetter)), dmsPath, "", -1)
}

// listURL will build the list url of list function
func listURL(client *golangsdk.ServiceClient, includeDeadLetter bool) string {
	return strings.Replace(client.ServiceURL(fmt.Sprintf("%s?include_deadletter=%t", resourcePath, includeDeadLetter)), dmsPath, "", -1)
}
