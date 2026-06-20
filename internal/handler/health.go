package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	c "nyeh-back/internal/core"
	d "nyeh-back/internal/domain"
)

// HealthCheckHandler godoc
//
//	@Summary	API Health Check
//	@Accept		json
//	@Produce	json
//	@Router		/healthCheck [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := d.HealthResponse{
		Status:  "ok",
		Port:    fmt.Sprintf("Running on port %v", c.Settings.PORT),
		Message: "API is fully operational",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
