package integration

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

func testStockService() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", getItemBySKU)

	return httptest.NewServer(mux)
}

func getItemBySKU(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	item := struct {
		Count uint16
	}{
		Count: 10,
	}

	res := struct {
		Message any `json:"message"`
	}{
		Message: item,
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}
