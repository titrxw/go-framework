package session

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gin-gonic/gin"
	config "github.com/titrxw/go-framework/src/Core/Config"
)

type Session struct {
	scs.SessionManager
	storageResolver func() scs.Store
}

func NewSession(sessionConfig config.Session, cookieConfig config.Cookie) *Session {
	session := &Session{
		SessionManager: *scs.New(),
	}
	session.ErrorFunc = func(writer http.ResponseWriter, request *http.Request, err error) {

	}
	if sessionConfig.Lifetime > 0 {
		session.Lifetime = sessionConfig.Lifetime
	}
	if sessionConfig.IdleTimeout > 0 {
		session.IdleTimeout = sessionConfig.IdleTimeout
	}
	if cookieConfig.Name != "" {
		session.Cookie.Name = cookieConfig.Name
	}
	if cookieConfig.Domain != "" {
		session.Cookie.Domain = cookieConfig.Domain
	}
	if cookieConfig.Path != "" {
		session.Cookie.Path = cookieConfig.Path
	}
	session.Cookie.HttpOnly = cookieConfig.HttpOnly
	session.Cookie.Persist = cookieConfig.Persist
	session.Cookie.Secure = cookieConfig.Secure
	if cookieConfig.SameSite == "Lax" {
		session.Cookie.SameSite = http.SameSiteLaxMode
	} else if cookieConfig.SameSite == "Strict" {
		session.Cookie.SameSite = http.SameSiteStrictMode
	} else if cookieConfig.SameSite == "None" {
		session.Cookie.SameSite = http.SameSiteNoneMode
	} else {
		session.Cookie.SameSite = http.SameSiteDefaultMode
	}

	session.SetStorageResolver(func() scs.Store {
		return memstore.New()
	})

	return session
}

func (session *Session) getContext(ctx *gin.Context) context.Context {
	return ctx.Request.Context()
}

func (session *Session) SetStorageResolver(storageResolver func() scs.Store) {
	session.storageResolver = storageResolver
}

func (session *Session) Init() {
	session.Store = session.storageResolver()
}

func (session *Session) Start(ctx *gin.Context) error {
	cookie, err := ctx.Cookie(session.Cookie.Name)
	if err != nil {
		cookie = ""
	}

	_ctx, err := session.Load(session.getContext(ctx), cookie)
	if err != nil {
		session.ErrorFunc(ctx.Writer, ctx.Request, err)
		return err
	}
	ctx.Request = ctx.Request.WithContext(_ctx)

	return nil
}

func (session *Session) Set(ctx *gin.Context, key string, value interface{}) error {
	session.Put(session.getContext(ctx), key, value)

	return session.saveAndResponse(ctx, false)
}

func (session *Session) Get(ctx *gin.Context, key string) interface{} {
	return session.SessionManager.Get(session.getContext(ctx), key)
}

func (session *Session) Delete(ctx *gin.Context, key string) error {
	session.Remove(session.getContext(ctx), key)

	return session.saveAndResponse(ctx, false)
}

func (session *Session) Destroy(ctx *gin.Context) error {
	err := session.SessionManager.Destroy(session.getContext(ctx))
	if err != nil {
		return err
	}

	return session.saveAndResponse(ctx, true)
}

func (session *Session) saveAndResponse(ctx *gin.Context, isDelete bool) error {
	token, expire, err := session.Commit(session.getContext(ctx))
	if err != nil {
		session.ErrorFunc(ctx.Writer, ctx.Request, err)
		return err
	}

	return session.responseCookie(ctx, token, expire, isDelete)
}

func (session *Session) responseCookie(ctx *gin.Context, token string, expire time.Time, isDelete bool) error {
	if isDelete {
		session.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, "", time.Time{})
		return nil
	}

	cookie, _ := ctx.Cookie(session.Cookie.Name)
	if cookie != "" {
		return nil
	}

	session.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, token, expire)
	return nil
}
