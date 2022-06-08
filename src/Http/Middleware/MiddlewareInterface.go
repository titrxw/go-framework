package middleware

import "github.com/gin-gonic/gin"

type MiddlewareInterface interface {
	Process(ctx *gin.Context)
}
