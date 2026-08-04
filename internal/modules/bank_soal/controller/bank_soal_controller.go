package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/bank_soal/dto"
	"backend/internal/modules/bank_soal/service"

	"github.com/gofiber/fiber/v2"
)

type BankSoalController struct {
	service service.BankSoalService
}

func NewBankSoalController(service service.BankSoalService) *BankSoalController {
	return &BankSoalController{service: service}
}

// CreateBankSoal godoc
// @Summary      Buat bank soal baru
// @Tags         Bank Soal
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateBankSoalRequest  true  "Data bank soal"
// @Success      201   {object}  helpers.Response{data=dto.BankSoalResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /bank-soal [post]
func (c *BankSoalController) CreateBankSoal(ctx *fiber.Ctx) error {
	var req dto.CreateBankSoalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateBankSoal(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create bank soal successfully", resp)
}

// GetAllBankSoal godoc
// @Summary      List bank soal
// @Tags         Bank Soal
// @Produce      json
// @Param        page       query     int  false  "Nomor halaman"          default(1)
// @Param        page_size  query     int  false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.BankSoalListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /bank-soal [get]
func (c *BankSoalController) GetAllBankSoal(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllBankSoal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all bank soal successfully", resp)
}

// GetBankSoalByMapel godoc
// @Summary      List bank soal berdasarkan mapel
// @Tags         Bank Soal
// @Produce      json
// @Param        mapel_id   path      string  true   "ID Mapel"
// @Param        page       query     int     false  "Nomor halaman"          default(1)
// @Param        page_size  query     int     false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.BankSoalListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /bank-soal/mapel/{mapel_id} [get]
func (c *BankSoalController) GetBankSoalByMapel(ctx *fiber.Ctx) error {
	mapelID := ctx.Params("mapel_id")
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetBankSoalByMapel(mapelID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get bank soal by mapel successfully", resp)
}

// GetBankSoalByID godoc
// @Summary      Detail bank soal
// @Tags         Bank Soal
// @Produce      json
// @Param        id   path      string  true  "ID Bank Soal"
// @Success      200  {object}  helpers.Response{data=dto.BankSoalResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /bank-soal/{id} [get]
func (c *BankSoalController) GetBankSoalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetBankSoalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get bank soal successfully", resp)
}

// UpdateBankSoal godoc
// @Summary      Update bank soal
// @Tags         Bank Soal
// @Accept       json
// @Produce      json
// @Param        id    path      string                     true  "ID Bank Soal"
// @Param        body  body      dto.UpdateBankSoalRequest  true  "Data bank soal"
// @Success      200   {object}  helpers.Response{data=dto.BankSoalResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /bank-soal/{id} [put]
func (c *BankSoalController) UpdateBankSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateBankSoalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateBankSoal(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update bank soal successfully", resp)
}

// DeleteBankSoal godoc
// @Summary      Hapus (soft delete) bank soal
// @Tags         Bank Soal
// @Produce      json
// @Param        id   path      string  true  "ID Bank Soal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /bank-soal/{id} [delete]
func (c *BankSoalController) DeleteBankSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteBankSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete bank soal successfully", nil)
}

// RestoreBankSoal godoc
// @Summary      Pulihkan bank soal yang sudah dihapus
// @Tags         Bank Soal
// @Produce      json
// @Param        id   path      string  true  "ID Bank Soal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /bank-soal/{id}/restore [patch]
func (c *BankSoalController) RestoreBankSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreBankSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore bank soal successfully", nil)
}
