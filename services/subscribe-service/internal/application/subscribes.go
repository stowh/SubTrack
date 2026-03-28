package application

import (
	"fmt"
	"sub/internal/repository"
)

type (
	CreateSubReq struct {
		AccId    int
		Name     string
		Title    string
		PerMonth int
	}

	RemoveSubReq struct {
		AccId int
		SubId int
	}

	SearchSubReq struct {
		AccId int
		Limit uint
	}

	PutSubReq struct {
		AccId int
		SubId int

		SubName     string
		SubTitle    string
		SubPerMonth int
		Status      int
	}
)

func CreateSubscribe(postgres *repository.Postgres, req *CreateSubReq) (int, error) {
	if req.AccId == 0 || req.Name == "" {
		return 0, fmt.Errorf("invalid argument")
	}

	return postgres.CreateSub(req.AccId, req.Name, req.Title, req.PerMonth)
}

func RemoveSubscribe(postgres *repository.Postgres, req *RemoveSubReq) error {
	if req.AccId == 0 || req.SubId == 0 {
		return fmt.Errorf("invalid argument")
	}

	return postgres.RemoveSub(req.SubId, req.AccId)
}

func SearchSubscribe(postgres *repository.Postgres, req *SearchSubReq) ([]*repository.Sub, error) {
	if req.AccId == 0 {
		return nil, fmt.Errorf("invalid argument")
	}

	return postgres.SearchSubs(req.AccId, req.Limit)
}

func PutSubscribe(postgres *repository.Postgres, req *PutSubReq) error {
	if req.AccId == 0 || req.SubId == 0 {
		return fmt.Errorf("invalid argument")
	}

	if req.SubName != "" {
		if err := postgres.PutNameSub(req.SubId, req.AccId, req.SubName); err != nil {
			return err
		}
	}
	if req.SubTitle != "" {
		if err := postgres.PutTitleSub(req.SubId, req.AccId, req.SubTitle); err != nil {
			return err
		}
	}
	if req.SubPerMonth != 0 {
		if err := postgres.PutPerMonth(req.SubId, req.AccId, req.SubPerMonth); err != nil {
			return err
		}
	}
	if req.Status != 0 {
		if err := postgres.PutStatusSub(req.SubId, req.AccId, req.Status); err != nil {
			return err
		}
	}

	return nil
}
