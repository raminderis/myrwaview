package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/raminderis/myrwaview/rand"
)

const (
	MinBytesPerToken = 32
)

type Session struct {
	ID        int
	UserID    int
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *sql.DB
	BytesPerToken int
}

func (service *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (service *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := service.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: service.hash(token),
	}
	row := service.DB.QueryRow(`
		UPDATE sessions
		SET token_hash = $2
		WHERE user_id = $1
		RETURNING id`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row := service.DB.QueryRow(`
			INSERT INTO sessions (user_id, token_hash)
			VALUES ($1, $2) RETURNING id`, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (service *SessionService) User(token string) (*User, error) {
	tokenHash := service.hash(token)
	var userID int
	row := service.DB.QueryRow(`
		SELECT user_id
		FROM sessions
		WHERE token_hash = $1`, tokenHash)
	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	var user User
	row = service.DB.QueryRow(`
		SELECT email, password_hash
		FROM users
		WHERE id = $1`, userID)
	err = row.Scan(&user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := service.hash(token)
	_, err := service.DB.Exec(`
		DELETE FROM sessions
		WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
