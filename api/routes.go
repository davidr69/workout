package api

import (
	"net/http"

	"workout.lavacro.net/database"
)

func Routes(dao *database.Dao) *http.ServeMux {
	db = dao
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/healthcheck", healthCheck)
	mux.HandleFunc("GET /api/v1/allprogress", allProgress)
	mux.HandleFunc("GET /api/v1/exercises", exercises)
	mux.HandleFunc("GET /api/v1/months", months)
	mux.HandleFunc("GET /api/v1/progress", getProgress)
	mux.HandleFunc("GET /api/v1/activity", getActivity)
	mux.HandleFunc("POST /api/v1/progress", newActivity)
	mux.HandleFunc("DELETE /api/v1/progress", deleteActivity)
	mux.HandleFunc("PUT /api/v1/progress", updateActivity)
	return mux
}
