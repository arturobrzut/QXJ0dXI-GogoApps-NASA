package url_collector

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

func Serv(config *serviceConfig) error {
	return router(config).Run(":" + strconv.Itoa(config.port))
}

func router(config *serviceConfig) *gin.Engine {
	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
	})
	router.GET("/pictures", func(c *gin.Context) {
		var dates picturesDate
		if err := c.ShouldBindWith(&dates, binding.Query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		urls, err := collect(config, dates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"urls": urls})
	})
	return router
}
