package resources

import (
	"encoding/json"
	"net/http"

	"github.com/krateoplatformops/sample-webservice/internal/handlers"
)

type handler struct {
	handlers.HandlerOptions
}

func List(opts handlers.HandlerOptions) http.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

type ListResponse struct {
	Resources []handlers.Resource `json:"resources"`
	Count     int                 `json:"count"`
}

var _ http.Handler = (*handler)(nil)

// @Summary Sample API GET list hardcoded resources
// @Description List hardcoded resources
// @ID list-resource
// @Produce json
// @Success 200 {object} []ListResponse
// @Router /resources [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the request
	h.Log.Info("Received request")

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Create a list response
	response := ListResponse{
		Resources: []handlers.Resource{
			{
				Name:        "Resource 1",
				Description: "This is the first sample resource",
			},
			{
				Name:        "Resource 2",
				Description: "This is the second sample resource",
			},
		},
		Count: 1,
	}

	h.Log.Info("List response", "resources", response.Resources, "count", response.Count)

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
