package controller

import (
	"path/filepath"
	"strconv"

	"backend/internal/helpers"
	jadwalrepo "backend/internal/modules/jadwal/repository"
	"backend/internal/modules/soal/dto"
	"backend/internal/modules/soal/service"

	"github.com/gofiber/fiber/v2"
)

type SoalController struct {
	service    service.SoalService
	jadwalRepo jadwalrepo.JadwalRepository
}

func NewSoalController(service service.SoalService, jadwalRepo jadwalrepo.JadwalRepository) *SoalController {
	return &SoalController{service: service, jadwalRepo: jadwalRepo}
}

// CreateSoal godoc
// @Summary      Buat soal baru
// @Description  Form-data (bukan JSON) karena mendukung upload gambar untuk soal maupun tiap opsi jawaban.
// @Tags         Soal
// @Accept       multipart/form-data
// @Produce      json
// @Param        id_bank_soal      formData  string  true   "ID Bank Soal"
// @Param        id_kategori_soal  formData  string  false  "ID Kategori Soal"
// @Param        no_soal       formData  int     false  "Nomor urut soal"
// @Param        soal          formData  string  false  "Teks soal"
// @Param        gambar_soal   formData  file    false  "Gambar soal"
// @Param        opsi_a        formData  string  false  "Teks opsi A"
// @Param        opsi_b        formData  string  false  "Teks opsi B"
// @Param        opsi_c        formData  string  false  "Teks opsi C"
// @Param        opsi_d        formData  string  false  "Teks opsi D"
// @Param        opsi_e        formData  string  false  "Teks opsi E"
// @Param        gambar_a      formData  file    false  "Gambar opsi A"
// @Param        gambar_b      formData  file    false  "Gambar opsi B"
// @Param        gambar_c      formData  file    false  "Gambar opsi C"
// @Param        gambar_d      formData  file    false  "Gambar opsi D"
// @Param        gambar_e      formData  file    false  "Gambar opsi E"
// @Param        kunci         formData  string  true   "Kunci jawaban (A/B/C/D/E)"
// @Success      201  {object}  helpers.Response{data=dto.SoalResponse}
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /soal [post]
func (c *SoalController) CreateSoal(ctx *fiber.Ctx) error {
	req := new(dto.CreateSoalRequest)

	req.IdBankSoal = ctx.FormValue("id_bank_soal")
	req.IdKategoriSoal = ctx.FormValue("id_kategori_soal")
	req.Soal = ctx.FormValue("soal")
	req.OpsiA = ctx.FormValue("opsi_a")
	req.OpsiB = ctx.FormValue("opsi_b")
	req.OpsiC = ctx.FormValue("opsi_c")
	req.OpsiD = ctx.FormValue("opsi_d")
	req.OpsiE = ctx.FormValue("opsi_e")
	req.Kunci = ctx.FormValue("kunci")

	noSoalStr := ctx.FormValue("no_soal")
	if noSoal, err := strconv.Atoi(noSoalStr); err == nil {
		req.NoSoal = noSoal
	}

	if file, err := ctx.FormFile("gambar_soal"); err == nil {
		req.GambarSoal = file
	}
	if file, err := ctx.FormFile("gambar_a"); err == nil {
		req.GambarA = file
	}
	if file, err := ctx.FormFile("gambar_b"); err == nil {
		req.GambarB = file
	}
	if file, err := ctx.FormFile("gambar_c"); err == nil {
		req.GambarC = file
	}
	if file, err := ctx.FormFile("gambar_d"); err == nil {
		req.GambarD = file
	}
	if file, err := ctx.FormFile("gambar_e"); err == nil {
		req.GambarE = file
	}

	resp, err := c.service.CreateSoal(req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create soal successfully", resp)
}

// GetSoalByID godoc
// @Summary      Detail soal
// @Description  Jika request membawa token JWT peserta yang sedang ujian dan jadwalnya diset acak_opsi=1, urutan opsi jawaban akan diacak secara konsisten untuk peserta tersebut (bukan diambil dari urutan asli di database).
// @Tags         Soal
// @Produce      json
// @Param        id   path      string  true  "ID Soal"
// @Success      200  {object}  helpers.Response{data=dto.SoalResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /soal/{id} [get]
func (c *SoalController) GetSoalByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetSoalByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	// Auto-acak opsi jika request datang dari peserta yang sedang ujian.
	// Backend cari nilai aktif peserta (via JWT user_id) → ambil acak_opsi dari jadwal-nya.
	if userID, ok := ctx.Locals("user_id").(string); ok && userID != "" {
		acakOpsi, err := c.jadwalRepo.GetAcakOpsiForPesertaSoal(userID, id)
		if err == nil && acakOpsi == 1 {
			resp = c.service.RandomizeOpsiForSoal(resp, userID)
		}
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get soal successfully", resp)
}

// GetAllSoal godoc
// @Summary      List soal
// @Tags         Soal
// @Produce      json
// @Param        page       query     int  false  "Nomor halaman"          default(1)
// @Param        page_size  query     int  false  "Jumlah data per halaman" default(10)
// @Success      200        {object}  helpers.Response{data=dto.SoalListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /soal [get]
func (c *SoalController) GetAllSoal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetAllSoal(pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all soal successfully", resp)
}

// GetSoalByBankSoal godoc
// @Summary      List soal berdasarkan bank soal
// @Tags         Soal
// @Produce      json
// @Param        bank_soal_id  path      string  true   "ID Bank Soal"
// @Param        page          query     int     false  "Nomor halaman"          default(1)
// @Param        page_size     query     int     false  "Jumlah data per halaman" default(10)
// @Success      200           {object}  helpers.Response{data=dto.SoalListResponse}
// @Failure      500           {object}  helpers.Response
// @Router       /soal/bank/{bank_soal_id} [get]
func (c *SoalController) GetSoalByBankSoal(ctx *fiber.Ctx) error {
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

	resp, err := c.service.GetSoalByBankSoal(bankSoalID, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get soal by bank successfully", resp)
}

// UpdateSoal godoc
// @Summary      Update soal
// @Description  Form-data (bukan JSON) karena mendukung upload gambar untuk soal maupun tiap opsi jawaban.
// @Tags         Soal
// @Accept       multipart/form-data
// @Produce      json
// @Param        id               path      string  true   "ID Soal"
// @Param        id_kategori_soal formData  string  false  "ID Kategori Soal"
// @Param        no_soal      formData  int     false  "Nomor urut soal"
// @Param        soal         formData  string  false  "Teks soal"
// @Param        gambar_soal  formData  file    false  "Gambar soal"
// @Param        opsi_a       formData  string  false  "Teks opsi A"
// @Param        opsi_b       formData  string  false  "Teks opsi B"
// @Param        opsi_c       formData  string  false  "Teks opsi C"
// @Param        opsi_d       formData  string  false  "Teks opsi D"
// @Param        opsi_e       formData  string  false  "Teks opsi E"
// @Param        gambar_a     formData  file    false  "Gambar opsi A"
// @Param        gambar_b     formData  file    false  "Gambar opsi B"
// @Param        gambar_c     formData  file    false  "Gambar opsi C"
// @Param        gambar_d     formData  file    false  "Gambar opsi D"
// @Param        gambar_e     formData  file    false  "Gambar opsi E"
// @Param        kunci        formData  string  true   "Kunci jawaban (A/B/C/D/E)"
// @Success      200  {object}  helpers.Response{data=dto.SoalResponse}
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /soal/{id} [put]
func (c *SoalController) UpdateSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	req := new(dto.UpdateSoalRequest)

	req.IdKategoriSoal = ctx.FormValue("id_kategori_soal")
	req.Soal = ctx.FormValue("soal")
	req.OpsiA = ctx.FormValue("opsi_a")
	req.OpsiB = ctx.FormValue("opsi_b")
	req.OpsiC = ctx.FormValue("opsi_c")
	req.OpsiD = ctx.FormValue("opsi_d")
	req.OpsiE = ctx.FormValue("opsi_e")
	req.Kunci = ctx.FormValue("kunci")

	noSoalStr := ctx.FormValue("no_soal")
	if noSoal, err := strconv.Atoi(noSoalStr); err == nil {
		req.NoSoal = noSoal
	}

	if file, err := ctx.FormFile("gambar_soal"); err == nil {
		req.GambarSoal = file
	}
	if file, err := ctx.FormFile("gambar_a"); err == nil {
		req.GambarA = file
	}
	if file, err := ctx.FormFile("gambar_b"); err == nil {
		req.GambarB = file
	}
	if file, err := ctx.FormFile("gambar_c"); err == nil {
		req.GambarC = file
	}
	if file, err := ctx.FormFile("gambar_d"); err == nil {
		req.GambarD = file
	}
	if file, err := ctx.FormFile("gambar_e"); err == nil {
		req.GambarE = file
	}

	resp, err := c.service.UpdateSoal(id, req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update soal successfully", resp)
}

// DeleteSoal godoc
// @Summary      Hapus (soft delete) soal
// @Tags         Soal
// @Produce      json
// @Param        id   path      string  true  "ID Soal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /soal/{id} [delete]
func (c *SoalController) DeleteSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete soal successfully", nil)
}

// RestoreSoal godoc
// @Summary      Pulihkan soal yang sudah dihapus
// @Tags         Soal
// @Produce      json
// @Param        id   path      string  true  "ID Soal"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /soal/{id}/restore [patch]
func (c *SoalController) RestoreSoal(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreSoal(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore soal successfully", nil)
}

// ImportSoalFromExcel godoc
// @Summary      Import soal massal dari file Excel
// @Description  Upload file .xls/.xlsx (maks 10MB) berisi banyak soal sekaligus untuk satu bank soal. Response berisi ringkasan jumlah berhasil/gagal beserta detail error per baris.
// @Tags         Soal
// @Accept       multipart/form-data
// @Produce      json
// @Param        id_bank_soal  formData  string  true  "ID Bank Soal"
// @Param        file          formData  file    true  "File Excel (.xls/.xlsx, maks 10MB)"
// @Success      200  {object}  helpers.Response{data=dto.ImportSoalResponse}
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /soal/import [post]
func (c *SoalController) ImportSoalFromExcel(ctx *fiber.Ctx) error {
	// 1. Parse multipart form
	file, err := ctx.FormFile("file")
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File tidak ditemukan", map[string]string{
			"error": "Silakan upload file excel",
		})
	}

	// 2. Validate file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "File terlalu besar", map[string]string{
			"error": "Max file size adalah 10MB",
		})
	}

	// 3. Validate file extension
	ext := filepath.Ext(file.Filename)
	if ext != ".xls" && ext != ".xlsx" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Format file tidak valid", map[string]string{
			"error": "File harus berupa .xls atau .xlsx",
		})
	}

	// 4. Build request
	req := &dto.ImportSoalRequest{
		IdBankSoal: ctx.FormValue("id_bank_soal"),
		File:       file,
	}

	// 5. Validate required fields
	if req.IdBankSoal == "" {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "id_bank_soal tidak ditemukan", nil)
	}

	// 6. Call service
	resp, err := c.service.ImportSoalFromExcel(ctx.Context(), req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Import soal gagal", map[string]string{
			"error": err.Error(),
		})
	}

	// 7. Return success response
	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Import soal berhasil", resp)
}
