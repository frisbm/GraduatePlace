package user

import "time"

type RegisterUser struct {
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	Username  string `json:"username,omitempty"`
}

type LoginUser struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthTokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshUser struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

type GetUser struct {
	UUID      string
	FirstName string
	LastName  string
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
