package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Err string `json:"error"`
}

func ErrorResponse(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := errorResponse{
		Err: err.Error(),
	}

	log.Println(err)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}
