package routes

import (
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	// Routes for GET method:
	a.Get("/rate", MailingHandler.GetCurrencyRate)

	// Routes for POST method:
	a.Post("/sendEmails", MailingHandler.SendRate)
	a.Post("/subscribe", MailingHandler.Subscribe)
}
