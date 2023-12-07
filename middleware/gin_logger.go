package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/make-money-fast/plugins/logger"
	"io"
	"strings"
	"time"
)

// responseWriter 记录每个请求的日志
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseWriter) isJsonResponse() bool {
	return strings.HasPrefix(w.Header().Get("Content-Type"), "application/json")
}

func newResponseWriter(ctx *gin.Context) *responseWriter {
	return &responseWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: ctx.Writer,
	}
}

type LogConfigure struct {
	EnableRequestBody  bool
	EnableResponseBody bool
	SkipLoggerFunc     func(ctx *gin.Context) bool
}

func Logger(log *logger.Entry, config LogConfigure) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if config.SkipLoggerFunc(ctx) {
			ctx.Next()
			return
		}
		var (
			start      = time.Now()
			respWriter *responseWriter
			fields     = []*logger.Field{
				logger.Any("scene", "http_server_request"),
				logger.Any("method", ctx.Request.Method),
				logger.Any("uri", ctx.Request.URL.Path),
				logger.Any("full_uri", ctx.Request.RequestURI),
				logger.Any("ip", ctx.ClientIP()),
				logger.Any("header", ctx.Request.Header),
			}
		)

		if config.EnableRequestBody {
			var data, err = ctx.GetRawData()
			if err != nil {
				log.Error(ctx.Request.Context(), "http request read error", logger.Any("scene", "http_server_request"), logger.Err(err))
				return
			}
			fields = append(fields, logger.Any("request_body", logger.RawJSON(data)))
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		}

		respWriter = newResponseWriter(ctx)
		ctx.Writer = respWriter
		ctx.Next()

		var (
			gErr   error
			ginErr = ctx.Errors.Last()
		)
		if ginErr != nil {
			gErr = ginErr.Unwrap()
			fields = append(fields, logger.Err(gErr))
		}

		fields = append(fields,
			logger.Any("status", ctx.Writer.Status()),
			logger.Any("size", ctx.Writer.Size()),
			logger.Any("latency", time.Since(start).Milliseconds()),
		)

		if config.EnableResponseBody && respWriter.isJsonResponse() {
			fields = append(fields,
				logger.Any("response_body", logger.RawJSON(respWriter.body.Bytes())),
			)

		}

		if ginErr != nil {
			log.Error(
				ctx,
				"http request and response",
				fields...,
			)
		} else {
			log.Info(
				ctx,
				"http request and response",
				fields...,
			)
		}
	}
}
