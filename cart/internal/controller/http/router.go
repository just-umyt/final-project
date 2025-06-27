package http

import (
	"net/http"
	"time"
)

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

func NewMux(c ICartController) *http.ServeMux {
	newMux := http.NewServeMux()

	newMux.HandleFunc("POST /cart/item/add", c.AddItem)
	newMux.HandleFunc("POST /cart/item/delete", c.DeleteItem)
	newMux.HandleFunc("POST /cart/list", c.CartList)
	newMux.HandleFunc("POST /cart/clear", c.CartClear)

	return newMux
}
