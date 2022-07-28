package routes

import (
	"genesis_test_case/src/pkg/delivery/http"
	"genesis_test_case/src/pkg/delivery/http/middleware"
	"genesis_test_case/src/pkg/repository"
	"genesis_test_case/src/pkg/usecase"

	"github.com/gofiber/fiber/v2"
)

var (
	MailingHandler *http.MailingHandler
)

func InitRoutes(app *fiber.App) error {
	repos, err := repository.NewRepositories()
	if err != nil {
		return err
	}

	usecases := usecase.NewUsecases(repos)
	MailingHandler = http.NewMailingHandler(usecases)

	middleware.FiberMiddleware(app)
	PublicRoutes(app)
	return nil
}
