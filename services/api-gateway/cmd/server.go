package main

import (
	"gateway/internal/config"
	"gateway/internal/services"
	"gateway/internal/transport"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.ParseEnv()
	if err != nil {
		log.Fatalln("[-] config:", err)
	}

	authClient, err := services.NewAuthorization(conf.GatewayAuthAddr)
	if err != nil {
		log.Fatalln("[-] services.auth:", err)
	}
	log.Println("[+] services: new auth")

	subClient, err := services.NewSubscribes(conf.GatewaySubAddr)
	if err != nil {
		log.Fatalln("[-] services.sub:", err)
	}
	log.Println("[+] services: new sub")

	engine := gin.Default()
	transport.Register(engine, authClient, subClient)
	transport.Listen(engine, conf)
}
