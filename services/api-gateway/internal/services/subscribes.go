package services

import (
	"context"
	"time"

	subscribespb "github.com/stowh/subtrack/grpc/generate/subscribes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewSubscribes(addr string) (subscribespb.SubscribesServiceClient, error) {
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

	return subscribespb.NewSubscribesServiceClient(client), nil
}
