package Response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
}

func (response *Response) JsonSuccessResponse(ctx *gin.Context) {
	response.JsonResponseWithoutError(ctx, "success")
}

func (response *Response) JsonResponseWithoutError(ctx *gin.Context, data interface{}) {
	response.JsonResponse(ctx, data, "", http.StatusOK)
}

func (response *Response) JsonResponseWithServerError(ctx *gin.Context, err interface{}) {
	response.JsonResponseWithError(ctx, err, http.StatusInternalServerError)
}

func (response *Response) JsonResponseWithError(ctx *gin.Context, err interface{}, statusCode int) {
	switch err.(type) {
	case error:
		response.JsonResponse(ctx, "", err.(error).Error(), statusCode)
	default:
		response.JsonResponse(ctx, "", err, statusCode)
	}

}

func (response *Response) JsonResponse(ctx *gin.Context, data interface{}, error interface{}, statusCode int) {
	ctx.JSON(statusCode, gin.H{
		"data": data,
		"code": statusCode,
		"msg":  error,
	})
	ctx.Abort()
}
