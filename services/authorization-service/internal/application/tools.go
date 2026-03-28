package application

import (
	"auth/internal/config"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	refresh_symbols = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789"
	refresh_len     = 128
)

type (
	JwtClaims struct {
		AccId       int
		Email       string
		DisplayName string
	}
)

func CheckData(email string, password string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if len(password) < 6 || len(password) > 64 {
		return false
	}

	return true
}

func GenerateRefresh() string {
	bytes := make([]byte, refresh_len)
	for in, _ := range bytes {
		number, _ := rand.Int(rand.Reader, big.NewInt(int64(len(refresh_symbols))))
		bytes[in] = refresh_symbols[number.Int64()]
	}
	return string(bytes)
}

func GenerateAccess(conf *config.Config, accId int, email string, name string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Duration(conf.JwtExpiredMin) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"acc_id":       accId,
		"email":        email,
		"display_name": name,
		"iat":          now.Unix(),
		"exp":          exp.Unix(),
	})

	access, err := token.SignedString([]byte(conf.JwtSecret))
	if err != nil {
		return "", err
	}

	return access, nil
}

func ParseAccess(conf *config.Config, access string) (*JwtClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	accIDRaw, ok := claims["acc_id"]
	if !ok {
		return nil, fmt.Errorf("missing acc_id in token")
	}
	accIDFloat, ok := accIDRaw.(float64)
	if !ok {
		return nil, fmt.Errorf("invalid acc_id in token")
	}

	displayNameRaw, ok := claims["display_name"]
	if !ok {
		return nil, fmt.Errorf("missing display_name in token")
	}
	displayName, ok := displayNameRaw.(string)
	if !ok {
		return nil, fmt.Errorf("invalid display_name in token")
	}

	emailRaw, ok := claims["email"]
	if !ok {
		return nil, fmt.Errorf("missing email in token")
	}
	email, ok := emailRaw.(string)
	if !ok {
		return nil, fmt.Errorf("invalid email in token")
	}

	return &JwtClaims{
		AccId:       int(accIDFloat),
		Email:       email,
		DisplayName: displayName,
	}, nil
}
