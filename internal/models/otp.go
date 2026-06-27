package models

import "time"

type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type VerifyReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type WaitingUser struct {
	User      User
	Otp       string
	CreatedAt time.Time
}
