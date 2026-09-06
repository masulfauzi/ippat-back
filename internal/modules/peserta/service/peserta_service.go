package service

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"backend/internal/constants"
	"backend/internal/modules/peserta/dto"
	"backend/internal/modules/peserta/model"
	"backend/internal/modules/peserta/repository"
	"backend/internal/utils"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func pesertaWithKelasToResponse(p *repository.PesertaWithKelas) *dto.PesertaResponse {
	return &dto.PesertaResponse{
		ID:        p.ID,
		Nama:      p.Nama,
		IDKelas:   p.IDKelas,
		NamaKelas: p.NamaKelas,
		Username:  p.Username,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type PesertaService interface {
	CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error)
	GetPesertaByID(id string) (*dto.PesertaResponse, error)
	GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error)
	UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error)
	DeletePeserta(id string) error
	RestorePeserta(id string) error
	ImportPesertaFromExcel(ctx context.Context, req *dto.ImportPesertaRequest) (*dto.ImportPesertaResponse, error)
}

type pesertaService struct {
	repo repository.PesertaRepository
}

func NewPesertaService(repo repository.PesertaRepository) PesertaService {
	return &pesertaService{repo: repo}
}

func (s *pesertaService) CreatePeserta(req *dto.CreatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetByUsername(req.Username)
	if err == nil && existing != nil {
		return nil, errors.New("username sudah digunakan")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	peserta := &model.Peserta{
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.repo.Create(peserta); err != nil {
		return nil, err
	}

	created, err := s.repo.GetByID(peserta.ID)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(created), nil
}

func (s *pesertaService) GetPesertaByID(id string) (*dto.PesertaResponse, error) {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}
	return pesertaWithKelasToResponse(peserta), nil
}

func (s *pesertaService) GetAllPeserta(page, pageSize int, idKelas string) (*dto.PesertaListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	pesertaList, total, err := s.repo.GetAll(page, pageSize, idKelas)
	if err != nil {
		return nil, err
	}

	var responses []dto.PesertaResponse
	for _, p := range pesertaList {
		responses = append(responses, *pesertaWithKelasToResponse(&p))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PesertaListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *pesertaService) UpdatePeserta(id string, req *dto.UpdatePesertaRequest) (*dto.PesertaResponse, error) {
	existing, err := s.repo.GetRawByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	if req.Username != existing.Username {
		taken, err := s.repo.GetByUsername(req.Username)
		if err == nil && taken != nil && taken.ID != id {
			return nil, errors.New("username sudah digunakan")
		}
	}

	peserta := &model.Peserta{
		ID:       existing.ID,
		Nama:     req.Nama,
		IDKelas:  req.IDKelas,
		Username: req.Username,
		Password: existing.Password,
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("gagal memproses password")
		}
		peserta.Password = hashedPassword
	}

	if err := s.repo.Update(peserta); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return pesertaWithKelasToResponse(updated), nil
}

func (s *pesertaService) DeletePeserta(id string) error {
	peserta, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}
	return s.repo.Delete(peserta.ID)
}

func (s *pesertaService) RestorePeserta(id string) error {
	return s.repo.Restore(id)
}

func (s *pesertaService) ImportPesertaFromExcel(ctx context.Context, req *dto.ImportPesertaRequest) (*dto.ImportPesertaResponse, error) {
	// 1. Validasi kelas exists
	exists, err := s.repo.GetKelasExists(ctx, req.IDKelas)
	if err != nil || !exists {
		return nil, errors.New("kelas tidak ditemukan")
	}

	// 2. Buka file dari request
	file, err := req.File.Open()
	if err != nil {
		return nil, errors.New("gagal membuka file")
	}
	defer file.Close()

	// 3. Parse excel file
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("file bukan format excel yang valid")
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, errors.New("gagal membaca sheet excel")
	}

	// 4. Parse & validasi tiap row (skip header, index 0)
	type pendingRow struct {
		row      *utils.ExcelPesertaRow
		username string
	}
	var pendings []pendingRow
	var errorDetails []dto.ImportPesertaErrorDetail
	var processedCount, failedCount int
	seenUsername := make(map[string]int) // username (lowercase) -> row pertama yang memakainya

	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]
		if len(row) == 0 {
			continue
		}
		processedCount++

		excelRow := utils.ParseExcelPesertaRow(row, rowIndex+1)

		validationErrors := utils.ValidatePesertaRow(excelRow)

		usernameKey := strings.ToLower(strings.TrimSpace(excelRow.Username))
		if usernameKey != "" {
			if firstRow, dup := seenUsername[usernameKey]; dup {
				validationErrors = append(validationErrors, "username duplikat dengan baris "+strconv.Itoa(firstRow))
			} else {
				seenUsername[usernameKey] = rowIndex + 1
			}
		}

		if len(validationErrors) > 0 {
			failedCount++
			errorDetails = append(errorDetails, dto.ImportPesertaErrorDetail{
				Row:   rowIndex + 1,
				Error: strings.Join(validationErrors, "; "),
			})
			continue
		}

		pendings = append(pendings, pendingRow{row: excelRow, username: excelRow.Username})
	}

	// 5. Cek username yang sudah dipakai peserta lain di database
	usernamesToCheck := make([]string, 0, len(pendings))
	for _, p := range pendings {
		usernamesToCheck = append(usernamesToCheck, p.username)
	}
	takenUsernames := make(map[string]bool)
	if len(usernamesToCheck) > 0 {
		taken, err := s.repo.GetByUsernames(usernamesToCheck)
		if err != nil {
			return nil, errors.New("gagal memeriksa username: " + err.Error())
		}
		for _, u := range taken {
			takenUsernames[strings.ToLower(u)] = true
		}
	}

	// 6. Hash password & siapkan baris final yang lolos semua validasi
	var pesertaList []model.Peserta
	successCount := 0
	for _, p := range pendings {
		if takenUsernames[strings.ToLower(p.username)] {
			failedCount++
			errorDetails = append(errorDetails, dto.ImportPesertaErrorDetail{
				Row:   p.row.RowIndex,
				Error: "username sudah digunakan",
			})
			continue
		}

		hashedPassword, err := utils.HashPassword(p.row.Password)
		if err != nil {
			failedCount++
			errorDetails = append(errorDetails, dto.ImportPesertaErrorDetail{
				Row:   p.row.RowIndex,
				Error: "gagal memproses password",
			})
			continue
		}

		pesertaList = append(pesertaList, model.Peserta{
			Nama:     p.row.Nama,
			IDKelas:  req.IDKelas,
			Username: p.username,
			Password: hashedPassword,
		})
		successCount++
	}

	// 7. Bulk insert ke database
	if len(pesertaList) > 0 {
		if err := s.repo.BulkCreate(ctx, pesertaList); err != nil {
			return nil, errors.New("gagal menyimpan data ke database: " + err.Error())
		}
	}

	// Limit error details ke max 100
	if len(errorDetails) > 100 {
		errorDetails = errorDetails[:100]
	}

	return &dto.ImportPesertaResponse{
		TotalProcessed: processedCount,
		TotalSuccess:   successCount,
		TotalFailed:    failedCount,
		IDKelas:        req.IDKelas,
		Timestamp:      time.Now(),
		Summary: map[string]int{
			"inserted": successCount,
			"errors":   failedCount,
		},
		Errors: errorDetails,
	}, nil
}
