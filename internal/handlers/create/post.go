package resources

import (
	"encoding/json"
	"net/http"

	"github.com/krateoplatformops/sample-webservice/internal/handlers"
)

type handler struct {
	handlers.HandlerOptions
}

func Create(opts handlers.HandlerOptions) http.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ http.Handler = (*handler)(nil)

// @Summary Sample API GET hardcoded resource.
// @Description Create a resource. It returns 201 Created if the resource is created successfully.
// @ID create-resource
// @Produce json
// @Success 201 {object} []handlers.Resource
// @Router /resource [post]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the request body
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
	// Decode the request body
	var resource handlers.Resource
	if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
		h.Log.Error("Failed to decode request body", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Validate the resource
	if resource.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if resource.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}
	// Log the resource creation
	h.Log.Info("Creating resource", "name", resource.Name, "description", resource.Description)
	// Simulate resource creation
	// In a real application, you would save the resource to a database or perform some other action here.
	// For this example, we'll just log the resource creation and return a success response.
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		h.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
