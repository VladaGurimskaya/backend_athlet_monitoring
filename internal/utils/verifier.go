package utils

import (
	"backend_athlet_monitoring/internal/models"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func ExtractHTTPRequest(ctx context.Context) (*http.Request, error) {
	httpReq, ok := ctx.Value("http-request").(*http.Request)
	if !ok {
		return nil, fmt.Errorf("HTTP-запрос не найден в контексте")
	}
	return httpReq, nil
}

func CreateJWTTokens(user models.User) (string, string, error) {
	accessToken, err := createAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := createRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func createAccessToken(user models.User) (string, error) {
	expiresAt, err := time.ParseDuration("15m")
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга времени: %w", err)
	}

	claims := &models.AccessToken{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("ошибка создания JWT токена: %w", err)
	}

	return tokenString, nil
}

func createRefreshToken(user models.User) (string, error) {
	expiresAt, err := time.ParseDuration("24h")
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга времени: %w", err)
	}

	claims := &models.RefreshToken{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fmt.Errorf("ошибка создания JWT токена: %w", err)
	}

	return tokenString, nil
}

func VerifyAccessToken(httpReq *http.Request) (*models.AccessToken, error) {
	accessCookie, err := httpReq.Cookie("access_token")
	if err != nil {
		return nil, fmt.Errorf("access token отсутствует: %v", err)
	}

	claims, err := verifyAccessToken(accessCookie.Value)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки JWT токена: %v", err)
	}

	return claims, nil
}

func verifyAccessToken(tokenString string) (*models.AccessToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AccessToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неправильный алгоритм: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JWT токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("неверный JWT токен")
	}

	claims, ok := token.Claims.(*models.AccessToken)
	if !ok {
		return nil, fmt.Errorf("ошибка парсинга JWT токена: %w", err)
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("JWT токен истек")
	}

	return claims, nil
}

func VerifyRefreshToken(httpReq *http.Request) (*models.RefreshToken, error) {
	refreshToken, err := httpReq.Cookie("refresh_token")
	if err != nil {
		return nil, fmt.Errorf("refresh token отсутствует: %v", err)
	}

	claims, err := verifyRefreshToken(refreshToken.Value)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки JWT токена: %v", err)
	}

	return claims, nil
}

func verifyRefreshToken(tokenString string) (*models.RefreshToken, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.RefreshToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неправильный алгоритм: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JWT токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("неверный JWT токен")
	}

	claims, ok := token.Claims.(*models.RefreshToken)
	if !ok {
		return nil, fmt.Errorf("ошибка парсинга JWT токена: %w", err)
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("JWT токен истек")
	}

	return claims, nil
}
