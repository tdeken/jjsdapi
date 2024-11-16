package server

import (
	"github.com/gofiber/fiber/v2"
	"jjsdapi/internal/fiber/mid"
)

type WebServer struct {
	Server *fiber.App
}

var Web = server()

func server() *WebServer {
	ser := fiber.New()

	ser.Use(
		mid.Recover(),
		mid.Logger(),
		mid.Cors(),
		mid.TraceId(),
	)

	return &WebServer{
		Server: ser,
	}
}
