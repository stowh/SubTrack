package main

import (
	"context"
	"log"
	"os/signal"
	"sub/internal/config"
	"sub/internal/repository"
	"sub/internal/transport"
	"syscall"
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
	log.Println("[+] postgres.connect:", conf.PostgresAddr)

	if err := postgres.Migration(); err != nil {
		log.Fatal("[-] postgres.migration:", err)
	}
	log.Println("[+] postgres.migration")
	defer func() {
		if err := postgres.Close(); err != nil {
			log.Println("[-] postgres.close:", err)
		}
	}()

	grpc := transport.CreateGRPC(conf, postgres)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("[+] grpc.serve:", conf.ServerAddr)
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- transport.ListenGRPC(grpc, conf)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			log.Fatal("[-] grpc.serve:", err)
		}
	case <-ctx.Done():
		log.Println("[*] grpc.shutdown")
		grpc.GracefulStop()
		if err := <-serveErr; err != nil {
			log.Fatal("[-] grpc.shutdown:", err)
		}
	}
}
