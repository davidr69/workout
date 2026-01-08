package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

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

func allProgress(w http.ResponseWriter, r *http.Request) {
	var resp []models.Progress
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

func months(w http.ResponseWriter, r *http.Request) {
	resp, dberr := database.YearMonths()
	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}
	writeResponse(w, models.Envelope{"dates": resp})
}

func progress(w http.ResponseWriter, r *http.Request) {
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")

	if year == "" || month == "" {
		http.Error(w, "year and month parameters are required", http.StatusBadRequest)
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		http.Error(w, "year parameter must be an integer", http.StatusBadRequest)
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		http.Error(w, "month parameter must be an integer", http.StatusBadRequest)
		return
	}

	resp, dberr := database.Progress(yearInt, monthInt)
	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"progress": resp})
}

func activity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "id parameter must be an integer", http.StatusBadRequest)
		return
	}

	var resp models.Progress
	resp, dberr := database.Activity(idInt)

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"activity": resp})
}
