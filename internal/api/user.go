package api

import (
	"errors"
	"net/http"
)

type User struct {
	UID               string
	Name              string
	EmailAddress      string
	EncryptedPassword string
}

// Retrieved from JWT
type AuthUser struct {
	UID  string
	Role string
}

type RegisterUserRequest struct {
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

type UserResponse struct {
	UID          string `json:"uid"`
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
}

type LoginRequest struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}

func (lr *LoginRequest) Bind(r *http.Request) error {
	if lr.EmailAddress == "" {
		return errors.New("empty email_address")
	}
	if lr.Password == "" {
		return errors.New("empty password")
	}
	return nil
}
