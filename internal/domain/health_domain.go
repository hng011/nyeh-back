package domain

// HealthResponse represents the JSON payload returned by the health check
type HealthResponse struct {
	Status  string `json:"status"`
	Port    string `json:"port"`
	Message string `json:"message"`
}
