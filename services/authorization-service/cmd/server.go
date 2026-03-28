package main

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/transport"
	"log"
)

func main() {
	conf, err := config.ParseEnv()
	if err != nil {
		log.Fatal("[-] config:", err)
	}

	postgres, err := repository.ConnectPostgres(conf)
	if err != nil {
		log.Fatal("[-] postgres.connect:", err)
	}

	err = postgres.Migration()
	if err != nil {
		log.Fatal("[-] postgres.migration:", err)
	}

	grpc := transport.CreateGRPC(conf, postgres)

	log.Println("[+] grpc.serve:", conf.ServerAddr)
	if err := transport.ListenGRPC(grpc, conf); err != nil {
		log.Fatal("[-] grpc.serve:", err)
	}
}
