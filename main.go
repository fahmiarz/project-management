package main

import (
	"log"

	"github.com/fahmiarz/project-management/config"
	"github.com/fahmiarz/project-management/controllers"
	"github.com/fahmiarz/project-management/database/seed"
	"github.com/fahmiarz/project-management/repositories"
	"github.com/fahmiarz/project-management/routes"
	"github.com/fahmiarz/project-management/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()

	app:= fiber.New()

	//user
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	//board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	//list
	listPosRepo := repositories.NewListPositionRepository()
	listRepo := repositories.NewListRepository()
	listService := services.NewListService(listRepo, boardRepo, listPosRepo)
	listController := controllers.NewListController(listService)

	//card
	cardRepo := repositories.NewCardRepository()
	cardService := services.NewCardService(cardRepo, listRepo, userRepo)
	cardController := controllers.NewCardController(cardService)

	routes.Setup(app, 
		userController,
		boardController,
		listController,
		cardController,
	)

	port:= config.AppConfig.AppPort
	log.Println("Server is running on port : ", port)
	log.Fatal(app.Listen(":" + port)) 
}