package app

const (
	HealthStatusOK      = "ok"
	HealthStatusWarning = "warning"
	HealthStatusError   = "error"
)

type Health struct {
	Name      string `json:"name"`
	Status    string `json:"status"` // ok|warning|error
	Core      bool   `json:"core"`   // true for core dependencies
	LatencyMs int64  `json:"latency_ms"`
	Data      struct {
		Message string `json:"message,omitempty"`
		Code    int    `json:"code,omitempty"` // remote HTTP status code
	} `json:"data"`
}
