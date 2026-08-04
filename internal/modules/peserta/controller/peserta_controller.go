package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/service"

	"github.com/gofiber/fiber/v2"
)

type PesertaController struct {
	service service.PesertaService
}

func NewPesertaController(service service.PesertaService) *PesertaController {
	return &PesertaController{service: service}
}

// CreatePeserta godoc
// @Summary      Buat akun peserta baru
// @Tags         Peserta
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreatePesertaRequest  true  "Data peserta"
// @Success      201   {object}  helpers.Response{data=dto.PesertaResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /peserta [post]
func (c *PesertaController) CreatePeserta(ctx *fiber.Ctx) error {
	var req dto.CreatePesertaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreatePeserta(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create peserta successfully", resp)
}

// GetAllPeserta godoc
// @Summary      List peserta
// @Tags         Peserta
// @Produce      json
// @Param        page       query     int     false  "Nomor halaman"          default(1)
// @Param        page_size  query     int     false  "Jumlah data per halaman" default(10)
// @Param        id_kelas   query     string  false  "Filter berdasarkan ID Kelas"
// @Success      200        {object}  helpers.Response{data=dto.PesertaListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /peserta [get]
func (c *PesertaController) GetAllPeserta(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idKelas := ctx.Query("id_kelas", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllPeserta(pageNum, pageSizeNum, idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all peserta successfully", resp)
}

// GetPesertaByID godoc
// @Summary      Detail peserta
// @Tags         Peserta
// @Produce      json
// @Param        id   path      string  true  "ID Peserta"
// @Success      200  {object}  helpers.Response{data=dto.PesertaResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /peserta/{id} [get]
func (c *PesertaController) GetPesertaByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetPesertaByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get peserta successfully", resp)
}

// UpdatePeserta godoc
// @Summary      Update data peserta
// @Tags         Peserta
// @Accept       json
// @Produce      json
// @Param        id    path      string                    true  "ID Peserta"
// @Param        body  body      dto.UpdatePesertaRequest  true  "Data peserta"
// @Success      200   {object}  helpers.Response{data=dto.PesertaResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /peserta/{id} [put]
func (c *PesertaController) UpdatePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdatePesertaRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdatePeserta(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update peserta successfully", resp)
}

// DeletePeserta godoc
// @Summary      Hapus (soft delete) peserta
// @Tags         Peserta
// @Produce      json
// @Param        id   path      string  true  "ID Peserta"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /peserta/{id} [delete]
func (c *PesertaController) DeletePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeletePeserta(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete peserta successfully", nil)
}

// RestorePeserta godoc
// @Summary      Pulihkan peserta yang sudah dihapus
// @Tags         Peserta
// @Produce      json
// @Param        id   path      string  true  "ID Peserta"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /peserta/{id}/restore [patch]
func (c *PesertaController) RestorePeserta(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestorePeserta(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore peserta successfully", nil)
}
