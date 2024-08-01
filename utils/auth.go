package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"math/rand"
	"time"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GenerateAuthTokens(user *models.User) (types.Tokens, error) {
	accessTokenClaims := types.JWTClaims{
		Role:      user.Role,
		UserID:    user.ID,
		TokenType: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	refreshTokenClaims := types.JWTClaims{
		Role:      user.Role,
		UserID:    user.ID,
		TokenType: "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(configs.GetEnvVariable("JWT_SECRET")))

	if err != nil {
		return types.Tokens{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(configs.GetEnvVariable("JWT_SECRET")))

	if err != nil {
		return types.Tokens{}, err
	}

	return types.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func EncryptData(data []byte) ([]byte, error) {
	key := []byte(configs.GetEnvVariable("ENCRYPTION_KEY"))
	iv := []byte(configs.GetEnvVariable("INITIALIZATION_VECTOR"))

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	cipherData := make([]byte, len(data))
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherData, data)
	return cipherData, nil
}

func DecryptData(data []byte) ([]byte, error) {
	key := []byte(configs.GetEnvVariable("ENCRYPTION_KEY"))
	iv := []byte(configs.GetEnvVariable("INITIALIZATION_VECTOR"))

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBDecrypter(block, iv)
	plainData := make([]byte, len(data))
	cfb.XORKeyStream(plainData, data)
	return plainData, nil
}

func ParseToken(tokenString string) (*types.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.GetEnvVariable("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(*types.JWTClaims)

	if !ok {
		return nil, err
	}

	return claims, nil
}

func GeneratePassword(length int, forUser bool) string {
	var chars string
	if forUser {
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	} else {
		chars = "0123456789"
	}

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = chars[rand.Intn(len(chars))]
	}

	return string(randomString)
}
