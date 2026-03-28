package handlers

import (
	"context"
	"fmt"
	"gateway/internal/transport/models"
	"time"

	"github.com/gin-gonic/gin"
	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
)

func HandleRegisterAcc(auth authorizationpb.AuthorizationServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.RegisterAccRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		rspn, err := auth.RegisterAccount(ctx2, &authorizationpb.RegisterAccountRequest{
			DisplayName: req.DisplayName,
			Email:       req.Email,
			Password:    req.Password,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "failed to register account"})
			return
		}

		ctx.JSON(201, models.Response{
			IsOk:    true,
			Message: "",
			Payload: models.Tokens{
				Access: rspn.AccessToken,
				Refrsh: rspn.RefreshToken,
			},
		})
	}
}

func HandleLoginAcc(auth authorizationpb.AuthorizationServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.LoginAccRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		rspn, err := auth.LoginAccount(ctx2, &authorizationpb.LoginAccountRequest{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "failed to login"})
			return
		}

		ctx.JSON(200, models.Response{
			IsOk:    true,
			Message: "",
			Payload: models.Tokens{
				Access: rspn.AccessToken,
				Refrsh: rspn.RefreshToken,
			},
		})
	}
}

func HandleRefresh(auth authorizationpb.AuthorizationServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.RefreshAccRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		rspn, err := auth.Refresh(ctx2, &authorizationpb.RefreshRequest{
			RefreshToken: req.Refresh,
		})
		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, models.Response{IsOk: false, Message: "failed to refresh session"})
			return
		}

		ctx.JSON(200, models.Response{
			IsOk:    true,
			Message: "",
			Payload: models.Tokens{
				Access: rspn.AccessToken,
				Refrsh: rspn.RefreshToken,
			},
		})
	}
}

func HandleLogout(auth authorizationpb.AuthorizationServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.RefreshAccRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		_, err := auth.Logout(ctx2, &authorizationpb.LogoutRequest{
			RefreshToken: req.Refresh,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "failed to logout"})
			return
		}

		ctx.JSON(200, models.Response{IsOk: true, Message: ""})
	}
}

func HandleMyAcc(auth authorizationpb.AuthorizationServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		acc := ctx.MustGet("acc").(*authorizationpb.MyAccountResponse)

		ctx.JSON(200, models.Response{
			IsOk:    true,
			Message: "",
			Payload: models.Account{
				AccId:       int(acc.AccountId),
				DisplayName: acc.DisplayName,
				Email:       acc.Email,
			},
		})
	}
}
