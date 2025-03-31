package routes

import (
	"net/http"

	"github.com/AbhinavXJ/go-url-shortener/api/database"
	"github.com/gin-gonic/gin"
)

func DeleteURL(c *gin.Context) {
	shortID := c.Param("shortID")

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortID).Result()

	if err != nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "cant delete,url data not found for this short url",
		})
		return
	}

	er := r.Del(database.Ctx, shortID)

	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cant delete,cant connect to redis server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "shortened url successfully deleted",
	})
}
