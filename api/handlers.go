package api

import (
	"fmt"
	"net/http"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Healthy!")
	if err != nil {
		fmt.Println(err)
	}
}
