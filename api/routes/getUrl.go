package routes

import (
	"net/http"

	"github.com/AbhinavXJ/go-url-shortener/api/database"
	"github.com/gin-gonic/gin"
)

func GetByShortID(c *gin.Context) {

	shortID := c.Param("shortID")
	r := database.CreateClient(0)

	val, err := r.Get(database.Ctx, shortID).Result()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Url not found for given ShortId",
		})
		return
	}

	c.Redirect(http.StatusFound, val)

}
