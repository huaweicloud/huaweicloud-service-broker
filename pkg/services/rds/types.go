package rds

// BindingCredential represent rds binding credential
type BindingCredential struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Name     string `json:"name,omitempty"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	URI      string `json:"uri,omitempty"`
	Type     string `json:"type,omitempty"`
}
