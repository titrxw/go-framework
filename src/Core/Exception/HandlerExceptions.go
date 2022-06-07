package exception

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type ExceptionRecoverLogger struct {
	io.Writer
	HandlerExceptions *HandlerExceptions
}

func (this *ExceptionRecoverLogger) Write(p []byte) (n int, err error) {
	this.HandlerExceptions.GetExceptionHandler().Reporter(this.HandlerExceptions.Logger, fmt.Errorf("%v", err))
	return len(p), nil
}

type HandlerExceptions struct {
	ExceptionHandler ExceptionHandlerInterface
	Logger           *zap.Logger
}

func (this *HandlerExceptions) SetExceptionHandler(exceptionHandler ExceptionHandlerInterface) {
	this.ExceptionHandler = exceptionHandler
}

func (this *HandlerExceptions) GetExceptionHandler() ExceptionHandlerInterface {
	if this.ExceptionHandler == nil {
		this.ExceptionHandler = new(ExceptionHandler)
	}
	return this.ExceptionHandler
}

func (this *HandlerExceptions) newCustomRecovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(&ExceptionRecoverLogger{
		HandlerExceptions: this,
	}, func(ctx *gin.Context, err interface{}) {
		this.GetExceptionHandler().Handle(ctx, fmt.Errorf("%v", err))
	})
}

func (this *HandlerExceptions) RegisterExceptionHandle(ctx *gin.Context) {
	if ctx == nil {
		ctx = &gin.Context{}
	}
	this.newCustomRecovery()(ctx)
}
