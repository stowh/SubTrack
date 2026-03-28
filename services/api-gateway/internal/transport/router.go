package transport

import (
	"gateway/internal/config"
	"gateway/internal/transport/handlers"
	"gateway/internal/transport/middleware"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authorizationpb "github.com/stowh/subtrack/grpc/generate/authorization"
	subscribespb "github.com/stowh/subtrack/grpc/generate/subscribes"
)

func Register(engine *gin.Engine, authService authorizationpb.AuthorizationServiceClient, subService subscribespb.SubscribesServiceClient) {
	api := engine.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/acc/my", middleware.NewAuth(authService), handlers.HandleMyAcc(authService))
			auth.POST("/acc/register", handlers.HandleRegisterAcc(authService))
			auth.POST("/acc/login", handlers.HandleLoginAcc(authService))
			auth.POST("/session/refresh", handlers.HandleRefresh(authService))
			auth.POST("/session/logout", handlers.HandleLogout(authService))
		}

		subs := api.Group("/subs")
		subs.Use(middleware.NewAuth(authService))
		{
			subs.POST("/create", handlers.HandleCreateSub(subService))
			subs.DELETE("/remove", handlers.HandleRemoveSub(subService))
			subs.PUT("/changedata", handlers.HandleChangeDataSub(subService))
			subs.GET("/", handlers.HandleSearchSub(subService))
		}
	}
}

func Listen(engine *gin.Engine, conf *config.Config) {
	server := http.Server{
		Addr:         conf.GatewayAddr,
		Handler:      engine,
		WriteTimeout: time.Duration(conf.GatewayWriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(conf.GatewayReadTimeout) * time.Second,
	}

	log.Println("[+] server.listen:", conf.GatewayAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("[-] server.listen:", err)
	}
}
