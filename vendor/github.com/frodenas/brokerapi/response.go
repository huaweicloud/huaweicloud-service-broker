package brokerapi

type EmptyResponse struct{}

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"description"`
}

type CatalogResponse struct {
	Services []Service `json:"services"`
}

type ProvisioningResponse struct {
	DashboardURL string `json:"dashboard_url,omitempty"`
}

type BindingResponse struct {
	Credentials    interface{} `json:"credentials"`
	SyslogDrainURL string      `json:"syslog_drain_url,omitempty"`
}

type CredentialsHash struct {
	Host     string `json:"host,omitempty"`
	Port     int64  `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URI      string `json:"uri,omitempty"`
	JDBCURI  string `json:"jdbcUrl,omitempty"`
}

type LastOperationResponse struct {
	State       string `json:"state"`
	Description string `json:"description,omitempty"`
}

const LastOperationInProgress = "in progress"
const LastOperationFailed = "failed"
const LastOperationSucceeded = "succeeded"
