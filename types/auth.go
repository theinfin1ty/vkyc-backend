package types

import "github.com/golang-jwt/jwt"

type JWTClaims struct {
	Role      string `json:"role"`
	UserID    uint   `json:"userId"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
