package Response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
}

func (this *Response) JsonSuccessResponse(ctx *gin.Context) {
	this.JsonResponseWithoutError(ctx, "success")
}

func (this *Response) JsonResponseWithoutError(ctx *gin.Context, data interface{}) {
	this.JsonResponse(ctx, data, "", http.StatusOK)
}

func (this *Response) JsonResponseWithServerError(ctx *gin.Context, err interface{}) {
	this.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
}

func (this *Response) JsonResponseWithError(ctx *gin.Context, err interface{}, statusCode int) {
	switch err.(type) {
	case error:
		this.JsonResponse(ctx, "", err.(error).Error(), statusCode)
	default:
		this.JsonResponse(ctx, "", err, statusCode)
	}

}

func (this *Response) JsonResponse(ctx *gin.Context, data interface{}, error interface{}, statusCode int) {
	ctx.JSON(statusCode, gin.H{
		"data": data,
		"code": statusCode,
		"msg":  error,
	})
	ctx.Abort()
}
