package application

import (
	"auth/internal/config"
	"auth/internal/repository"
	"fmt"
	"time"
)

func Refresh(postgres *repository.Postgres, conf *config.Config, refresh string) (*Tokens, error) {
	session, err := postgres.SearchSession(refresh)
	if err != nil {
		return nil, err
	}

	postgres.RemoveSession(refresh)
	if session.Expired < time.Now().Unix() {
		return nil, fmt.Errorf("refresh expired")
	}

	refresh = GenerateRefresh()
	expiredRefresh := time.Now().Add(time.Duration(conf.RefreshExpiredDays) * 24 * time.Hour).Unix()
	if err := postgres.CreateSession(session.AccId, refresh, expiredRefresh); err != nil {
		return nil, err
	}

	acc, err := postgres.SearchAccById(session.AccId)
	if err != nil {
		return nil, err
	}

	access, err := GenerateAccess(conf, acc.Id, acc.Email, acc.Name)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccId:   session.AccId,
		Access:  access,
		Refresh: refresh,
	}, nil
}

func Logout(postgres *repository.Postgres, refresh string) error {
	return postgres.RemoveSession(refresh)
}
