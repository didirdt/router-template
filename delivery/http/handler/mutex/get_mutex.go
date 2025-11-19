package mutex

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMutex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}
