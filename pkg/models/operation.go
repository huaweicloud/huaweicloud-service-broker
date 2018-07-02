package models

const (
	// OperationProvisioning for asyc provision
	OperationProvisioning string = "provisioning"

	// OperationDeprovisioning for asyc deprovision
	OperationDeprovisioning string = "deprovisioning"

	// OperationUpdating for asyc update
	OperationUpdating string = "updating"
)

const (
	// OperationAsyncDCS for default way
	OperationAsyncDCS bool = true

	// OperationAsyncDMSQueue for default way
	OperationAsyncDMSQueue bool = false

	// OperationAsyncDMSInstance for default way
	OperationAsyncDMSInstance bool = true

	// OperationAsyncOBS for default way
	OperationAsyncOBS bool = false

	// OperationAsyncRDS for default way
	OperationAsyncRDS bool = true
)
