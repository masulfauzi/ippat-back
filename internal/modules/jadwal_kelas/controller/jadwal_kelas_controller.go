package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jadwal_kelas/dto"
	"backend/internal/modules/jadwal_kelas/service"

	"github.com/gofiber/fiber/v2"
)

type JadwalKelasController struct {
	service service.JadwalKelasService
}

func NewJadwalKelasController(service service.JadwalKelasService) *JadwalKelasController {
	return &JadwalKelasController{service: service}
}

// CreateJadwalKelas godoc
// @Summary      Kaitkan jadwal ujian dengan kelas
// @Tags         Jadwal Kelas
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateJadwalKelasRequest  true  "Data jadwal-kelas"
// @Success      201   {object}  helpers.Response{data=dto.JadwalKelasResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal-kelas [post]
func (c *JadwalKelasController) CreateJadwalKelas(ctx *fiber.Ctx) error {
	var req dto.CreateJadwalKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJadwalKelas(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jadwal kelas successfully", resp)
}

// GetAllJadwalKelas godoc
// @Summary      List relasi jadwal-kelas
// @Tags         Jadwal Kelas
// @Produce      json
// @Param        page       query     int     false  "Nomor halaman"          default(1)
// @Param        page_size  query     int     false  "Jumlah data per halaman" default(10)
// @Param        id_jadwal  query     string  false  "Filter berdasarkan ID Jadwal"
// @Param        id_kelas   query     string  false  "Filter berdasarkan ID Kelas"
// @Success      200        {object}  helpers.Response{data=dto.JadwalKelasListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /jadwal-kelas [get]
func (c *JadwalKelasController) GetAllJadwalKelas(ctx *fiber.Ctx) error {
	page     := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idJadwal := ctx.Query("id_jadwal", "")
	idKelas  := ctx.Query("id_kelas", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJadwalKelas(pageNum, pageSizeNum, idJadwal, idKelas)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jadwal kelas successfully", resp)
}

// GetJadwalKelasByID godoc
// @Summary      Detail relasi jadwal-kelas
// @Tags         Jadwal Kelas
// @Produce      json
// @Param        id   path      string  true  "ID Jadwal Kelas"
// @Success      200  {object}  helpers.Response{data=dto.JadwalKelasResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /jadwal-kelas/{id} [get]
func (c *JadwalKelasController) GetJadwalKelasByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJadwalKelasByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal kelas successfully", resp)
}

// UpdateJadwalKelas godoc
// @Summary      Update relasi jadwal-kelas
// @Tags         Jadwal Kelas
// @Accept       json
// @Produce      json
// @Param        id    path      string                        true  "ID Jadwal Kelas"
// @Param        body  body      dto.UpdateJadwalKelasRequest  true  "Data jadwal-kelas"
// @Success      200   {object}  helpers.Response{data=dto.JadwalKelasResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal-kelas/{id} [put]
func (c *JadwalKelasController) UpdateJadwalKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJadwalKelasRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJadwalKelas(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jadwal kelas successfully", resp)
}

// DeleteJadwalKelas godoc
// @Summary      Hapus (soft delete) relasi jadwal-kelas
// @Tags         Jadwal Kelas
// @Produce      json
// @Param        id   path      string  true  "ID Jadwal Kelas"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal-kelas/{id} [delete]
func (c *JadwalKelasController) DeleteJadwalKelas(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJadwalKelas(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jadwal kelas successfully", nil)
}
