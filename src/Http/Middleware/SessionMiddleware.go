package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	session "github.com/titrxw/go-framework/src/Http/Session"
)

type SessionMiddleware struct {
	MiddlewareAbstract
	Session *session.Session
}

func NewSessionMiddleware(appSession *session.Session) *SessionMiddleware {
	appSession.Init()
	return &SessionMiddleware{Session: appSession}
}

func (sessionMiddleware SessionMiddleware) Process(ctx *gin.Context) {
	err := sessionMiddleware.Session.Start(ctx)
	if err != nil {
		sessionMiddleware.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.Next()
}
