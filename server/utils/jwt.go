package utils

import (
	"errors"
	"time"

	"go-admin/global"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 JWT Claims
type CustomClaims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
// 返回 accessToken 和 refreshToken
func GenerateToken(userID int64, username string) (accessToken string, refreshToken string, err error) {
	jwtCfg := global.Config.JWT

	// Access Token
	accessClaims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtCfg.Expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtCfg.Refresh) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString([]byte(jwtCfg.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseToken 解析 JWT Token，返回 Claims
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 Token
// 验证 refreshToken 有效性后，生成新的 accessToken 和 refreshToken
func RefreshToken(refreshTokenString string) (accessToken string, refreshToken string, err error) {
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return "", "", errors.New("refresh token 已失效，请重新登录")
	}

	// 生成新的 Token 对
	return GenerateToken(claims.UserID, claims.Username)
}
