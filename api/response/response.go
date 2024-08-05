package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

type List struct {
	Items   interface{} `json:"items"`
	Total   int64       `json:"total"`
	Page    int64       `json:"page"`
	PerPage int64       `json:"perPage"`
}

func (r *Response) String(ctx *gin.Context) string {
	return fmt.Sprintf("code: %d, msg: %s, data: %s", r.Code, r.Msg, r.Data)
}

func successRender(ctx *gin.Context, status int, data ...interface{}) {
	r := &Response{Code: status, Msg: http.StatusText(status)}
	if len(data) != 0 {
		r.Data = data[0]
	}
	ctx.JSON(status, r)
}

func errorRender(ctx *gin.Context, status int, data ...interface{}) {
	if len(data) == 0 {
		ctx.Status(status)
		return
	}
	r := &Response{Code: status, Msg: data[0]}
	if len(data) > 1 {
		r.Data = data[1]
	}
	if err, ok := data[0].(error); ok {
		r.Msg = err.Error()
	}
	if status < 500 {
		log.Warn().Msg(r.String(ctx))
	} else {
		log.Error().Msg(r.String(ctx))
	}
	ctx.JSON(status, r)
}

//successResponses

func Data(ctx *gin.Context, data ...interface{}) {
	successRender(ctx, http.StatusOK, data...)
}

func Created(ctx *gin.Context, data ...interface{}) {

	successRender(ctx, http.StatusCreated, data...)
}

func Accepted(ctx *gin.Context, data ...interface{}) {
	successRender(ctx, http.StatusAccepted, data...)
}

func NoContent(ctx *gin.Context, data ...interface{}) {
	successRender(ctx, http.StatusNoContent, data...)
}

//errorResponses

func BadRequest(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusBadRequest, data...)
	ctx.Abort()
}

func NotAcceptable(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusNotAcceptable, data...)
	ctx.Abort()
}

func Conflict(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusConflict, data...)
	ctx.Abort()
}

func NotFound(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusNotFound, data...)
	ctx.Abort()
}

func Unauthorized(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusUnauthorized, data...)
	ctx.Abort()
}

func Forbidden(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusForbidden, data...)
	ctx.Abort()
}

func MethodNotAllowed(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusMethodNotAllowed, data...)
	ctx.Abort()
}

func InternalServerError(ctx *gin.Context, data ...interface{}) {
	errorRender(ctx, http.StatusInternalServerError, data...)
	ctx.Abort()
}
