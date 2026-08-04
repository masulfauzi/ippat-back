package service

import (
	"errors"
	"math"

	"backend/internal/constants"
	"backend/internal/modules/kategori_soal/dto"
	"backend/internal/modules/kategori_soal/model"
	"backend/internal/modules/kategori_soal/repository"

	"gorm.io/gorm"
)

type KategoriSoalService interface {
	CreateKategoriSoal(req *dto.CreateKategoriSoalRequest) (*dto.KategoriSoalResponse, error)
	GetKategoriSoalByID(id string) (*dto.KategoriSoalResponse, error)
	GetAllKategoriSoal(page, pageSize int) (*dto.KategoriSoalListResponse, error)
	UpdateKategoriSoal(id string, req *dto.UpdateKategoriSoalRequest) (*dto.KategoriSoalResponse, error)
	DeleteKategoriSoal(id string) error
	RestoreKategoriSoal(id string) error
}

type kategoriSoalService struct {
	repo repository.KategoriSoalRepository
}

func NewKategoriSoalService(repo repository.KategoriSoalRepository) KategoriSoalService {
	return &kategoriSoalService{repo: repo}
}

func (s *kategoriSoalService) CreateKategoriSoal(req *dto.CreateKategoriSoalRequest) (*dto.KategoriSoalResponse, error) {
	kategoriSoal := &model.KategoriSoal{
		Kategori: req.Kategori,
		Benar:    req.Benar,
		Salah:    req.Salah,
	}

	if err := s.repo.Create(kategoriSoal); err != nil {
		return nil, err
	}

	return s.modelToResponse(kategoriSoal), nil
}

func (s *kategoriSoalService) GetKategoriSoalByID(id string) (*dto.KategoriSoalResponse, error) {
	kategoriSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	return s.modelToResponse(kategoriSoal), nil
}

func (s *kategoriSoalService) GetAllKategoriSoal(page, pageSize int) (*dto.KategoriSoalListResponse, error) {
	kategoriSoals, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var responses []dto.KategoriSoalResponse
	for _, kategoriSoal := range kategoriSoals {
		responses = append(responses, *s.modelToResponse(&kategoriSoal))
	}

	totalPage := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.KategoriSoalListResponse{
		Data:      responses,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *kategoriSoalService) UpdateKategoriSoal(id string, req *dto.UpdateKategoriSoalRequest) (*dto.KategoriSoalResponse, error) {
	kategoriSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(constants.ErrNotFound)
		}
		return nil, err
	}

	kategoriSoal.Kategori = req.Kategori
	kategoriSoal.Benar = req.Benar
	kategoriSoal.Salah = req.Salah

	if err := s.repo.Update(kategoriSoal); err != nil {
		return nil, err
	}

	return s.modelToResponse(kategoriSoal), nil
}

func (s *kategoriSoalService) DeleteKategoriSoal(id string) error {
	kategoriSoal, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(constants.ErrNotFound)
		}
		return err
	}

	return s.repo.Delete(kategoriSoal.ID)
}

func (s *kategoriSoalService) RestoreKategoriSoal(id string) error {
	return s.repo.Restore(id)
}

func (s *kategoriSoalService) modelToResponse(kategoriSoal *model.KategoriSoal) *dto.KategoriSoalResponse {
	return &dto.KategoriSoalResponse{
		ID:        kategoriSoal.ID,
		Kategori:  kategoriSoal.Kategori,
		Benar:     kategoriSoal.Benar,
		Salah:     kategoriSoal.Salah,
		CreatedAt: kategoriSoal.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: kategoriSoal.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
