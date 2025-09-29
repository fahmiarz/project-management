package seed

import (
	"log"

	"github.com/fahmiarz/project-management/config"
	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/utils"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin123")
	
	admin := models.User {
		Name: "Super admin",
		Email: "admin@example.com",
		Password: password,
		Role: "admin",
	}
	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error; err != nil {
		log.Println("Failed too seed admin", err)
	}else{
		log.Println("Admin User seeded successfully")
	}
}