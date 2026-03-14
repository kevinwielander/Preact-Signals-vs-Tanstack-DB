package main

import (
	_ "embed"
	"log"
	"net/http"

	"backend/api"
	"backend/events"
	"backend/ws"
)

//go:embed openapi.yaml
var openapiSpec []byte

func main() {
	hub := ws.NewHub()

	store := events.NewStore(func(evt events.Event) {
		hub.Broadcast(evt)
	})

	alarmHandler := api.NewAlarmHandler(store)
	resourceHandler := api.NewResourceHandler(store, alarmHandler.ListAll)

	mux := http.NewServeMux()

	// Alarm routes
	mux.HandleFunc("POST /alarms", alarmHandler.Create)
	mux.HandleFunc("GET /alarms", alarmHandler.List)
	mux.HandleFunc("GET /alarms/{id}", alarmHandler.Get)
	mux.HandleFunc("PATCH /alarms/{id}", alarmHandler.Update)
	mux.HandleFunc("GET /alarms/{id}/events", alarmHandler.GetEvents)

	// Resource routes
	mux.HandleFunc("POST /resources", resourceHandler.Create)
	mux.HandleFunc("GET /resources", resourceHandler.List)
	mux.HandleFunc("GET /resources/{id}", resourceHandler.Get)
	mux.HandleFunc("PATCH /resources/{id}", resourceHandler.Update)
	mux.HandleFunc("GET /resources/{id}/events", resourceHandler.GetEvents)

	// Resource -> Alarms
	mux.HandleFunc("GET /resources/{id}/alarms", resourceHandler.GetAlarms)

	// WebSocket
	mux.HandleFunc("GET /ws", hub.HandleWS)

	// OpenAPI spec
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(openapiSpec)
	})

	// Swagger UI
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html><head>
<title>API Docs</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head><body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>SwaggerUIBundle({url:"/openapi.yaml",dom_id:"#swagger-ui"})</script>
</body></html>`))
	})

	handler := corsMiddleware(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
