package resources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/krateoplatformops/sample-webservice/internal/handlers"
)

type handler struct {
	handlers.HandlerOptions
}

func Get(opts handlers.HandlerOptions) http.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ http.Handler = (*handler)(nil)

// @Summary Sample API GET hardcoded resource
// @Description Get a hardcoded resource
// @ID get-resource
// @Param name query string true "Name of the resource"
// @Produce json
// @Success 200 {object} []handlers.Resource
// @Router /resource [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	// Log the request
	h.Log.Info("Received request", "name", name)

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	index := slices.IndexFunc(*h.ResourceStore, func(r handlers.Resource) bool {
		return r.Name == name
	})

	fmt.Println("Length of resource store:", len(*h.ResourceStore))

	if index == -1 {
		h.Log.Error("Resource not found", "name", name)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	res := *h.ResourceStore
	resource := res[index]
	// Write the response
	if err := json.NewEncoder(w).Encode(resource); err != nil {
		h.Log.Error("Failed to encode response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
