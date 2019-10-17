package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID              int    `json:"id"`
	Email           string `json:"email"`
	Password        string `json:"password, omitempty"`
	EncryptPassword string `json:"-"`
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptPassword = enc
	}

	return nil
}

func encryptString(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func (u *User) Saintize() {
	u.Password = ""
}
