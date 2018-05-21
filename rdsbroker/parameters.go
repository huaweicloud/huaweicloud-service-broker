package rdsbroker

type ProvisionParameters struct {
}

type UpdateParameters struct {
	VolumeSize int64 `mapstructure:"volume_size"`
}

type BindParameters struct {
}
