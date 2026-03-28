package handlers

import (
	"context"
	"gateway/internal/transport/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
	subscribespb "github.com/stowh/subtrack/grpc/generate/subscribes"
)

func HandleCreateSub(subService subscribespb.SubscribesServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.CreateSubRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		acc := ctx.MustGet("acc").(*authorizationpb.MyAccountResponse)

		ctx2, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*15)
		defer cancel()

		rspn, err := subService.CreateSub(ctx2, &subscribespb.CreateSubRequest{
			AccountId:               acc.GetAccountId(),
			SubscriptionName:        req.Name,
			SubscriptionTitle:       req.Title,
			SubscriptionPayPerMonth: req.PerMonth,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "Failed to create subscribe"})
			return
		}

		ctx.JSON(201, models.Response{IsOk: true, Message: "", Payload: map[string]any{"sub_id": rspn.GetSubscriptionId()}})
	}
}

func HandleRemoveSub(subService subscribespb.SubscribesServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.RemoveSubRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		acc := ctx.MustGet("acc").(*authorizationpb.MyAccountResponse)

		ctx2, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*15)
		defer cancel()

		_, err := subService.RemoveSub(ctx2, &subscribespb.RemoveSubRequest{
			AccountId:      acc.GetAccountId(),
			SubscriptionId: req.SubId,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "Failed to remove subscribe"})
			return
		}

		ctx.JSON(200, models.Response{IsOk: true, Message: ""})
	}
}

func HandleChangeDataSub(subService subscribespb.SubscribesServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req models.ChangeDataRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "invalid json"})
			return
		}

		acc := ctx.MustGet("acc").(*authorizationpb.MyAccountResponse)

		ctx2, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*15)
		defer cancel()

		_, err := subService.PutSub(ctx2, &subscribespb.PutSubRequest{
			AccountId:               acc.GetAccountId(),
			SubscriptionId:          req.SubId,
			SubscriptionName:        req.Name,
			SubscriptionTitle:       req.Title,
			SubscriptionPayPerMonth: req.PerMonth,
			SubscriptionStatus:      req.Status,
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "Failed to change data for subscribe"})
			return
		}

		ctx.JSON(200, models.Response{IsOk: true, Message: ""})
	}
}

func HandleSearchSub(subService subscribespb.SubscribesServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		acc := ctx.MustGet("acc").(*authorizationpb.MyAccountResponse)

		limitStr := ctx.Query("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "limit is not int"})
			return
		}
		if limit <= 0 {
			ctx.JSON(400, models.Response{IsOk: false, Message: "limit must be > 0"})
			return
		}

		ctx2, cancel := context.WithTimeout(ctx.Request.Context(), time.Second*15)
		defer cancel()

		rspn, err := subService.GetSubs(ctx2, &subscribespb.GetSubsRequest{
			AccountId: acc.GetAccountId(),
			Limit:     uint32(limit),
		})
		if err != nil {
			ctx.JSON(400, models.Response{IsOk: false, Message: "Failed to get subscribes"})
			return
		}

		var subs []*models.Subscribe
		for _, sub := range rspn.Items {
			subs = append(subs, &models.Subscribe{
				SubIt:       int(sub.SubscriptionId),
				SubName:     sub.SubscriptionName,
				SubTitle:    sub.SubscriptionTitle,
				SubPerMonth: int(sub.SubscriptionPayPerMonth),
				SubStatus:   int(sub.SubscriptionStatus),
				CreatedAt:   int(sub.CreatedAtUnix),
			})
		}

		ctx.JSON(200, models.Response{
			IsOk: true, Message: "",
			Payload: subs,
		})
	}
}
