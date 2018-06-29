package obs

import (
	"fmt"
	"strings"
)

// BuildBindingCredential from different obs bucket
func BuildBindingCredential(
	region string,
	url string,
	bucketname string,
	ak string,
	sk string,
	servicetype string) (BindingCredential, error) {

	// Example: https://{BucketName}.{Endpoint}
	parts := strings.Split(url, "//")
	if len(parts) != 2 {
		return BindingCredential{}, fmt.Errorf("unvalid url: %s", url)
	}
	link := fmt.Sprintf("%s//%s.%s", parts[0], bucketname, parts[1])

	// Init BindingCredential
	bc := BindingCredential{
		Region:     region,
		URL:        link,
		BucketName: bucketname,
		AK:         ak,
		SK:         sk,
		Type:       servicetype,
	}
	return bc, nil
}
