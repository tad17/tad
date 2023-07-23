package fiber

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewFiber() *fiber.App {
	// инициализация fiber
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(favicon.New())
	app.Use(compress.New())
	fmt.Printf("fiber init\n")
	return app
}
