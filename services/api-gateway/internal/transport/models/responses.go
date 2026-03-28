package models

type (
	Response struct {
		IsOk    bool   `json:"is_ok"`
		Message string `json:"msg"`
		Payload any    `json:"payload"`
	}

	Tokens struct {
		Access string `json:"access"`
		Refrsh string `json:"refresh"`
	}

	Account struct {
		AccId       int    `json:"id"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
	Subscribe struct {
		SubIt       int    `json:"sub_id"`
		SubName     string `json:"sub_name"`
		SubTitle    string `json:"sub_title"`
		SubPerMonth int    `json:"sub_pay_per_month"`
		SubStatus   int    `json:"sub_status"`
		CreatedAt   int    `json:"created_at"`
	}
)
