package api

import (
	"net/http"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/healthcheck", healthCheck)
	return mux
}
