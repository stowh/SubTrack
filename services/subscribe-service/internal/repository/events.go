package repository

import "time"

type (
	Sub struct {
		Id          int
		AccId       int
		Name        string
		Title       string
		Status      int
		PayPerMonth int
		CreatedAt   int
	}
)

func (p *Postgres) CreateSub(accId int, name string, title string, payPerMonth int) (int, error) {
	var lastId int
	err := p.db.QueryRow(
		"INSERT INTO subs (acc_id, subName, subTitle, subPay_per_month, created_at, subStatus) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		accId, name, title, payPerMonth, time.Now().Unix(), 0,
	).Scan(&lastId)
	return lastId, err
}

func (p *Postgres) RemoveSub(id int, accId int) error {
	_, err := p.db.Exec("DELETE FROM subs WHERE id=$1 AND acc_id=$2", id, accId)
	return err
}

func (p *Postgres) SearchSubs(accId int, limit uint) ([]*Sub, error) {
	if limit == 0 {
		limit = 25
	}
	if limit > 200 {
		limit = 200
	}

	rows, err := p.db.Query(
		"SELECT id, acc_id, subName, subTitle, subStatus, subPay_per_month, created_at FROM subs WHERE acc_id=$1 ORDER BY created_at DESC LIMIT $2",
		accId, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]*Sub, 0, limit)
	for rows.Next() {
		sub := &Sub{}
		if err = rows.Scan(
			&sub.Id,
			&sub.AccId,
			&sub.Name,
			&sub.Title,
			&sub.Status,
			&sub.PayPerMonth,
			&sub.CreatedAt,
		); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subs, nil
}

func (p *Postgres) PutNameSub(id int, accId int, name string) error {
	_, err := p.db.Exec("UPDATE subs SET subName=$1 WHERE id=$2 AND acc_id=$3", name, id, accId)
	return err
}

func (p *Postgres) PutTitleSub(id int, accId int, title string) error {
	_, err := p.db.Exec("UPDATE subs SET subTitle=$1 WHERE id=$2 AND acc_id=$3", title, id, accId)
	return err
}

func (p *Postgres) PutStatusSub(id int, accId int, status int) error {
	_, err := p.db.Exec("UPDATE subs SET subStatus=$1 WHERE id=$2 AND acc_id=$3", status, id, accId)
	return err
}

func (p *Postgres) PutPerMonth(id int, accId int, pay int) error {
	_, err := p.db.Exec("UPDATE subs SET subPay_per_month=$1 WHERE id=$2 AND acc_id=$3", pay, id, accId)
	return err
}
