package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type api struct {
	//db *database.APIServerDB
}

type DataProduct struct {
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	Type        string `firestore:"type"`
	URI         string `firestore:"uri"`
}

func (a *api) dataproducts(w http.ResponseWriter, r *http.Request) {
	dataproducts := []DataProduct{{Name: "my_dp", Description: "description of dp", Type: "bigquery", URI: "https://google.com/dp/x"}}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(dataproducts)

	if err != nil {
		log.Errorf("encoding dataproducts response: %v", err)
		respondf(w, http.StatusInternalServerError, "unable to get device config\n")
		return
	}
}

func respondf(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)

	if _, wErr := w.Write([]byte(fmt.Sprintf(format, args...))); wErr != nil {
		log.Errorf("unable to write response: %v", wErr)
	}
}
