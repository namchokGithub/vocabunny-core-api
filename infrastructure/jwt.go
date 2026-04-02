package infrastructure

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/namchokGithub/vocabunny-core-api/configs"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
)

type JWTManager struct {
	secret          []byte
	issuer          string
	audience        string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type AccessClaims struct {
	TokenUse string `json:"token_use"`
	Scope    string `json:"scope,omitempty"`
	jwt.RegisteredClaims
}

func NewJWTManager(cfg configs.JWTConfig) *JWTManager {
	return &JWTManager{
		secret:          []byte(cfg.Secret),
		issuer:          cfg.Issuer,
		audience:        cfg.Audience,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (m *JWTManager) GenerateAccessToken(subject string) (string, error) {
	now := time.Now()
	claims := AccessClaims{
		TokenUse: domain.TokenUseAccess,
		Scope:    domain.TokenScopeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    m.issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings{m.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) AccessTokenTTLSeconds() int64 {
	return int64(m.accessTokenTTL / time.Second)
}

func (m *JWTManager) GenerateRefreshToken(subject string) (string, error) {
	return "", helper.Internal(
		"refresh_token_not_implemented",
		"refresh token issuance is not implemented yet",
		nil,
	)
}

func (m *JWTManager) RefreshTokenTTLSeconds() int64 {
	return int64(m.refreshTokenTTL / time.Second)
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
			if claims.TokenUse != domain.TokenUseAccess {
				return helper.RespondError(c, helper.Unauthorized("invalid_token", "token is not an access token", nil))
			}

			c.Set("jwt_subject", claims.Subject)
			return next(c)
		}
	}
}
