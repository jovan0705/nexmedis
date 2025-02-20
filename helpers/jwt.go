package helpers

import (
    "time"
    "github.com/dgrijalva/jwt-go"
    "github.com/jovan0705/nexmedis/models"
)

var secretKey = []byte("your_secret_key")

func GenerateJWT(user models.User) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })
    return token.SignedString(secretKey)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	return token, err
}