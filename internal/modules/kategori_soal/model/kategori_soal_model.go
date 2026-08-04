package model

import (
	"database/sql"
	"time"
)

type KategoriSoal struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Kategori  string         `gorm:"type:varchar(255)" json:"kategori"`
	Benar     float64        `gorm:"type:numeric(5,2)" json:"benar"`
	Salah     float64        `gorm:"type:numeric(5,2)" json:"salah"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time     `gorm:"index" json:"deleted_at"`
	CreatedBy sql.NullString `gorm:"type:uuid" json:"created_by"`
	UpdatedBy sql.NullString `gorm:"type:uuid" json:"updated_by"`
}

func (KategoriSoal) TableName() string {
	return "kategori_soal"
}
