package main

import (
	"github.com/laurentsbrndn/accounting-app/rest-api/internal/api"
	"github.com/laurentsbrndn/accounting-app/rest-api/internal/config"
	"github.com/laurentsbrndn/accounting-app/rest-api/internal/connection"
	"github.com/laurentsbrndn/accounting-app/rest-api/internal/repository"
	"github.com/laurentsbrndn/accounting-app/rest-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

	userRepository := repository.NewUser(dbConnection)

	authService := service.NewAuth(cnf, userRepository)

	api.NewAuth(app, authService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
