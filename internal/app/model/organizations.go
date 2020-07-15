package model

type Organizations struct {
	ID      int    `json:"id" db:"id"`
	OrgName string `json:"org_name" db:"org_name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
	Email   string `json:"email" db:"email"`
}
