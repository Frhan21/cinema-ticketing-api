package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Daftar kode warna ANSI untuk CLI
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
)

// MethodColor memberikan warna spesifik untuk setiap tipe HTTP Method
func MethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorCyan
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorGreen
	case "OPTIONS":
		return colorPurple
	default:
		return colorReset
	}
}

// StatusCodeColor memberikan warna spesifik untuk setiap rentang status code
func StatusCodeColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen
	case code >= 300 && code < 400:
		return colorCyan
	case code >= 400 && code < 500:
		return colorYellow
	default:
		return colorRed
	}
}

// LoggerMiddleware formats request logs using default log configuration with Colors
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// Dapatkan warna yang sesuai
		mColor := MethodColor(method)
		sColor := StatusCodeColor(statusCode)

		if len(errorMessage) > 0 {
			log.Printf("%s[ERROR]%s %v |%s %3d %s| %13v | %15s |%s %-7s %s| %s | %s%s%s",
				colorRed, colorReset,
				end.Format("2006/01/02 - 15:04:05"),
				sColor, statusCode, colorReset,
				latency,
				clientIP,
				mColor, method, colorReset,
				path,
				colorYellow, errorMessage, colorReset,
			)
		} else {
			log.Printf("%s[INFO]%s  %v |%s %3d %s| %13v | %15s |%s %-7s %s| %s",
				colorGreen, colorReset,
				end.Format("2006/01/02 - 15:04:05"),
				sColor, statusCode, colorReset,
				latency,
				clientIP,
				mColor, method, colorReset,
				path,
			)
		}
	}
}
