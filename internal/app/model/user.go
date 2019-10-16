package model

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password, omitempty"`
}

func (u *User) Saintize() {
	u.Password = ""
}
