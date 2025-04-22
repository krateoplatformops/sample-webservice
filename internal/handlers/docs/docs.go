package docs

import (
	"net/http"
	"strings"
)

var WrapHandler = Handler()

type handler struct {
}

func Handler() http.Handler {
	return &handler{}
}

var _ http.Handler = (*handler)(nil)

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the request is for the OpenAPI JSON file
	if strings.HasSuffix(r.URL.Path, "openapi.json") {
		// Serve the OpenAPI JSON
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/v3/openapi.json")
		return
	}

	// Default: serve the Swagger UI HTML
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>API Documentation</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-standalone-preset.js"></script>
    <script>
        SwaggerUIBundle({
            url: './openapi.json',  // Use relative URL to work with any base path
            dom_id: '#swagger-ui',
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout"
        })
    </script>
</body>
</html>
    `))
}
