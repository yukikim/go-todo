package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"go-todo/internal/model"
)

var ErrInvalidCredentials = errors.New("invalid username or password")
var ErrInvalidToken = errors.New("invalid token")

type AuthService struct {
	username string
	password string
	secret   []byte
	now      func() time.Time
}

type tokenClaims struct {
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
}

func NewAuthService(username, password, secret string) *AuthService {
	return &AuthService{
		username: username,
		password: password,
		secret:   []byte(secret),
		now:      time.Now,
	}
}

func (s *AuthService) Login(req model.LoginRequest) (model.LoginResponse, error) {
	if req.Username != s.username || req.Password != s.password {
		return model.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(req.Username)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return model.LoginResponse{Token: token}, nil
}

func (s *AuthService) GenerateToken(username string) (string, error) {
	now := s.now()
	claims := tokenClaims{
		Subject:   username,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(24 * time.Hour).Unix(),
	}

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	header := base64.RawURLEncoding.EncodeToString(headerJSON)
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := header + "." + payload
	signature := s.sign(signingInput)

	return signingInput + "." + signature, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := s.sign(signingInput)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return "", ErrInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalidToken
	}

	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return "", ErrInvalidToken
	}
	if header.Algorithm != "HS256" || header.Type != "JWT" {
		return "", ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidToken
	}

	var claims tokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return "", ErrInvalidToken
	}
	if claims.Subject == "" || claims.ExpiresAt <= s.now().Unix() {
		return "", ErrInvalidToken
	}

	return claims.Subject, nil
}

func (s *AuthService) sign(signingInput string) string {
	mac := hmac.New(sha256.New, s.secret)
	fmt.Fprint(mac, signingInput)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
