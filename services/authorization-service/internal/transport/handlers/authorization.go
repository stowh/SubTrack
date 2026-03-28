package handlers

import (
	"auth/internal/application"
	"auth/internal/config"
	"auth/internal/repository"
	"context"

	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
)

type (
	AuthServer struct {
		conf     *config.Config
		postgres *repository.Postgres
		authorizationpb.UnimplementedAuthorizationServiceServer
	}
)

func NewService(conf *config.Config, postgres *repository.Postgres) *AuthServer {
	return &AuthServer{
		conf:     conf,
		postgres: postgres,
	}
}

func (a *AuthServer) RegisterAccount(ctx context.Context, in *authorizationpb.RegisterAccountRequest) (*authorizationpb.RegisterAccountResponse, error) {
	tokens, err := application.CreateAccount(a.postgres, a.conf, &application.CreateAccReq{
		DisplayName: in.GetDisplayName(),
		Email:       in.GetEmail(),
		Password:    in.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &authorizationpb.RegisterAccountResponse{
		AccountId:    int64(tokens.AccId),
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}, nil
}

func (a *AuthServer) LoginAccount(ctx context.Context, in *authorizationpb.LoginAccountRequest) (*authorizationpb.LoginAccountResponse, error) {
	tokens, err := application.LoginAccount(a.postgres, a.conf, &application.LoginAccReq{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &authorizationpb.LoginAccountResponse{
		AccountId:    int64(tokens.AccId),
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}, nil

}

func (a *AuthServer) MyAccount(ctx context.Context, in *authorizationpb.MyAccountRequest) (*authorizationpb.MyAccountResponse, error) {
	acc, err := application.MyAccount(a.conf, in.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return &authorizationpb.MyAccountResponse{
		AccountId:   int64(acc.AccId),
		DisplayName: acc.Name,
		Email:       acc.Email,
	}, nil
}

func (a *AuthServer) Refresh(ctx context.Context, in *authorizationpb.RefreshRequest) (*authorizationpb.RefreshResponse, error) {
	tokens, err := application.Refresh(a.postgres, a.conf, in.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authorizationpb.RefreshResponse{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}, nil
}

func (a *AuthServer) Logout(ctx context.Context, in *authorizationpb.LogoutRequest) (*authorizationpb.LogoutResponse, error) {
	err := application.Logout(a.postgres, in.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authorizationpb.LogoutResponse{Success: true}, nil
}
