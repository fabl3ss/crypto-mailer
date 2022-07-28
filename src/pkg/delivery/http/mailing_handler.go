package http

import (
	"genesis_test_case/src/pkg/domain"
	"genesis_test_case/src/pkg/types/errors"
	"genesis_test_case/src/pkg/usecase"
	"genesis_test_case/src/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type MailingHandler struct {
	usecases *usecase.Usecases
}

func NewMailingHandler(u *usecase.Usecases) *MailingHandler {
	return &MailingHandler{
		usecases: u,
	}
}

func (m *MailingHandler) SendRate(c *fiber.Ctx) error {
	unsent, err := m.usecases.Mailing.SendCurrencyRate()
	if err != nil {
		// Return status 400
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if len(unsent) > 0 {
		return c.JSON(fiber.Map{
			"unsent": unsent,
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (m *MailingHandler) Subscribe(c *fiber.Ctx) error {
	recipient := new(domain.Recipient)

	if err := c.BodyParser(recipient); err != nil {
		// Return status 400 and error message
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Code instance
	validate := utils.GetValidator()

	// Validate code fields
	if err := validate.Struct(recipient); err != nil {
		// Return, if some fields are not valid
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	err := m.usecases.Mailing.Subscribe(recipient)
	if err != nil {
		if err == errors.ErrAlreadyExists {
			// Return status 409
			return c.SendStatus(fiber.StatusConflict)
		}
		// Return status 500 and error message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (m *MailingHandler) GetCurrencyRate(c *fiber.Ctx) error {
	rate, err := m.usecases.Crypto.GetConfigCurrencyRate()
	if err != nil {
		// Return status 400
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(rate)
}
