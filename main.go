package main

import (
	"chat-app-go/model/entity"
	"chat-app-go/route"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
    hub := entity.NewHub()
    go hub.Run()

	app := fiber.New()
	app.Use("/ws", route.AllowUpgrade)
	app.Use("/ws/chat", websocket.New(entity.SendMessage(hub)))

	log.Fatal(app.Listen(":8000"))
}

