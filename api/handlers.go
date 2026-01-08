package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"workout.lavacro.net/database"
	"workout.lavacro.net/models"
)

func writeResponse(w http.ResponseWriter, data models.Envelope) {
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

func readBody(r *http.Request, obj any) (any, error) {
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(obj); err != nil {
		return nil, err
	}

	err := dec.Decode(&struct{}{})
	if err != io.EOF {
		return nil, errors.New("trailing garbage after JSON object")
	}

	return dec, nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, models.Envelope{"status": "ok"})
}

func progress(w http.ResponseWriter, r *http.Request) {
	var resp []models.AllProgress
	resp, dberr := database.AllProgress()

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"progress": resp})
}

func exercises(w http.ResponseWriter, r *http.Request) {
	var resp []models.Exercises
	resp, dberr := database.Exercises()

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"exercises": resp})
}
