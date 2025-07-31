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

type ServerConfig struct {
	Address             string
	Handler             http.Handler
	ReaderHeaderTimeout time.Duration
}

func NewGatewayServer(serverConfig *ServerConfig) *http.Server {
	server := &http.Server{
		Addr:              serverConfig.Address,
		Handler:           serverConfig.Handler,
		ReadHeaderTimeout: serverConfig.ReaderHeaderTimeout,
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
