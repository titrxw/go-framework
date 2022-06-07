package middleware

import (
	"github.com/gin-gonic/gin"
	session "github.com/titrxw/go-framework/src/Core/Session"
	"net/http"
)

type SessionMiddleware struct {
	MiddlewareAbstract
	Session *session.Session
}

func NewSessionMiddleware(appSession *session.Session) *SessionMiddleware {
	appSession.Init()
	return &SessionMiddleware{Session: appSession}
}

func (this SessionMiddleware) Process(ctx *gin.Context) {
	err := this.Session.Start(ctx)
	if err != nil {
		this.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.Next()
}
