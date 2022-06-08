package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	exception "github.com/titrxw/go-framework/src/Core/Exception"
	"io"
	"net/http"
)

type ExceptionRecoverLogger struct {
	io.Writer
	HandlerExceptions *exception.HandlerExceptions
}

func (this *ExceptionRecoverLogger) Write(p []byte) (n int, err error) {
	this.HandlerExceptions.GetExceptionHandler().Reporter(this.HandlerExceptions.Logger, fmt.Errorf("%v", err), string(this.HandlerExceptions.Stack(4)))
	return len(p), nil
}

type ExceptionMiddleware struct {
	MiddlewareAbstract
	HandlerExceptions *exception.HandlerExceptions
}

func (this ExceptionMiddleware) Process(ctx *gin.Context) {
	gin.CustomRecoveryWithWriter(&ExceptionRecoverLogger{
		HandlerExceptions: this.HandlerExceptions,
	}, func(ctx *gin.Context, err interface{}) {
		if gin.Mode() == gin.DebugMode {
			this.JsonResponseWithError(ctx, string(this.HandlerExceptions.Stack(4)), http.StatusInternalServerError)
		} else {
			this.JsonResponseWithError(ctx, "系统内部错误", http.StatusInternalServerError)
		}
	})
}
