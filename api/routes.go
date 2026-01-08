package api

import (
	"net/http"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/healthcheck", healthCheck)
	mux.HandleFunc("GET /api/v1/allprogress", allProgress)
	mux.HandleFunc("GET /api/v1/exercises", exercises)
	mux.HandleFunc("GET /api/v1/months", months)
	mux.HandleFunc("GET /api/v1/progress", progress)
	mux.HandleFunc("GET /api/v1/activity", activity)
	/*
		mux.HandleFunc("GET /api/v1/stats/{when}", nil)
		mux.HandleFunc("POST /api/v1/activity", nil)
		mux.HandleFunc("PUT /api/v1/activity", nil)
		mux.HandleFunc("DELETE /api/v1/activity", nil)
	*/
	return mux
}
