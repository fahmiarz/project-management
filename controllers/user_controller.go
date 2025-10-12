package controllers

import (

	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/services"
	"github.com/fahmiarz/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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
	var UserResp models.UserResponse
	_= copier.Copy(&UserResp, &user)
	return utils.Success(ctx, "Registration Success", UserResp)
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Invalid Data", err.Error())
	}
	user, err  := c.service.Login(body.Email, body.Password)
	if err !=nil {
		return utils.Unauthorized(ctx, "Login Failed", err.Error())
	}

	token,_ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken,_ := utils.GenerateRefreshToken(user.InternalID)

	var userResp models.UserResponse
	_= copier.Copy(&userResp, &user)
	return utils.Success(ctx, "Login Successfully",fiber.Map{
		"access_token" 	: token,
		"refresh_token" : refreshToken,
		"user" 			: userResp,
	} )
}