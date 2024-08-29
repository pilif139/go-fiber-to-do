package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(userId int) (string, error) {
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":    time.Now().Unix(),
	})

	token, err := claims.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid authorization token")
	}

	return token, nil
}

func GetUserId(token *jwt.Token) (int, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userId, ok := claims["userId"].(float64) //float64 because json formats all numbers to float64
	if !ok {
		return 0, fmt.Errorf("invalid userId")
	}

	return int(userId), nil
}
