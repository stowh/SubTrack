package repository

import (
	"analytics/internal/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type (
	Postgres struct {
		db *sql.DB
	}
)

func ConnectPostgres(conf *config.Config) (*Postgres, error) {
	str := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		conf.PostgresUser, conf.PostgresPass, conf.PostgresAddr, conf.PostgresName)

	db, err := sql.Open("pgx", str)
	if err != nil {
		return nil, err
	}

	for range 10 {
		if err = db.Ping(); err != nil {
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) Migration() error {
	_, err := p.db.Exec(`CREATE TABLE IF NOT EXISTS acc (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		hash TEXT NOT NULL,
		created_at INTEGER
	)`)

	return err
}
