package middleware

import (
	"github.com/gin-gonic/gin"
	exception "github.com/titrxw/go-framework/src/Core/Exception"
)

type ExceptionMiddleware struct {
	MiddlewareAbstract
	HandlerExceptions *exception.HandlerExceptions
}

func (this ExceptionMiddleware) Process(ctx *gin.Context) {
	this.HandlerExceptions.RegisterExceptionHandle(ctx)
}
