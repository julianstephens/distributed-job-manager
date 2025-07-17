package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type ValidationError struct {
	FieldError validator.FieldError
}

func NewError(ctx *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

func (v ValidationError) NewFieldError() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + v.FieldError.Field() + "'")
	sb.WriteString(", condition: " + v.FieldError.ActualTag())

	if v.FieldError.Param() != "" {
		sb.WriteString(" { " + v.FieldError.Param() + " }")
	}

	if v.FieldError.Value() != nil && v.FieldError.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", v.FieldError.Value()))
	}

	return sb.String()

}

func HandleFieldError(ctx *gin.Context, err error) {
	for _, fieldErr := range err.(validator.ValidationErrors) {
		msg := ValidationError{FieldError: fieldErr}.NewFieldError()
		NewError(ctx, http.StatusBadRequest, errors.New(msg))
		return
	}
}
