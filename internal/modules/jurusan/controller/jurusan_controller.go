package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jurusan/dto"
	"backend/internal/modules/jurusan/service"

	"github.com/gofiber/fiber/v2"
)

type JurusanController struct {
	service service.JurusanService
}

func NewJurusanController(service service.JurusanService) *JurusanController {
	return &JurusanController{service: service}
}

// CreateJurusan godoc
// @Summary      Buat jurusan baru
// @Tags         Jurusan
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateJurusanRequest  true  "Data jurusan"
// @Success      201   {object}  helpers.Response{data=dto.JurusanResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jurusan [post]
func (c *JurusanController) CreateJurusan(ctx *fiber.Ctx) error {
	var req dto.CreateJurusanRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJurusan(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jurusan successfully", resp)
}

// GetAllJurusan godoc
// @Summary      List jurusan
// @Tags         Jurusan
// @Produce      json
// @Param        page       query     int  false  "Nomor halaman"          default(1)
// @Param        page_size  query     int  false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.JurusanListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /jurusan [get]
func (c *JurusanController) GetAllJurusan(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllJurusan(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jurusan successfully", resp)
}

// GetJurusanByID godoc
// @Summary      Detail jurusan
// @Tags         Jurusan
// @Produce      json
// @Param        id   path      string  true  "ID Jurusan"
// @Success      200  {object}  helpers.Response{data=dto.JurusanResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /jurusan/{id} [get]
func (c *JurusanController) GetJurusanByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJurusanByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jurusan successfully", resp)
}

// UpdateJurusan godoc
// @Summary      Update jurusan
// @Tags         Jurusan
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "ID Jurusan"
// @Param        body  body      dto.UpdateJurusanRequest  true  "Data jurusan"
// @Success      200   {object}  helpers.Response{data=dto.JurusanResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jurusan/{id} [put]
func (c *JurusanController) UpdateJurusan(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJurusanRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJurusan(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jurusan successfully", resp)
}

// DeleteJurusan godoc
// @Summary      Hapus (soft delete) jurusan
// @Tags         Jurusan
// @Produce      json
// @Param        id   path      string  true  "ID Jurusan"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jurusan/{id} [delete]
func (c *JurusanController) DeleteJurusan(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJurusan(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jurusan successfully", nil)
}

// RestoreJurusan godoc
// @Summary      Pulihkan jurusan yang sudah dihapus
// @Tags         Jurusan
// @Produce      json
// @Param        id   path      string  true  "ID Jurusan"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jurusan/{id}/restore [patch]
func (c *JurusanController) RestoreJurusan(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJurusan(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jurusan successfully", nil)
}
