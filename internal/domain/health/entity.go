package health

import "time"

type HealthStatus struct {
	Status       string                 `json:"status"`
	Timestamp    time.Time              `json:"timestamp"`
	Service      string                 `json:"service"`
	Database     string                 `json:"database"`
	DatabaseInfo map[string]interface{} `json:"database_info,omitempty"`
	Error        string                 `json:"error,omitempty"`
}
