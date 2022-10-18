package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const accessTokenTypeName = "access_token"
const refreshTokenTypeName = "refresh_token"

var JWTKey []byte

type JWTClaim struct {
	jwt.RegisteredClaims
	UserID    uint   `json:"userID"`
	Username  string `json:"username"`
	TokenType string `json:"tokenType"`
}

func GenerateAccessToken(userID uint, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		UserID:    userID,
		Username:  username,
		TokenType: accessTokenTypeName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JWTKey)
	return
}

func GenerateRefreshToken(userID uint, username string) (tokenString string, err error) {
	claims := &JWTClaim{
		UserID:    userID,
		Username:  username,
		TokenType: refreshTokenTypeName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JWTKey)
	return
}

func GenerateTokenPair(userID uint, username string) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(userID, username)
	if err != nil {
		return
	}
	refreshToken, err = GenerateRefreshToken(userID, username)
	return
}

func ValidateAccessToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.TokenType != accessTokenTypeName {
		err = errors.New("invalid access token")
		return
	}
	if claims.ExpiresAt.Before(time.Now().Local()) {
		err = errors.New("token expired")
		return
	}
	return
}

func ValidateRefreshToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return JWTKey, nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.TokenType != refreshTokenTypeName {
		err = errors.New("invalid refresh token")
		return
	}
	return
}
