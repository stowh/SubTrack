package repository

import (
	"database/sql"
	"fmt"
	"sub/internal/config"
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
	_, err := p.db.Exec(`CREATE TABLE IF NOT EXISTS subs (
		id SERIAL PRIMARY KEY,
		acc_id INTEGER NOT NULL,
		subName TEXT NOT NULL,
		subTitle TEXT,
		subStatus INTEGER,
		subPay_per_month INTEGER NOT NULL,
		created_at INTEGER
	)`)

	return err
}
