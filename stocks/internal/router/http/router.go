package http

import (
	"net/http"
	"time"
)

type IStockController interface {
	AddStock(w http.ResponseWriter, r *http.Request)
	DeleteStockBySKU(w http.ResponseWriter, r *http.Request)
	GetItemsByLocation(w http.ResponseWriter, r *http.Request)
	GetItemBySKU(w http.ResponseWriter, r *http.Request)
}

type ServerConfig struct {
	Address           string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
}

func NewServer(serverConfig *ServerConfig) *http.Server {
	server := &http.Server{
		Addr:              serverConfig.Address,
		Handler:           serverConfig.Handler,
		ReadHeaderTimeout: serverConfig.ReadHeaderTimeout,
	}

	return server
}

func NewMux(c IStockController) *http.ServeMux {
	newMux := http.NewServeMux()

	newMux.HandleFunc("POST /stocks/item/add", c.AddStock)
	newMux.HandleFunc("POST /stocks/item/get", c.GetItemBySKU)
	newMux.HandleFunc("POST /stocks/item/delete", c.DeleteStockBySKU)
	newMux.HandleFunc("POST /stocks/list/location", c.GetItemsByLocation)

	return newMux
}
