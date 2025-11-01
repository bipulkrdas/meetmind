package utils

import (
    "time"

    "github.com/dgrijalva/jwt-go"
)

func GenerateJWT(userID, email, secret string) (string, time.Time, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &jwt.StandardClaims{
        Subject:   userID,
        ExpiresAt: expirationTime.Unix(),
        Issuer:    email,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secret))

    return tokenString, expirationTime, err
}
