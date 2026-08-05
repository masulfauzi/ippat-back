package controller

import (
	"strconv"

	"backend/internal/helpers"
	"backend/internal/modules/jawaban/dto"
	"backend/internal/modules/jawaban/service"

	"github.com/gofiber/fiber/v2"
)

type JawabanController struct {
	service service.JawabanService
}

func NewJawabanController(service service.JawabanService) *JawabanController {
	return &JawabanController{service: service}
}

// CreateJawaban godoc
// @Summary      Submit jawaban peserta untuk satu soal
// @Description  Dipanggil peserta saat memilih jawaban pertama kali untuk suatu soal dalam satu attempt ujian (id_nilai). Jika baris jawaban untuk kombinasi id_nilai+id_soal sudah ada (baris ini otomatis dibuat kosong saat ujian dimulai), gunakan PUT /jawaban/{id} untuk mengubahnya.
// @Tags         Jawaban
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateJawabanRequest  true  "Data jawaban"
// @Success      201   {object}  helpers.Response{data=dto.JawabanResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jawaban [post]
func (c *JawabanController) CreateJawaban(ctx *fiber.Ctx) error {
	var req dto.CreateJawabanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.CreateJawaban(&req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusCreated, "Create jawaban successfully", resp)
}

// GetAllJawaban godoc
// @Summary      List jawaban
// @Tags         Jawaban
// @Produce      json
// @Param        page       query     int     false  "Nomor halaman"          default(1)
// @Param        page_size  query     int     false  "Jumlah data per halaman" default(10)
// @Param        id_nilai   query     string  false  "Filter berdasarkan ID Nilai (attempt ujian)"
// @Param        id_peserta query     string  false  "Filter berdasarkan ID Peserta"
// @Param        id_soal    query     string  false  "Filter berdasarkan ID Soal"
// @Success      200        {object}  helpers.Response{data=dto.JawabanListResponse}
// @Failure      500        {object}  helpers.Response
// @Router       /jawaban [get]
func (c *JawabanController) GetAllJawaban(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("page_size", "10")
	idNilai := ctx.Query("id_nilai", "")
	idPeserta := ctx.Query("id_peserta", "")
	idSoal := ctx.Query("id_soal", "")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum <= 0 {
		pageSizeNum = 10
	}

	resp, err := c.service.GetAllJawaban(pageNum, pageSizeNum, idNilai, idPeserta, idSoal)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get all jawaban successfully", resp)
}

// GetJawabanByID godoc
// @Summary      Detail jawaban
// @Tags         Jawaban
// @Produce      json
// @Param        id   path      string  true  "ID Jawaban"
// @Success      200  {object}  helpers.Response{data=dto.JawabanResponse}
// @Failure      404  {object}  helpers.Response
// @Router       /jawaban/{id} [get]
func (c *JawabanController) GetJawabanByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.GetJawabanByID(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusNotFound, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban successfully", resp)
}

// GetJawabanByNilai godoc
// @Summary      List jawaban dalam satu attempt ujian
// @Description  Mengembalikan seluruh baris jawaban (termasuk yang masih kosong) untuk satu attempt ujian (id_nilai), terurut sesuai no_urut — dipakai frontend untuk render soal per soal dan navigasi antar soal.
// @Tags         Jawaban
// @Produce      json
// @Param        id_nilai  path      string  true  "ID Nilai (attempt ujian)"
// @Success      200       {object}  helpers.Response{data=[]dto.JawabanResponse}
// @Failure      500       {object}  helpers.Response
// @Router       /jawaban/nilai/{id_nilai} [get]
func (c *JawabanController) GetJawabanByNilai(ctx *fiber.Ctx) error {
	idNilai := ctx.Params("id_nilai")

	resp, err := c.service.GetJawabanByNilai(idNilai)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by nilai successfully", resp)
}

// GetJawabanByPeserta godoc
// @Summary      List jawaban milik seorang peserta
// @Tags         Jawaban
// @Produce      json
// @Param        id_peserta  path      string  true   "ID Peserta"
// @Param        page        query     int     false  "Nomor halaman"          default(1)
// @Param        page_size   query     int     false  "Jumlah data per halaman" default(10)
// @Success      200         {object}  helpers.Response{data=dto.JawabanListResponse}
// @Failure      500         {object}  helpers.Response
// @Router       /jawaban/peserta/{id_peserta} [get]
func (c *JawabanController) GetJawabanByPeserta(ctx *fiber.Ctx) error {
	idPeserta := ctx.Params("id_peserta")
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

	resp, err := c.service.GetJawabanByPeserta(idPeserta, pageNum, pageSizeNum)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Get jawaban by peserta successfully", resp)
}

// UpdateJawaban godoc
// @Summary      Ubah pilihan jawaban peserta
// @Description  Dipanggil saat peserta mengganti pilihan jawaban untuk soal yang sama (baris jawaban sudah ada — lihat GET /jawaban/nilai/{id_nilai} untuk mendapatkan id-nya). Jawaban dievaluasi ulang (is_benar) terhadap kunci soal.
// @Tags         Jawaban
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "ID Jawaban"
// @Param        body  body      dto.UpdateJawabanRequest  true  "Pilihan jawaban baru (A/B/C/D/E)"
// @Success      200   {object}  helpers.Response{data=dto.JawabanResponse}
// @Failure      400   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jawaban/{id} [put]
func (c *JawabanController) UpdateJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req dto.UpdateJawabanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request format", nil)
	}

	resp, err := c.service.UpdateJawaban(id, &req)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Update jawaban successfully", resp)
}

// ClearJawaban godoc
// @Summary      Kosongkan pilihan jawaban (batalkan jawaban)
// @Description  Mengosongkan isi jawaban (jawaban & is_benar diset null) TANPA menghapus baris soal dari attempt ujian — soal tetap muncul di daftar ujian sebagai "belum dijawab". Dipakai saat peserta ragu dengan jawabannya dan ingin membatalkannya agar tidak kena penalti poin salah dari kategori soal (soal kosong dihitung 0). Peserta tetap bisa memilih jawaban lain lagi lewat PUT /jawaban/{id}.
// @Tags         Jawaban
// @Produce      json
// @Param        id   path      string  true  "ID Jawaban"
// @Success      200  {object}  helpers.Response{data=dto.JawabanResponse}
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jawaban/{id}/jawaban [delete]
func (c *JawabanController) ClearJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	resp, err := c.service.ClearJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Jawaban berhasil dikosongkan", resp)
}

// DeleteJawaban godoc
// @Summary      Hapus (soft delete) baris jawaban
// @Description  Endpoint administratif: menghapus seluruh baris jawaban (soal ikut hilang dari daftar attempt ujian tersebut). Untuk membatalkan pilihan jawaban peserta tanpa menghilangkan soalnya, gunakan DELETE /jawaban/{id}/jawaban.
// @Tags         Jawaban
// @Produce      json
// @Param        id   path      string  true  "ID Jawaban"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jawaban/{id} [delete]
func (c *JawabanController) DeleteJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.DeleteJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Delete jawaban successfully", nil)
}

// RestoreJawaban godoc
// @Summary      Pulihkan baris jawaban yang sudah dihapus
// @Tags         Jawaban
// @Produce      json
// @Param        id   path      string  true  "ID Jawaban"
// @Success      200  {object}  helpers.Response
// @Failure      400  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /jawaban/{id}/restore [patch]
func (c *JawabanController) RestoreJawaban(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.service.RestoreJawaban(id)
	if err != nil {
		return helpers.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error(), nil)
	}

	return helpers.SuccessResponse(ctx, fiber.StatusOK, "Restore jawaban successfully", nil)
}
