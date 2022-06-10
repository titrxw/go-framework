package session

import (
	"context"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gin-gonic/gin"
	config "github.com/titrxw/go-framework/src/Core/Config"
	"net/http"
	"time"
)

type Session struct {
	scs.SessionManager
	storageResolver func() scs.Store
}

func NewSession(sessionConfig *config.Session, cookieConfig *config.Cookie) *Session {
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

func (this *Session) getContext(ctx *gin.Context) context.Context {
	return ctx.Request.Context()
}

func (this *Session) SetStorageResolver(storageResolver func() scs.Store) {
	this.storageResolver = storageResolver
}

func (this *Session) Init() {
	this.Store = this.storageResolver()
}

func (this *Session) Start(ctx *gin.Context) error {
	cookie, err := ctx.Cookie(this.Cookie.Name)
	if err != nil {
		cookie = ""
	}

	_ctx, err := this.Load(this.getContext(ctx), cookie)
	if err != nil {
		this.ErrorFunc(ctx.Writer, ctx.Request, err)
		return err
	}
	ctx.Request = ctx.Request.WithContext(_ctx)

	return nil
}

func (this *Session) Set(ctx *gin.Context, key string, value interface{}) error {
	this.Put(this.getContext(ctx), key, value)

	return this.saveAndResponse(ctx, false)
}

func (this *Session) Get(ctx *gin.Context, key string) interface{} {
	return this.SessionManager.Get(this.getContext(ctx), key)
}

func (this *Session) Delete(ctx *gin.Context, key string) error {
	this.Remove(this.getContext(ctx), key)

	return this.saveAndResponse(ctx, false)
}

func (this *Session) Destroy(ctx *gin.Context) error {
	err := this.SessionManager.Destroy(this.getContext(ctx))
	if err != nil {
		return err
	}

	return this.saveAndResponse(ctx, true)
}

func (this *Session) saveAndResponse(ctx *gin.Context, isDelete bool) error {
	token, expire, err := this.Commit(this.getContext(ctx))
	if err != nil {
		this.ErrorFunc(ctx.Writer, ctx.Request, err)
		return err
	}

	return this.responseCookie(ctx, token, expire, isDelete)
}

func (this *Session) responseCookie(ctx *gin.Context, token string, expire time.Time, isDelete bool) error {
	if isDelete {
		this.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, "", time.Time{})
		return nil
	}

	cookie, _ := ctx.Cookie(this.Cookie.Name)
	if cookie != "" {
		return nil
	}

	this.WriteSessionCookie(ctx.Request.Context(), ctx.Writer, token, expire)
	return nil
}
