package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/AbhinavXJ/go-url-shortener/api/database"
	"github.com/AbhinavXJ/go-url-shortener/api/models"
	"github.com/AbhinavXJ/go-url-shortener/api/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func ShortenURL(c *gin.Context) {
	var body models.Request

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot parse JSON"})
		return
	}

	//used to store ip add of client
	r2 := database.CreateClient(1)
	defer r2.Close()

	_, err := r2.Get(database.Ctx, c.ClientIP()).Result()

	if err == redis.Nil {
		r2.Set(database.Ctx, c.ClientIP(), os.Getenv("API_QUOTA"), 30*60*time.Second)
	} else {
		val, _ := r2.Get(database.Ctx, c.ClientIP()).Result()
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.ClientIP()).Result()
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
			return
		}

	}

	if !govalidator.IsURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid url",
		})
		return
	}
	//yaha ! hai
	if !utils.IsDifferentDomain(body.URL) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "you cant hack this sytem",
		})
		return
	}

	body.URL = utils.EnsureHttpPrefix(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ := r.Get(database.Ctx, id).Result()

	if val != "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "custom short url already there in database",
		})
		return
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to connect to redis server",
		})
		return
	}

	resp := models.Response{
		Expiry: body.Expiry,
		//30 min baad limit reset hogi
		XRateLimitReset: 30,
		// us 30 min me user 10 requests bhej sakta hai
		XRateRemaining: 10,
		URL:            body.URL,
		CustomShort:    "",
	}

	r2.Decr(database.Ctx, c.ClientIP())

	val, _ = r2.Get(database.Ctx, c.ClientIP()).Result()

	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.ClientIP()).Result()

	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	c.JSON(http.StatusOK, resp)

}
