package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(500*time.Second),
		timeout.WithResponse(testResponse),
	)
}

func testResponse(c *gin.Context) {
	c.String(http.StatusRequestTimeout, "timeout")
}
