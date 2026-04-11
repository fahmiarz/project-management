package controllers

import (
	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/services"
	"github.com/fahmiarz/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController{
	return &ListController{service: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error{
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil{
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}
	if err := c.service.Create(list); err !=nil {
		return utils.BadRequest(ctx, "Gagal membuat list", err.Error())
	}
	return utils.Success(ctx, "List berhasil dibuat", list)
}

func (c *ListController) UpdateList (ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new (models.List)

	//gabungkan apa yang ada di body parameter ke dalam variable list
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	//validasi publid id
	if _,err := uuid.Parse(publicID) ; err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err. Error())
	}

	//verifikasi list yang di update ada atau tidak
	existingList , err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditermukan", err.Error())
	}

	//gabungkan apa yang ada di variable existingList ke dalam variable list
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.Update(list); err !=nil {
		return utils.BadRequest(ctx, "Gagal update list", err.Error())
	}

	//mengambil data yang sudah di update
	updatedList , err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	return utils.Success(ctx , "Berhasil memperbarui list" , updatedList)

}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error{
	boardPublicID := ctx.Params("board_id")
	//validasi boardPublicID 
	if _,err := uuid.Parse(boardPublicID ) ; err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err. Error())
	}

	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx , "List tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "Data berhasil diambil" , lists)
}

func (c * ListController) DeleteList (ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	//validasi boardPublicID 
	if _,err := uuid.Parse(publicID ) ; err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err. Error())
	}

	lists, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx , "List tidak ditemukan" , err.Error())
	}

	if err := c.service.Delete(uint(lists.InternalID)) ; err !=nil {
		return utils.InternalServerError(ctx , "Gagal menghapus list", err.Error())
	}

	return utils.Success(ctx, "List berhasil dihapus", publicID)
}