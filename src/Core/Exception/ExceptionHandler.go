package exception

import (
	"go.uber.org/zap"
)

type ExceptionHandler struct {
	ExceptionHandlerInterface
}

func (this *ExceptionHandler) Reporter(logger *zap.Logger, err error, trace string) {
	logger.Debug(err.Error())
}

func (this *ExceptionHandler) Handle(err error, trace string) {

}
