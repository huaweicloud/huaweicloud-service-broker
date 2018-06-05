package obs

// BuildBindingCredential from different obs bucket
func BuildBindingCredential(
	region string,
	url string,
	bucketname string,
	ak string,
	sk string,
	servicetype string) (BindingCredential, error) {

	// Init BindingCredential
	bc := BindingCredential{
		Region:     region,
		URL:        url,
		BucketName: bucketname,
		AK:         ak,
		SK:         sk,
		Type:       servicetype,
	}
	return bc, nil
}
