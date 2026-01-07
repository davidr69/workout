package api

import (
	"encoding/json"
	"log"
	"net/http"

	"workout.lavacro.net/database"
	"workout.lavacro.net/models"
)

func writeResponse(w http.ResponseWriter, r *http.Request, data models.Envelope) {
	js, jsErr := json.Marshal(data)
	if jsErr != nil {
		http.Error(w, jsErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	js = append(js, '\n')
	_, werr := w.Write(js)

	if werr != nil {
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, r, models.Envelope{"status": "ok"})
}

func exercises(w http.ResponseWriter, r *http.Request) {
	var resp []models.AllProgress
	resp, dberr := database.AllProgress()

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, r, models.Envelope{"progress": resp})
}
