package transport

import (
	"net"
	"sub/internal/config"
	"sub/internal/repository"
	"sub/internal/transport/handlers"

	subscribespb "github.com/stowh/subtrack/grpc/generate/subscribes"
	"google.golang.org/grpc"
)

func CreateGRPC(conf *config.Config, postgres *repository.Postgres) *grpc.Server {
	server := grpc.NewServer()
	service := handlers.NewService(conf, postgres)
	subscribespb.RegisterSubscribesServiceServer(server, service)

	return server
}

func ListenGRPC(server *grpc.Server, conf *config.Config) error {
	listen, err := net.Listen("tcp", conf.ServerAddr)
	if err != nil {
		return err
	}

	return server.Serve(listen)
}
