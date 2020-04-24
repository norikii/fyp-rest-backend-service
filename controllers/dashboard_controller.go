package controllers

import (
	"encoding/json"
	"net/http"
)

func Dashboard(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response)
}
