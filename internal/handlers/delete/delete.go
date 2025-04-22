package resources

import (
	"net/http"
	"slices"

	"github.com/krateoplatformops/sample-webservice/internal/handlers"
)

type handler struct {
	handlers.HandlerOptions
}

func Delete(opts handlers.HandlerOptions) http.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ http.Handler = (*handler)(nil)

// @Summary Sample API DELETE hardcoded resource.
// @Description Delete a resource. It returns 204 No Content if the resource is deleted successfully.
// @ID delete-resource
// @Param name query string true "Name of the resource"
// @Produce json
// @Success 204
// @Router /resource [delete]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	// Log the request
	h.Log.Info("Received request", "name", name)

	index := slices.IndexFunc(*h.ResourceStore, func(r handlers.Resource) bool {
		return r.Name == name
	})

	if index == -1 {
		h.Log.Error("", "name", name)

		return
	}

	res := *h.ResourceStore
	res = slices.Delete(res, index, index+1)

	// Simulate resource deletion
	// In a real application, you would delete the resource from a database or perform some other action here.
	// For this example, we'll just log the resource deletion and return a success response.
	w.WriteHeader(http.StatusNoContent)
}
