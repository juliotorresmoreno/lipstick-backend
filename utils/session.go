package utils

import (
	"context"
	"time"

	"github.com/juliotorresmoreno/tana-api/cache"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/models"
)

var SessionFields = []string{"id", "name", "last_name", "email", "photo_url", "phone"}

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	PhotoURL string `json:"photo_url"`
	Phone    string `json:"phone"`
}

type Session struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func ValidateSession(token string) (*Session, error) {
	ctx := context.Background()
	cmd := cache.DefaultClient.Get(ctx, "session-"+token)
	email := cmd.Val()
	if email == "" {
		return &Session{}, StatusUnauthorized
	}

	conn := db.DefaultClient
	user := &models.User{}
	tx := conn.Select(SessionFields).First(user, "email = ? AND deleted_at IS NULL", email)
	if tx.Error != nil {
		return &Session{}, StatusInternalServerError
	}

	cache.DefaultClient.Set(ctx, "session-"+token, email, 24*time.Hour)

	return ParseSession(token, user), nil
}

func ParseSession(token string, user *models.User) *Session {
	return &Session{
		Token: token,
		User: &User{
			ID:       user.ID,
			Name:     user.Name,
			LastName: user.LastName,
			Email:    user.Email,
			PhotoURL: user.PhotoURL,
			Phone:    user.Phone,
		},
	}
}

func MakeSession(user *models.User) (*Session, error) {
	token, err := GenerateRandomString(128)
	if err != nil {
		return &Session{}, StatusInternalServerError
	}

	cmd := cache.DefaultClient.Set(
		context.Background(),
		"session-"+token,
		user.Email, 24*time.Hour,
	)
	if cmd.Err() != nil {
		return &Session{}, StatusInternalServerError
	}

	return ParseSession(token, user), nil
}
