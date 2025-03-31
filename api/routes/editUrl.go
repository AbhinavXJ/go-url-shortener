package routes

import (
	"net/http"
	"time"

	"github.com/AbhinavXJ/go-url-shortener/api/database"
	"github.com/AbhinavXJ/go-url-shortener/api/models"
	"github.com/gin-gonic/gin"
)

func EditURL(c *gin.Context) {

	shortID := c.Param("shortID")
	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot parse json",
		})
		return
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortID).Result()
	if err != nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "short id not found",
		})
		return
	}

	_, err = r.Set(database.Ctx, shortID, val, body.Expiry*3600&time.Second).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot connect to redis server",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "short url content has been updated",
	})
}
