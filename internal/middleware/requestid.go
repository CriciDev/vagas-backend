package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	HeaderRequestID = "X-Request-ID"
	ContextKey      = "request_id"
)

var (
	validRequestID  = regexp.MustCompile(`^[A-Za-z0-9._-]{1,128}$`)
	fallbackCounter atomic.Uint64
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetHeader(HeaderRequestID)
		if !validRequestID.MatchString(id) {
			id = generateRequestID()
		}

		ctx.Set(ContextKey, id)
		ctx.Header(HeaderRequestID, id)
		ctx.Next()
	}
}

func FromContext(ctx *gin.Context) string {
	if value, ok := ctx.Get(ContextKey); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}

func AccessLogFormatter(param gin.LogFormatterParams) string {
	id, _ := param.Keys[ContextKey].(string)
	return fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %#v request_id=%s\n",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
		id,
	)
}

func generateRequestID() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err == nil {
		return hex.EncodeToString(buffer)
	}
	return fmt.Sprintf("fallback-%d-%d", time.Now().UnixNano(), fallbackCounter.Add(1))
}
