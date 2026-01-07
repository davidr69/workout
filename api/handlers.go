package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"workout.lavacro.net/database"
	"workout.lavacro.net/models"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}
	rv, jerr := json.Marshal(data)

	if jerr != nil {
		http.Error(w, jerr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rv = append(rv, '\n')
	_, werr := w.Write(rv)

	if werr != nil {
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}
}

func exercises(w http.ResponseWriter, r *http.Request) {
	var resp []models.AllProgress
	resp, dberr := database.AllProgress()

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	js, jsErr := json.Marshal(resp)
	if jsErr != nil {
		log.Fatal("Problem encoding data ...", js)
	}

	_, respErr := fmt.Fprintln(w, js)
	if respErr != nil {
		log.Fatal("Problem writing response ...", respErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(js)
	if err != nil {
		log.Fatal("Problem writing response ...", err)
	}
}
