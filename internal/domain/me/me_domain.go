package me

// HealthResponse represents the JSON payload returned by the health check
type NyehResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
