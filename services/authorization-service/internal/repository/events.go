package repository

import "time"

type (
	Account struct {
		Id    int
		Name  string
		Email string
		Hash  string
	}

	Session struct {
		AccId   int
		Token   string
		Expired int64
	}
)

func (p *Postgres) CreateAcc(displayName string, email string, password string) (int, error) {
	var id int
	err := p.db.QueryRow("INSERT INTO acc (name, email, hash, created_at) VALUES($1, $2, $3, $4) RETURNING id",
		displayName, email, password, time.Now().Unix()).Scan(&id)
	return id, err
}

func (p *Postgres) SearchAccByEmail(email string) (*Account, error) {
	var id int
	var name, hash string

	err := p.db.QueryRow("SELECT id, name, hash FROM acc WHERE email=$1", email).Scan(&id, &name, &hash)
	if err != nil {
		return nil, err
	}

	return &Account{Id: id, Name: name, Email: email, Hash: hash}, nil
}

func (p *Postgres) SearchAccById(id int) (*Account, error) {
	var name, email, hash string

	err := p.db.QueryRow("SELECT email, name, hash FROM acc WHERE id=$1", id).Scan(&email, &name, &hash)
	if err != nil {
		return nil, err
	}

	return &Account{Id: id, Name: name, Email: email, Hash: hash}, nil
}

func (p *Postgres) CreateSession(accId int, token string, expired int64) error {
	_, err := p.db.Exec("INSERT INTO session (acc_id, token, expired_at, created_at) VALUES ($1, $2, $3, $4)",
		accId, token, expired, time.Now().Unix())
	return err
}

func (p *Postgres) SearchSession(token string) (*Session, error) {
	var accId int
	var expired int64
	err := p.db.QueryRow("SELECT acc_id, expired_at FROM session WHERE token=$1", token).Scan(&accId, &expired)
	if err != nil {
		return nil, err
	}

	return &Session{
		AccId:   accId,
		Token:   token,
		Expired: expired,
	}, nil
}

func (p *Postgres) RemoveSession(token string) error {
	_, err := p.db.Exec("DELETE FROM session WHERE token=$1", token)
	return err
}
