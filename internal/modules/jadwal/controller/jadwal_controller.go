package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jadwal/dto"
	"backend/internal/modules/jadwal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JadwalController struct {
	service service.JadwalService
}

func NewJadwalController(service service.JadwalService) *JadwalController {
	return &JadwalController{service: service}
}

// CreateJadwal godoc
// @Summary      Buat jadwal ujian baru
// @Tags         Jadwal
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateJadwalRequest  true  "Data jadwal"
// @Success      201   {object}  helpers.Response{data=dto.JadwalResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal [post]
func (c *JadwalController) CreateJadwal(ctx *fiber.Ctx) error {
	var req dto.CreateJadwalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJadwal(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jadwal successfully", resp)
}

// GetAllJadwal godoc
// @Summary      List jadwal ujian
// @Tags         Jadwal
// @Produce      json
// @Param        page       query     int  false  "Nomor halaman"          default(1)
// @Param        page_size  query     int  false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.JadwalListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /jadwal [get]
func (c *JadwalController) GetAllJadwal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllJadwal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jadwal successfully", resp)
}

// GetJadwalByBankSoal godoc
// @Summary      List jadwal berdasarkan bank soal
// @Tags         Jadwal
// @Produce      json
// @Param        bank_soal_id  path      string  true   "ID Bank Soal"
// @Param        page          query     int     false  "Nomor halaman"          default(1)
// @Param        page_size     query     int     false  "Jumlah data per halaman" default(10)
// @Success      200           {object}  helpers.Response{data=dto.JadwalListResponse}
// @Failure      500           {object}  helpers.Response
// @Router       /jadwal/bank-soal/{bank_soal_id} [get]
func (c *JadwalController) GetJadwalByBankSoal(ctx *fiber.Ctx) error {
	bankSoalID := ctx.Params("bank_soal_id")
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

	resp, err := c.service.GetJadwalByBankSoal(bankSoalID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal by bank soal successfully", resp)
}

// GetJadwalByID godoc
// @Summary      Detail jadwal ujian
// @Tags         Jadwal
// @Produce      json
// @Param        id   path      string  true  "ID Jadwal"
// @Success      200  {object}  helpers.Response{data=dto.JadwalResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /jadwal/{id} [get]
func (c *JadwalController) GetJadwalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJadwalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal successfully", resp)
}

// UpdateJadwal godoc
// @Summary      Update jadwal ujian
// @Tags         Jadwal
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "ID Jadwal"
// @Param        body  body      dto.UpdateJadwalRequest  true  "Data jadwal"
// @Success      200   {object}  helpers.Response{data=dto.JadwalResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal/{id} [put]
func (c *JadwalController) UpdateJadwal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJadwalRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJadwal(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jadwal successfully", resp)
}

// DeleteJadwal godoc
// @Summary      Hapus (soft delete) jadwal ujian
// @Tags         Jadwal
// @Produce      json
// @Param        id   path      string  true  "ID Jadwal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal/{id} [delete]
func (c *JadwalController) DeleteJadwal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJadwal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jadwal successfully", nil)
}

// RestoreJadwal godoc
// @Summary      Pulihkan jadwal yang sudah dihapus
// @Tags         Jadwal
// @Produce      json
// @Param        id   path      string  true  "ID Jadwal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal/{id}/restore [patch]
func (c *JadwalController) RestoreJadwal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJadwal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jadwal successfully", nil)
}

// GetJadwalAktifHariIni godoc
// @Summary      Jadwal ujian aktif hari ini untuk peserta yang login
// @Description  Mengambil daftar jadwal ujian yang berlangsung hari ini untuk kelas peserta yang sedang login, lengkap dengan status pengerjaannya (belum mulai/sedang ujian/selesai).
// @Tags         Jadwal
// @Produce      json
// @Success      200  {object}  helpers.Response{data=[]dto.JadwalAktifResponse}
// @Failure      401  {object}  helpers.Response
// @Failure      500  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jadwal/aktif/hari-ini [get]
func (c *JadwalController) GetJadwalAktifHariIni(ctx *fiber.Ctx) error {
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid token claims", nil)
	}

	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "user_id tidak ditemukan di token", nil)
	}

	resp, err := c.service.GetJadwalAktifHariIniByUser(userID)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jadwal aktif hari ini successfully", resp)
}
