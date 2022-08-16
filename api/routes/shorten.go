package routes

import (
	"time"

	"github.com/gofiber/fiber"
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

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Cannot parse JSON."})
	}

	// implement rate limiting

	// check if the domain is an actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid URL."})
	}

	// check for domain error

	if !helpers.RemoveDomainError(body.URL){
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Wrong domain, don't try to play me. HeHeHe"})
	}

	// enforce https, SSL

	body.URL := helpers.EnforceHTTPS(body.URL)

	// check if the URL is already in the database

	if helpers.CheckIf
}
