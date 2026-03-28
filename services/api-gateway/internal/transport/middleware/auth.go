package middleware

import (
	"context"
	"gateway/internal/transport/models"
	"time"

	"github.com/gin-gonic/gin"
	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
)

func NewAuth(a authorizationpb.AuthorizationServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		access := ctx.GetHeader("Authorization")
		if access == "" {
			ctx.JSON(400, models.Response{IsOk: false, Message: "failed to get auth header"})
			ctx.Abort()
		}

		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()

		rspn, err := a.MyAccount(ctx2, &authorizationpb.MyAccountRequest{
			AccessToken: access,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "failed to get account"})
			return
		}

		ctx.Set("acc", rspn)
		ctx.Next()
	}
}
