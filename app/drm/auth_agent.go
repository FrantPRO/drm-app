package drm

import (
	"fmt"
	"strings"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type AuthAgent struct {
	users map[string]*User
}

func NewAuthAgent() *AuthAgent {
	return &AuthAgent{
		users: map[string]*User{
			"admin-token": {ID: "1", Name: "Admin", Role: "admin"},
			"user-token":  {ID: "2", Name: "User", Role: "user"},
			"guest-token": {ID: "3", Name: "Guest", Role: "guest"},
		},
	}
}

func (a *AuthAgent) ValidateToken(token string) (*User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	user, exists := a.users[token]
	if !exists {
		return nil, fmt.Errorf("invalid token")
	}

	return user, nil
}