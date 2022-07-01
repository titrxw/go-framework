package middleware

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	exception "github.com/titrxw/go-framework/src/Core/Exception"
)

type ExceptionRecoverLogger struct {
	io.Writer
	HandlerExceptions *exception.HandlerExceptions
}

func (exceptionRecoverLogger *ExceptionRecoverLogger) Write(p []byte) (n int, err error) {
	exceptionRecoverLogger.HandlerExceptions.GetExceptionHandler().Reporter(exceptionRecoverLogger.HandlerExceptions.Logger, fmt.Errorf("%v", err), string(exceptionRecoverLogger.HandlerExceptions.Stack(4)))
	return len(p), nil
}

type ExceptionMiddleware struct {
	MiddlewareAbstract
	HandlerExceptions *exception.HandlerExceptions
}

func (exceptionMiddleware ExceptionMiddleware) Process(ctx *gin.Context) {
	gin.CustomRecoveryWithWriter(&ExceptionRecoverLogger{
		HandlerExceptions: exceptionMiddleware.HandlerExceptions,
	}, func(ctx *gin.Context, err interface{}) {
		if gin.Mode() == gin.DebugMode {
			exceptionMiddleware.JsonResponseWithError(ctx, string(exceptionMiddleware.HandlerExceptions.Stack(4)), http.StatusInternalServerError)
		} else {
			exceptionMiddleware.JsonResponseWithError(ctx, "系统内部错误", http.StatusInternalServerError)
		}
	})
}
