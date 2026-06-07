package apiUtils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = os.Getenv("SECRET_KEY")

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func GenerateToken(email string, userId string) (*TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Minute * 15).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenStr := uuid.New().String()

	return &TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}, nil
}

func VerifyToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected errors")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	isValidToken := parsedToken.Valid

	if !isValidToken {
		return "nil", errors.New("Invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Invalid token")
	}

	userId := claims["userId"].(string)

	return userId, nil
}
