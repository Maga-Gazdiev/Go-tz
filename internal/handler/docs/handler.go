package docs

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /docs/openapi.yaml", openAPISpec)
	mux.HandleFunc("GET /docs/", swaggerUI)
}

func openAPISpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/openapi.yaml")
}

func swaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <title>Subscriptions API documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>SwaggerUIBundle({url: "/docs/openapi.yaml", dom_id: "#swagger-ui"})</script>
</body>
</html>`))
}
