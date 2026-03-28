package models

// Account
type (
	RegisterAccRequest struct {
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
	}

	LoginAccRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AccessAccRequest struct {
		Access string `json:"access"`
	}

	RefreshAccRequest struct {
		Refresh string `json:"refresh"`
	}
)

// Subscrubes
type (
	CreateSubRequest struct {
		Name     string `json:"sub_name"`
		Title    string `json:"sub_title"`
		PerMonth int64  `json:"sub_pay_per_month"`
	}

	RemoveSubRequest struct {
		SubId int64 `json:"sub_id"`
	}

	ChangeDataRequest struct {
		SubId    int64  `json:"sub_id"`
		Name     string `json:"sub_name"`
		Title    string `json:"sub_title"`
		PerMonth int64  `json:"sub_pay_per_month"`
		Status   int64  `json:"sub_status"`
	}
)
