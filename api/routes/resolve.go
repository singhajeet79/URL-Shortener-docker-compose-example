package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber"
	"github.com/hamees-sayed/URL-Shortener/database"
)

func ResolveURL(c *fiber.Ctx) error {

	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found in database."})
	} else if err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "could not connect to database."})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)

}

