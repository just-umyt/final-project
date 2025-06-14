package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message any
	Code    int
}

func SuccesResponse(w http.ResponseWriter, msg any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := response{
		Message: msg,
		Code:    http.StatusOK,
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}
