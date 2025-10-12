package routes

import (
	"log"

	"github.com/fahmiarz/project-management/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Setup(app *fiber.App, uc *controllers.UserController) {
	err := godotenv.Load()
	if err !=nil {
		log.Fatal("Error loading .env file")
	}
	app.Post("/v1/auth/register", uc.Register)
}