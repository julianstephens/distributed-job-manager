package httputil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Key   string
	Value any
}

// GetUserId parses the 'userId' path param from Gin context
func GetUserId(ctx *gin.Context) (userId string) {
	val, ok := ctx.Get("userId")
	if !ok {
		NewError(ctx, http.StatusUnauthorized, errors.New("no user id provided"))
		return
	}

	userId = val.(string)

	return
}

// GetId parses the 'id' path param from Gin context
func GetId(ctx *gin.Context) (id string) {
	id = ctx.Param("id")
	if id == "" {
		NewError(ctx, http.StatusUnauthorized, errors.New("no id provided"))
		return
	}
	return
}
