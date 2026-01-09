package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"workout.lavacro.net/database"
	"workout.lavacro.net/models"
)

var db *database.Dao

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

func readBody(r *http.Request, obj any) error {
	body, err := io.ReadAll(r.Body)
	defer func() {
		err = r.Body.Close()
	}()

	if err != nil {
		return err
	}

	println(string(body))
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, models.Envelope{"status": "ok"})
}

func allProgress(w http.ResponseWriter, r *http.Request) {
	var resp []models.Progress
	resp, dberr := db.AllProgress()

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"progress": resp})
}

func exercises(w http.ResponseWriter, r *http.Request) {
	var resp []models.Exercises
	resp, dberr := db.Exercises()

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"exercises": resp})
}

func months(w http.ResponseWriter, r *http.Request) {
	resp, dberr := db.YearMonths()
	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}
	writeResponse(w, models.Envelope{"dates": resp})
}

func getProgress(w http.ResponseWriter, r *http.Request) {
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

	resp, dberr := db.Progress(yearInt, monthInt)
	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"progress": resp})
}

func getActivity(w http.ResponseWriter, r *http.Request) {
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
	resp, dberr := db.Activity(idInt)

	if dberr != nil {
		log.Fatal("Problem getting data from database ...", dberr)
	}

	writeResponse(w, models.Envelope{"activity": resp})
}

func newActivity(w http.ResponseWriter, r *http.Request) {
	var act models.NewActivity
	err := readBody(r, &act)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, dberr := db.NewActivity(act)
	if dberr != nil {
		http.Error(w, dberr.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, models.Envelope{"id": id})
}

func deleteActivity(w http.ResponseWriter, r *http.Request) {
	log.Println("deleteActivity")
	id := r.URL.Query().Get("id")

	log.Println("id = ", id)
	if id == "" {
		log.Println("id is empty")
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	numId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("id is not an integer")
		http.Error(w, "id parameter must be an integer", http.StatusBadRequest)
		return
	}

	rows, err := db.DeleteActivity(numId)
	if err != nil {
		log.Println("error deleting activity")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("deleted rows = ", rows)
	writeResponse(w, models.Envelope{"deleted rows": rows})
}
