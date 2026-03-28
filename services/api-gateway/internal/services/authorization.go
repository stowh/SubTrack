package services

import (
	"context"
	"time"

	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthorization(addr string) (authorizationpb.AuthorizationServiceClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return authorizationpb.NewAuthorizationServiceClient(client), nil
}
