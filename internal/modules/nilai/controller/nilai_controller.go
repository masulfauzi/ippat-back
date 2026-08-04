package controller

import (
	"fmt"
	"strconv"
	"strings"

	"backend/internal/helpers"
	"backend/internal/modules/nilai/dto"
	"backend/internal/modules/nilai/service"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type NilaiController struct {
	service service.NilaiService
}

func NewNilaiController(service service.NilaiService) *NilaiController {
	return &NilaiController{service: service}
}

// CreateNilai godoc
// @Summary      Buat data nilai secara manual
// @Description  Umumnya nilai/attempt ujian dibuat otomatis lewat POST /nilai/mulai-ujian/{id_jadwal}. Endpoint ini untuk input manual oleh admin/guru.
// @Tags         Nilai
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateNilaiRequest  true  "Data nilai"
// @Success      201   {object}  helpers.Response{data=dto.NilaiResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /nilai [post]
func (c *NilaiController) CreateNilai(ctx *fiber.Ctx) error {
	var req dto.CreateNilaiRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateNilai(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create nilai successfully", resp)
}

// GetAllNilai godoc
// @Summary      List nilai/attempt ujian
// @Tags         Nilai
// @Produce      json
// @Param        page       query     int     false  "Nomor halaman"          default(1)
// @Param        page_size  query     int     false  "Jumlah data per halaman" default(10)
// @Param        id_peserta query     string  false  "Filter berdasarkan ID Peserta"
// @Param        id_jadwal  query     string  false  "Filter berdasarkan ID Jadwal"
// @Success      200        {object}  helpers.Response{data=dto.NilaiListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /nilai [get]
func (c *NilaiController) GetAllNilai(ctx *fiber.Ctx) error {
	page      := ctx.Query("page", "1")
	pageSize  := ctx.Query("page_size", "10")
	idPeserta := ctx.Query("id_peserta", "")
	idJadwal  := ctx.Query("id_jadwal", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllNilai(pageNum, pageSizeNum, idPeserta, idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all nilai successfully", resp)
}

// GetNilaiByID godoc
// @Summary      Detail nilai/attempt ujian
// @Tags         Nilai
// @Produce      json
// @Param        id   path      string  true  "ID Nilai"
// @Success      200  {object}  helpers.Response{data=dto.NilaiResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /nilai/{id} [get]
func (c *NilaiController) GetNilaiByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetNilaiByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai successfully", resp)
}

// GetNilaiByPeserta godoc
// @Summary      List riwayat nilai seorang peserta
// @Tags         Nilai
// @Produce      json
// @Param        id_peserta  path      string  true   "ID Peserta"
// @Param        page        query     int     false  "Nomor halaman"          default(1)
// @Param        page_size   query     int     false  "Jumlah data per halaman" default(10)
// @Success      200         {object}  helpers.Response{data=dto.NilaiListResponse}
// @Failure      500         {object}  helpers.Response
// @Router       /nilai/peserta/{id_peserta} [get]
func (c *NilaiController) GetNilaiByPeserta(ctx *fiber.Ctx) error {
	idPeserta := ctx.Params("id_peserta")
	page      := ctx.Query("page", "1")
	pageSize  := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetNilaiByPeserta(idPeserta, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai by peserta successfully", resp)
}

// GetNilaiByJadwal godoc
// @Summary      List nilai seluruh peserta untuk satu jadwal ujian
// @Tags         Nilai
// @Produce      json
// @Param        id_jadwal  path      string  true   "ID Jadwal"
// @Param        page       query     int     false  "Nomor halaman"          default(1)
// @Param        page_size  query     int     false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.NilaiListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /nilai/jadwal/{id_jadwal} [get]
func (c *NilaiController) GetNilaiByJadwal(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	page     := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetNilaiByJadwal(idJadwal, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get nilai by jadwal successfully", resp)
}

// UpdateNilai godoc
// @Summary      Update data nilai
// @Tags         Nilai
// @Accept       json
// @Produce      json
// @Param        id    path      string                  true  "ID Nilai"
// @Param        body  body      dto.UpdateNilaiRequest  true  "Data nilai"
// @Success      200   {object}  helpers.Response{data=dto.NilaiResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /nilai/{id} [put]
func (c *NilaiController) UpdateNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateNilaiRequest

	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateNilai(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update nilai successfully", resp)
}

// DeleteNilai godoc
// @Summary      Hapus (soft delete) data nilai
// @Tags         Nilai
// @Produce      json
// @Param        id   path      string  true  "ID Nilai"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /nilai/{id} [delete]
func (c *NilaiController) DeleteNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteNilai(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete nilai successfully", nil)
}

// RestoreNilai godoc
// @Summary      Pulihkan data nilai yang sudah dihapus
// @Tags         Nilai
// @Produce      json
// @Param        id   path      string  true  "ID Nilai"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /nilai/{id}/restore [patch]
func (c *NilaiController) RestoreNilai(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreNilai(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore nilai successfully", nil)
}

// MulaiUjian godoc
// @Summary      Mulai atau lanjutkan ujian
// @Description  Dipanggil peserta saat membuka ujian. Jika belum ada attempt, membuat data nilai baru + baris jawaban kosong untuk tiap soal (status 201). Jika attempt sudah ada dan belum selesai, mengembalikan data attempt yang sama untuk dilanjutkan (status 200). Mengembalikan error jika ujian untuk jadwal ini sudah pernah diselesaikan. id_peserta diambil dari token JWT peserta yang login.
// @Tags         Nilai
// @Produce      json
// @Param        id_jadwal  path      string  true  "ID Jadwal"
// @Success      200        {object}  helpers.Response{data=dto.NilaiResponse}  "Melanjutkan ujian yang sudah berjalan"
// @Success      201        {object}  helpers.Response{data=dto.NilaiResponse}  "Ujian baru dimulai"
// @Failure      400         {object}  helpers.Response
// @Failure      401         {object}  helpers.Response
// @Security     BearerAuth
// @Router       /nilai/mulai-ujian/{id_jadwal} [post]
func (c *NilaiController) MulaiUjian(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	if idJadwal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
	}

	// Ambil user_id (= id_peserta) dari JWT
	userToken, ok := ctx.Locals("user").(*jwt.Token)
	if !ok || userToken == nil {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "Invalid token claims", nil)
	}
	idPeserta, ok := claims["user_id"].(string)
	if !ok || idPeserta == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusUnauthorized, "user_id tidak ditemukan di token", nil)
	}

	resp, isNew, err := c.service.MulaiUjian(idPeserta, idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	if isNew {
		return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Mulai ujian successfully", resp)
	}
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Lanjutkan ujian successfully", resp)
}

// ExportNilai godoc
// @Summary      Export nilai satu jadwal ujian ke file (ZIP)
// @Description  Menghasilkan file ZIP berisi rekap nilai seluruh peserta untuk satu jadwal ujian. Response bukan JSON, melainkan file binary (Content-Type application/zip) yang langsung ter-download.
// @Tags         Nilai
// @Produce      application/zip
// @Param        id_jadwal  path      string  true  "ID Jadwal"
// @Success      200        {file}    file    "File ZIP hasil export"
// @Failure      400        {object}  helpers.Response
// @Security     BearerAuth
// @Router       /nilai/export/{id_jadwal} [get]
func (c *NilaiController) ExportNilai(ctx *fiber.Ctx) error {
	idJadwal := ctx.Params("id_jadwal")
	if idJadwal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_jadwal tidak boleh kosong", nil)
	}

	result, err := c.service.ExportNilaiByJadwal(idJadwal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	safeNama := strings.ReplaceAll(result.NamaUjian, " ", "_")
	filename := fmt.Sprintf("export_nilai_%s.zip", safeNama)

	ctx.Set("Content-Type", "application/zip")
	ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return ctx.Send(result.ZipBytes)
}
