package controllers

import (
	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/services"
	"github.com/fahmiarz/project-management/utils"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	service services.UserService
}

func NewUserController (s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Data Parsing Failed", err.Error())
	}
	if err := c.service.Register(user); err !=nil {
		return utils.BadRequest(ctx, "Registration Failed", err.Error())
	}
	return utils.Success(ctx, "Registration Success", user)
}