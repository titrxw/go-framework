package exception

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExceptionHandlerInterface interface {
	Handle(ctx *gin.Context, err error)
	Reporter(logger *zap.Logger, err error)
}
