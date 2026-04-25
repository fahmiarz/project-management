package controllers

import (
	"time"

	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/services"
	"github.com/fahmiarz/project-management/utils"
	"github.com/gofiber/fiber/v2"
)

type CardController struct {
	service services.CardService
}

func NewCardController(s services.CardService) *CardController {
	return &CardController{service: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string `json:"list_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		DueDate time.Time `json:"due_date"`
		Position int `json:"position"`
	}

	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err !=nil{
		return utils.BadRequest(ctx, "Gagal mengambil data", err.Error())
	}
	
	card := &models.Card {
		Title: req.Title,
		Description: req.Description,
		DueDate: &req.DueDate,
		Position: req.Position,
	}

	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal membuat card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dibuat", card)
}