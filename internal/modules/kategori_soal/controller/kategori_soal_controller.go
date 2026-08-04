package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/kategori_soal/dto"
	"backend/internal/modules/kategori_soal/service"

	"github.com/gofiber/fiber/v2"
)

type KategoriSoalController struct {
	service service.KategoriSoalService
}

func NewKategoriSoalController(service service.KategoriSoalService) *KategoriSoalController {
	return &KategoriSoalController{service: service}
}

// CreateKategoriSoal godoc
// @Summary      Buat kategori soal baru
// @Tags         Kategori Soal
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateKategoriSoalRequest  true  "Data kategori soal"
// @Success      201   {object}  helpers.Response{data=dto.KategoriSoalResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /kategori-soal [post]
func (c *KategoriSoalController) CreateKategoriSoal(ctx *fiber.Ctx) error {
	var req dto.CreateKategoriSoalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateKategoriSoal(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create kategori soal successfully", resp)
}

// GetAllKategoriSoal godoc
// @Summary      List kategori soal
// @Tags         Kategori Soal
// @Produce      json
// @Param        page       query     int  false  "Nomor halaman"          default(1)
// @Param        page_size  query     int  false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.KategoriSoalListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /kategori-soal [get]
func (c *KategoriSoalController) GetAllKategoriSoal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllKategoriSoal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all kategori soal successfully", resp)
}

// GetKategoriSoalByID godoc
// @Summary      Detail kategori soal
// @Tags         Kategori Soal
// @Produce      json
// @Param        id   path      string  true  "ID Kategori Soal"
// @Success      200  {object}  helpers.Response{data=dto.KategoriSoalResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /kategori-soal/{id} [get]
func (c *KategoriSoalController) GetKategoriSoalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetKategoriSoalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get kategori soal successfully", resp)
}

// UpdateKategoriSoal godoc
// @Summary      Update kategori soal
// @Tags         Kategori Soal
// @Accept       json
// @Produce      json
// @Param        id    path      string                          true  "ID Kategori Soal"
// @Param        body  body      dto.UpdateKategoriSoalRequest  true  "Data kategori soal"
// @Success      200   {object}  helpers.Response{data=dto.KategoriSoalResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /kategori-soal/{id} [put]
func (c *KategoriSoalController) UpdateKategoriSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateKategoriSoalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateKategoriSoal(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update kategori soal successfully", resp)
}

// DeleteKategoriSoal godoc
// @Summary      Hapus (soft delete) kategori soal
// @Tags         Kategori Soal
// @Produce      json
// @Param        id   path      string  true  "ID Kategori Soal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /kategori-soal/{id} [delete]
func (c *KategoriSoalController) DeleteKategoriSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteKategoriSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete kategori soal successfully", nil)
}

// RestoreKategoriSoal godoc
// @Summary      Pulihkan kategori soal yang sudah dihapus
// @Tags         Kategori Soal
// @Produce      json
// @Param        id   path      string  true  "ID Kategori Soal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /kategori-soal/{id}/restore [patch]
func (c *KategoriSoalController) RestoreKategoriSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreKategoriSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore kategori soal successfully", nil)
}
