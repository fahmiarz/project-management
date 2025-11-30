package controllers

import (
	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/services"
	"github.com/fahmiarz/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	var userID uuid.UUID
	var err error

	board := new(models.Board)
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err !=nil {
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}

	if userID, err = uuid.Parse(claims["pub_id"].(string)); err !=nil {
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}
	board.OwnerPublicID = userID

	if err := c.service.Create(board); err !=nil {
		return utils.BadRequest(ctx, "Gagal menyimpan Data", err.Error())
	}

	return utils.Success(ctx, "Board berhasil dibuat", board)
	
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new (models.Board)

	if err := ctx.BodyParser(board); err !=nil {
		return utils.BadRequest(ctx, "Gagal parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err !=nil {
		return utils.BadRequest(ctx, "ID tidak Valid", err.Error())
	}

	existingBoard, err  := c.service.GetByPublicID(publicID) 
	 if err !=nil {
		return utils.NotFound(ctx, "Board tidak ditemukan", err.Error())
	 }

	 board.InternalID = existingBoard.InternalID
	 board.PublicID = existingBoard.PublicID
	 board.OwnerID = existingBoard.OwnerID
	 board.OwnerPublicID = existingBoard.OwnerPublicID
	 board.CreatedAt = existingBoard.CreatedAt

	 if err := c.service.Update(board); err !=nil {
		return utils.BadRequest(ctx, "Gagal update board", err.Error())
	 }

	 return utils.Success(ctx, "Board berhasil diperbaharui", board)
}

func (c* BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string

	if err := ctx.BodyParser(&userIDs); err !=nil {
		return utils.BadRequest(ctx, "Gagal parsing Data", err.Error())
	}

	if err := c.service.AddMembers(publicID, userIDs) ; err !=nil {
		return utils.BadRequest(ctx, "Gagal menambahkan members", err.Error())
	}

	return utils.Success(ctx, "Member berhasil ditambahkan", nil)
}

func (c* BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string

	if err := ctx.BodyParser(&userIDs); err !=nil {
		return utils.BadRequest(ctx, "Gagal parsing Data", err.Error())
	}

	if err := c.service.RemoveMembers(publicID, userIDs) ; err !=nil {
		return utils.BadRequest(ctx, "Gagal menghapus members", err.Error())
	}

	return utils.Success(ctx, "Members berhasil dihapus", nil)
}

