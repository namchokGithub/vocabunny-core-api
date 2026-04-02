package infrastructure

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/configs"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

type AccessClaims struct {
	Subject string `json:"sub"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg configs.JWTConfig) *JWTManager {
	return &JWTManager{
		secret: []byte(cfg.Secret),
		issuer: cfg.Issuer,
		ttl:    cfg.AccessTokenTTL,
	}
}

func (m *JWTManager) GenerateAccessToken(subject string) (string, error) {
	claims := AccessClaims{
		Subject: subject,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) AccessTokenTTLSeconds() int64 {
	return int64(m.ttl / time.Second)
}

func (m *JWTManager) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return helper.RespondError(c, helper.Unauthorized("missing_token", "authorization token is required", nil))
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return helper.RespondError(c, helper.Unauthorized("invalid_token", "authorization header must use bearer token", nil))
			}

			token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}

				return m.secret, nil
			})
			if err != nil || !token.Valid {
				return helper.RespondError(c, helper.Unauthorized("invalid_token", "token is invalid", err))
			}

			claims, ok := token.Claims.(*AccessClaims)
			if !ok {
				return helper.RespondError(c, helper.Unauthorized("invalid_token", "token claims are invalid", nil))
			}

			c.Set("jwt_subject", claims.Subject)
			return next(c)
		}
	}
}
