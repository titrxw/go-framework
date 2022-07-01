package exception

import (
	"go.uber.org/zap"
)

type ExceptionHandler struct {
	ExceptionHandlerInterface
}

func (exceptionHandler *ExceptionHandler) Reporter(logger *zap.Logger, err error, trace string) {
	logger.Debug(err.Error())
}

func (exceptionHandler *ExceptionHandler) Handle(err error, trace string) {

}
