package grpc

import (
	"context"
	"net/http"
	"time"

	pb "cart/pkg/api/cart"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	readHeaderTimeout = 3
)

type ServerConfig struct {
	Address string
	Handler http.Handler
}

func NewGatewayServer(serverConfig *ServerConfig) *http.Server {
	server := &http.Server{
		Addr:              serverConfig.Address,
		Handler:           serverConfig.Handler,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	return server
}

func NewMux(ctx context.Context, grpcAddress string) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterCartServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return nil, err
	}

	return mux, nil
}
