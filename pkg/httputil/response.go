package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPResponse[T any] struct {
	Message string `json:"message" `
	Data    T      `json:"data"`
}

type Options struct {
	IsCrudHandler bool
	HttpMsgMethod HTTPMethod
	Status        int
}

type HTTPMethod int64

const (
	Get HTTPMethod = iota
	Post
	Patch
	Put
	Delete
)

func (s HTTPMethod) String() string {
	switch s {
	case Get:
		return "get"
	case Post:
		return "post"
	case Patch:
		return "patch"
	case Put:
		return "put"
	case Delete:
		return "delete"
	}
	return "unknown"
}

func NewResponse[T any](ctx *gin.Context, data T, opts Options) {
	status := If(opts.Status > 0, opts.Status, http.StatusOK)

	if !opts.IsCrudHandler {
		ctx.JSON(status, data)
		return
	}

	var message string

	switch opts.HttpMsgMethod {
	default:
		message = "resource(s) retrieved successfully"
	case Post:
		status = http.StatusCreated
		message = "resource(s) created successfully"
	case Put:
		message = "resource(s) updated successfully"
	case Delete:
		message = "resource(s) deleted successfully"
	}

	res := HTTPResponse[T]{
		Message: message,
		Data:    data,
	}
	ctx.JSON(status, res)
}
