package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccesResp struct {
	Message any `json:"message"`
}

func SuccessResponse(w http.ResponseWriter, msg any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := SuccesResp{
		Message: msg,
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}
