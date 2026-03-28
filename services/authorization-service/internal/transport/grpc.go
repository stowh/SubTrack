package transport

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/transport/handlers"
	"net"

	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
	"google.golang.org/grpc"
)

func CreateGRPC(conf *config.Config, postgres *repository.Postgres) *grpc.Server {
	server := grpc.NewServer()
	service := handlers.NewService(conf, postgres)
	authorizationpb.RegisterAuthorizationServiceServer(server, service)

	return server
}

func ListenGRPC(server *grpc.Server, conf *config.Config) error {
	listen, err := net.Listen("tcp", conf.ServerAddr)
	if err != nil {
		return err
	}

	return server.Serve(listen)
}
