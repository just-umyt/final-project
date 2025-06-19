package utils

import (
	"encoding/json"
	"net/http"
	"stocks/pkg/logger"
)

type errorResponse struct {
	Err string `json:"error"`
}

func Error(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := errorResponse{
		Err: err.Error(),
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		logger.Log.Error(err)
	}
}
