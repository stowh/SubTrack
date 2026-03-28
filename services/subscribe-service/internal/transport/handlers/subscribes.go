package handlers

import (
	"context"
	"sub/internal/application"
	"sub/internal/config"
	"sub/internal/repository"

	subscribespb "github.com/stowh/subtrack/grpc/generate/subscribes"
)

type (
	SubServer struct {
		conf     *config.Config
		postgres *repository.Postgres
		subscribespb.UnimplementedSubscribesServiceServer
	}
)

func NewService(conf *config.Config, postgres *repository.Postgres) *SubServer {
	return &SubServer{
		conf:     conf,
		postgres: postgres,
	}
}

func (s *SubServer) CreateSub(ctx context.Context, in *subscribespb.CreateSubRequest) (*subscribespb.CreateSubResponse, error) {
	lastId, err := application.CreateSubscribe(s.postgres, &application.CreateSubReq{
		AccId:    int(in.GetAccountId()),
		Name:     in.GetSubscriptionName(),
		Title:    in.GetSubscriptionTitle(),
		PerMonth: int(in.GetSubscriptionPayPerMonth()),
	})
	if err != nil {
		return nil, err
	}

	return &subscribespb.CreateSubResponse{SubscriptionId: int64(lastId)}, nil
}

func (s *SubServer) RemoveSub(ctx context.Context, in *subscribespb.RemoveSubRequest) (*subscribespb.RemoveSubResponse, error) {
	err := application.RemoveSubscribe(s.postgres, &application.RemoveSubReq{
		AccId: int(in.GetAccountId()),
		SubId: int(in.GetSubscriptionId()),
	})
	if err != nil {
		return nil, err
	}

	return &subscribespb.RemoveSubResponse{Success: true}, nil
}

func (s *SubServer) PutSub(ctx context.Context, in *subscribespb.PutSubRequest) (*subscribespb.PutSubResponse, error) {
	err := application.PutSubscribe(s.postgres, &application.PutSubReq{
		AccId:       int(in.GetAccountId()),
		SubId:       int(in.GetSubscriptionId()),
		SubName:     in.GetSubscriptionName(),
		SubTitle:    in.GetSubscriptionTitle(),
		SubPerMonth: int(in.GetSubscriptionPayPerMonth()),
		Status:      int(in.GetSubscriptionStatus()),
	})
	if err != nil {
		return nil, err
	}

	return &subscribespb.PutSubResponse{Success: true}, nil
}

func (s *SubServer) GetSubs(ctx context.Context, in *subscribespb.GetSubsRequest) (*subscribespb.GetSubsResponse, error) {
	subs, err := application.SearchSubscribe(s.postgres, &application.SearchSubReq{
		AccId: int(in.GetAccountId()),
		Limit: uint(in.GetLimit()),
	})
	if err != nil {
		return nil, err
	}

	var rspn subscribespb.GetSubsResponse
	for _, sub := range subs {
		rspn.Items = append(rspn.Items, &subscribespb.Subscription{
			SubscriptionId:          int64(sub.Id),
			SubscriptionName:        sub.Name,
			SubscriptionTitle:       sub.Title,
			SubscriptionPayPerMonth: int64(sub.PayPerMonth),
			SubscriptionStatus:      int64(sub.Status),
			CreatedAtUnix:           int64(sub.CreatedAt),
		})
	}

	return &rspn, nil
}
