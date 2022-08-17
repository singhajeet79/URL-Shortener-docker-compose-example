package routes

import (
	"os"
	"strconv"
	"time"
	"uuid"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber"
	"github.com/hamees-sayed/URL-Shortener/database"
	"github.com/hamees-sayed/URL-Shortener/helpers"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"customShort"`
	Expiry          time.Duration `json:"expiry"`
	XRateLimit      int           `json:"xRateLimit`
	XRateLimitReset time.Duration `json:"xRateLimitReset"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(Request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON."})
	}

	// implement rate limiting

	r2 := database.CreateClient(1)
	defer r2.Close()
	value, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		value, _ = r2.Get(database.Ctx, c.IP()).Result()
		valueInt, _ := strconv.Atoi(value)
		if valueInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate Limit exceeded.",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// check if the domain is an actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL."})
	}

	// check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Wrong domain, don't try to play me. HeHeHe"})
	}

	// enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	}

	r2.Decr(database.Ctx, c.IP())
}
