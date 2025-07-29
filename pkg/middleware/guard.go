package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julianstephens/distributed-job-manager/pkg/config"
	"github.com/julianstephens/distributed-job-manager/pkg/httputil"
	"github.com/julianstephens/distributed-job-manager/pkg/logger"
	"github.com/julianstephens/distributed-job-manager/pkg/models"
)

var (
	ErrInvalidToken = errors.New("authorization token is invalid")
)

type JWTManager struct {
	conf *models.Config
}

func NewJWTManager() *JWTManager {
	conf := config.GetConfig()

	return &JWTManager{
		conf: conf,
	}
}

func (j *JWTManager) extractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func (j *JWTManager) Validate(c *gin.Context) (string, error) {
	tokenString := j.extractToken(c)
	if tokenString == "" {
		logger.Errorf("received malformated token")
		return "", ErrInvalidToken
	}

	k, err := keyfunc.NewDefaultCtx(c.Request.Context(), []string{j.conf.Auth0.JWKSUrl})
	if err != nil {
		return "", fmt.Errorf("failed to create a keyfunc from the jwks URL: %v", err)
	}

	parsed, err := jwt.Parse(tokenString, k.Keyfunc)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return "", ErrInvalidToken
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return "", fmt.Errorf("invalid signature")
		case errors.Is(err, jwt.ErrTokenExpired):
			return "", fmt.Errorf("token expired")
		default:
			return "", fmt.Errorf("failed to parse jwt: %v", err)
		}
	}

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
		iss, err := claims.GetIssuer()
		if err != nil || iss != fmt.Sprintf("https://%s/", j.conf.Auth0.Domain) {
			return "", errors.New("invalid token issuer")
		}

		aud, err := claims.GetAudience()
		if err != nil || !slices.Contains(aud, j.conf.Auth0.Audience) {
			return "", errors.New("invalid token audience")
		}

		return claims.GetSubject()
	} else {
		return "", errors.New("failed to parse token claims")
	}
}

func Guard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jM := NewJWTManager()
		uid, err := jM.Validate(ctx)
		if err != nil {
			httputil.NewError(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}
		ctx.Set("userId", uid)
		ctx.Next()
	}
}
