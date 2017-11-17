package rds

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("instances")
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("instances", id)
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("instances", id)
}

func updateVolumeURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("instances", id, "action")
}
