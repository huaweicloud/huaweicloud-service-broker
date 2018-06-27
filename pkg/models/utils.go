package models

import (
	"encoding/json"
	"fmt"
)

// ToJson for log
func ToJson(obj interface{}) string {
	// Marshal
	bs, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf("marshal rds instance failed. Error: %s", err)
	}
	return string(bs)
}
