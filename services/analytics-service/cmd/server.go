package main

import (
	"analytics/internal/config"
	"analytics/internal/repository"
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

}
