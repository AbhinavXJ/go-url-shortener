package routes

import (
	"encoding/json"
	"net/http"

	"github.com/AbhinavXJ/go-url-shortener/api/database"
	"github.com/gin-gonic/gin"
)

type TagRequest struct {
	ShortID string `json:"shortID"`
	Tag     string `json:"tag"`
}

func AddTag(c *gin.Context) {
	// Add tag to URL
	var tagRequest TagRequest
	if err := c.BindJSON(&tagRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shortId := tagRequest.ShortID
	tag := tagRequest.Tag

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, shortId).Result()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found for given shorturl"})
		return
	}
	//in redits data is stored in string format so unmarshal to format it
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		//if data is not in json object,assume it as plain string
		data = make(map[string]interface{})
		data["data"] = val
	}

	//check if tag is already present and its a slice of strings

	var tags []string

	if existingTags, ok := data["tags"].([]interface{}); ok {
		for _, t := range existingTags {
			if strTag, ok := t.(string); ok {
				tags = append(tags, strTag)
			}
		}
	}
	//check for duplicate tags
	for _, existingTag := range tags {
		if existingTag == tag {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "tag already exists",
			})

			return
		}
	}

	//adding new tag to tag slice
	tags = append(tags, tag)
	data["tags"] = tags

	//marshall the updated data back to json
	updatedData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong while parsing data to redis",
		})
		return
	}

	err = r.Set(database.Ctx, shortId, updatedData, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	//respond with the updated data
	c.JSON(http.StatusOK, gin.H{
		"message": "tag added successfully",
		"data":    data,
	})

}
