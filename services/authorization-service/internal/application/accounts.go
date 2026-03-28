package application

import (
	"auth/internal/config"
	"auth/internal/repository"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type (
	CreateAccReq struct {
		DisplayName string
		Email       string
		Password    string
	}
	LoginAccReq struct {
		Email    string
		Password string
	}

	Account struct {
		AccId int
		Name  string
		Email string
	}
	Tokens struct {
		AccId   int
		Access  string
		Refresh string
	}
)

func CreateAccount(postgres *repository.Postgres, conf *config.Config, req *CreateAccReq) (*Tokens, error) {
	if !CheckData(req.Email, req.Password) {
		return nil, fmt.Errorf("email or password is invalid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	accId, err := postgres.CreateAcc(req.DisplayName, req.Email, string(hash))
	if err != nil {
		return nil, err
	}

	refresh := GenerateRefresh()
	expiredRefresh := time.Now().Add(time.Duration(conf.RefreshExpiredDays) * 24 * time.Hour).Unix()
	if err := postgres.CreateSession(accId, refresh, expiredRefresh); err != nil {
		return nil, err
	}

	access, err := GenerateAccess(conf, accId, req.Email, req.DisplayName)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccId:   accId,
		Access:  access,
		Refresh: refresh,
	}, nil
}

func LoginAccount(postgres *repository.Postgres, conf *config.Config, req *LoginAccReq) (*Tokens, error) {
	if !CheckData(req.Email, req.Password) {
		return nil, fmt.Errorf("email or password is invalid")
	}

	acc, err := postgres.SearchAccByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Hash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	refresh := GenerateRefresh()
	expiredRefresh := time.Now().Add(time.Duration(conf.RefreshExpiredDays) * 24 * time.Hour).Unix()
	if err := postgres.CreateSession(acc.Id, refresh, expiredRefresh); err != nil {
		return nil, err
	}

	access, err := GenerateAccess(conf, acc.Id, acc.Email, acc.Name)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccId:   acc.Id,
		Access:  access,
		Refresh: refresh,
	}, nil
}

func MyAccount(conf *config.Config, access string) (*Account, error) {
	session, err := ParseAccess(conf, access)
	if err != nil {
		return nil, err
	}

	return &Account{
		AccId: session.AccId,
		Email: session.Email,
		Name:  session.DisplayName,
	}, nil
}
