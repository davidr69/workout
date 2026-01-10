package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"workout.lavacro.net/database"
	"workout.lavacro.net/models"
)

var db *database.Dao

func writeResponse(w http.ResponseWriter, data models.Envelope) {
	js, jsErr := json.Marshal(data)
	if jsErr != nil {
		slog.Error("Error marshalling JSON", "error", jsErr)
		http.Error(w, jsErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	js = append(js, '\n')
	_, werr := w.Write(js)

	if werr != nil {
		slog.Error("Error writing response", "error", werr)
		http.Error(w, werr.Error(), http.StatusInternalServerError)
	}
}

func errorWriter(w http.ResponseWriter, err string) {
	slog.Error("Error", "message", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, werr := w.Write([]byte(`{"error": "` + err + `"}` + "\n"))
	if werr != nil {
		slog.Error("Error writing response", "error", werr)
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

	slog.Info("Request body", "body", string(body))
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
	slog.Info("allProgress")
	var resp []models.Progress
	resp, dberr := db.AllProgress()

	if dberr != nil {
		slog.Error("Error getting progress ...", dberr)
		errorWriter(w, "database error")
		return
	}

	writeResponse(w, models.Envelope{"progress": resp})
}

func exercises(w http.ResponseWriter, r *http.Request) {
	slog.Info("exercises")
	var resp []models.Exercises
	resp, dberr := db.Exercises()

	if dberr != nil {
		slog.Error("exercises:", "error", dberr.Error())
		errorWriter(w, "database error")
		return
	}

	writeResponse(w, models.Envelope{"exercises": resp})
}

func months(w http.ResponseWriter, r *http.Request) {
	slog.Info("months")
	resp, dberr := db.YearMonths()
	if dberr != nil {
		slog.Error("months:", "error", dberr.Error())
		errorWriter(w, "database error")
		return
	}
	writeResponse(w, models.Envelope{"dates": resp})
}

func getProgress(w http.ResponseWriter, r *http.Request) {
	slog.Info("getProgress")
	year := r.URL.Query().Get("year")
	month := r.URL.Query().Get("month")

	if year == "" || month == "" {
		errorWriter(w, "year and month parameters are required")
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		errorWriter(w, "year parameter must be an integer")
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		errorWriter(w, "month parameter must be an integer")
		return
	}

	slog.Info("getProgress:", "year", yearInt, "month", monthInt)
	resp, dberr := db.Progress(yearInt, monthInt)
	if dberr != nil {
		slog.Error("Error getting progress ...", dberr)
		errorWriter(w, "database error")
	}

	writeResponse(w, models.Envelope{"progress": resp})
}

func getActivity(w http.ResponseWriter, r *http.Request) {
	slog.Info("getActivity")
	id := r.URL.Query().Get("id")

	if id == "" {
		errorWriter(w, "id parameter is required")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		errorWriter(w, "id parameter must be an integer")
		return
	}

	var resp models.Progress
	resp, dberr := db.Activity(idInt)

	if dberr != nil {
		slog.Error("Error getting activity ...", dberr)
		errorWriter(w, "database error")
		return
	}

	writeResponse(w, models.Envelope{"activity": resp})
}

func newActivity(w http.ResponseWriter, r *http.Request) {
	slog.Info("newActivity")
	var act models.Activity
	err := readBody(r, &act)
	if err != nil {
		errorWriter(w, err.Error())
		return
	}

	id, dberr := db.NewActivity(act)
	if dberr != nil {
		slog.Error("Error inserting activity ...", dberr)
		errorWriter(w, "database error")
	} else {
		writeResponse(w, models.Envelope{"id": id})
	}
}

func updateActivity(w http.ResponseWriter, r *http.Request) {
	slog.Info("updateActivity")
	var act models.Activity
	err := readBody(r, &act)
	if err != nil {
		errorWriter(w, err.Error())
		return
	}

	count, uerr := db.UpdateActivity(act)
	if uerr != nil {
		slog.Error("Error updating activity ...", uerr)
		errorWriter(w, "database error")
		return
	}

	writeResponse(w, models.Envelope{"updated": count})
}

func deleteActivity(w http.ResponseWriter, r *http.Request) {
	slog.Info("deleteActivity")
	id := r.URL.Query().Get("id")

	slog.Info("id = ", id)
	if id == "" {
		errorWriter(w, "id parameter is required")
		return
	}

	numId, err := strconv.Atoi(id)
	if err != nil {
		errorWriter(w, err.Error())
		return
	}

	rows, err := db.DeleteActivity(numId)
	if err != nil {
		slog.Error("Error deleting activity ...", err)
		errorWriter(w, "database error")
		return
	}

	slog.Info("deleted rows = ", rows)
	writeResponse(w, models.Envelope{"deleted rows": rows})
}
