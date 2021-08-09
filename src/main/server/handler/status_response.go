package handler

type HealthCheckResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

var (
	activeHealthCheckResponse   = HealthCheckResponse{Status: "active", Version: "1.0.1"}
	inactiveHealthCheckResponse = HealthCheckResponse{Status: "inactive"}
)
