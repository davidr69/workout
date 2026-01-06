package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	_, err := fmt.Fprintf(w, "Exercises ...!")
	if err != nil {
		fmt.Println(err)
	}
}
