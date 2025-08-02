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
			return "", fmt.Errorf("invalid token signature")
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

		parsedScopes := strings.Split(claims["scope"].(string), " ")
		c.Set("scopes", parsedScopes)

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
			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("unauthorized request"))
			ctx.Abort()
			return
		}
		ctx.Set("userId", uid)
		ctx.Next()
	}
}

// RequireScopes returns a Gin middleware that checks if the user has the required OAuth scopes.
// If the user has the 'admin' scope, access is always granted regardless of requiredScopes.
func RequireScopes(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userScopes []string

		rawScopes, exists := c.Get("scopes")
		if !exists {
			httputil.NewError(c, http.StatusForbidden, errors.New("scopes not found in context"))
			c.Abort()
			return
		}

		switch v := rawScopes.(type) {
		case string:
			userScopes = strings.Fields(v)
		case []string:
			userScopes = v
		default:
			httputil.NewError(c, http.StatusForbidden, errors.New("invalid scopes format in context"))
			c.Abort()
			return
		}

		if slices.Contains(userScopes, "admin") {
			c.Set("isAdmin", true)
			c.Next()
			return
		}

		scopeSet := make(map[string]struct{}, len(userScopes))
		for _, s := range userScopes {
			scopeSet[s] = struct{}{}
		}

		for _, scope := range requiredScopes {
			if _, ok := scopeSet[scope]; !ok {
				httputil.NewError(c, http.StatusForbidden, fmt.Errorf("missing required scope: %s", scope))
				c.Abort()
				return
			}
		}

		c.Set("isAdmin", false)
		c.Next()
	}
}
