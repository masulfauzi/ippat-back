package repository

import (
	"backend/internal/modules/kategori_soal/model"

	"gorm.io/gorm"
)

type KategoriSoalRepository interface {
	Create(kategoriSoal *model.KategoriSoal) error
	GetByID(id string) (*model.KategoriSoal, error)
	GetAll(page, pageSize int) ([]model.KategoriSoal, int64, error)
	Update(kategoriSoal *model.KategoriSoal) error
	Delete(id string) error
	Restore(id string) error
	HardDelete(id string) error
}

type kategoriSoalRepository struct {
	db *gorm.DB
}

func NewKategoriSoalRepository(db *gorm.DB) KategoriSoalRepository {
	return &kategoriSoalRepository{db: db}
}

func (r *kategoriSoalRepository) Create(kategoriSoal *model.KategoriSoal) error {
	return r.db.Create(kategoriSoal).Error
}

func (r *kategoriSoalRepository) GetByID(id string) (*model.KategoriSoal, error) {
	var kategoriSoal model.KategoriSoal
	err := r.db.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&kategoriSoal).Error
	if err != nil {
		return nil, err
	}
	return &kategoriSoal, nil
}

func (r *kategoriSoalRepository) GetAll(page, pageSize int) ([]model.KategoriSoal, int64, error) {
	var kategoriSoals []model.KategoriSoal
	var total int64

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Count total records (excluding soft deleted)
	err := r.db.
		Model(&model.KategoriSoal{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get records with pagination (excluding soft deleted)
	err = r.db.
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(pageSize).
		Find(&kategoriSoals).Error

	return kategoriSoals, total, err
}

func (r *kategoriSoalRepository) Update(kategoriSoal *model.KategoriSoal) error {
	return r.db.Save(kategoriSoal).Error
}

func (r *kategoriSoalRepository) Delete(id string) error {
	// Soft delete - GORM will automatically set deleted_at
	return r.db.Delete(&model.KategoriSoal{}, "id = ?", id).Error
}

func (r *kategoriSoalRepository) Restore(id string) error {
	// Restore - clear deleted_at using direct update
	return r.db.Table("kategori_soal").Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *kategoriSoalRepository) HardDelete(id string) error {
	return r.db.Unscoped().Delete(&model.KategoriSoal{}, "id = ?", id).Error
}
