package entities

type RequestUser struct {
	ID    uint     `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}
