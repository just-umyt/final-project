package utils

import (
	"encoding/json"
	"net/http"
	"stocks/pkg/logger"
)

type response struct {
	Message any
	Code    int
}

func SuccessResponse(w http.ResponseWriter, msg any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := response{
		Message: msg,
		Code:    http.StatusOK,
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err)
	}
}
