package resources

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/krateoplatformops/sample-webservice/internal/handlers"
)

type handler struct {
	handlers.HandlerOptions
}

func Update(opts handlers.HandlerOptions) http.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ http.Handler = (*handler)(nil)

// @Summary Sample API PATCH hardcoded resource.
// @Description Patch a resource. It returns 200 OK if the resource is updated successfully.
// @ID patch-resource
// @Produce json
// @Success 200 {object} []handlers.Resource
// @Router /resource [patch]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the request body
	if r.Method != http.MethodPatch {
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
	h.Log.Info("Pathing resource", "name", resource.Name, "description", resource.Description)

	index := slices.IndexFunc(*h.ResourceStore, func(r handlers.Resource) bool {
		return r.Name == resource.Name
	})

	if index == -1 {
		h.Log.Error("Resource not found", "name", resource.Name)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	res := *h.ResourceStore
	res[index] = resource

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		h.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
