package exception

import (
	"go.uber.org/zap"
)

type ExceptionHandlerInterface interface {
	Reporter(logger *zap.Logger, err error, trace string)
	Handle(err error, trace string)
}
