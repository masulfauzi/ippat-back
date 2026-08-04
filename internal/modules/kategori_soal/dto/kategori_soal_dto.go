package dto

type CreateKategoriSoalRequest struct {
	Kategori string  `json:"kategori" validate:"required,max=255"`
	Benar    float64 `json:"benar" validate:"required"`
	Salah    float64 `json:"salah" validate:"required"`
}

type UpdateKategoriSoalRequest struct {
	Kategori string  `json:"kategori" validate:"required,max=255"`
	Benar    float64 `json:"benar" validate:"required"`
	Salah    float64 `json:"salah" validate:"required"`
}

type KategoriSoalResponse struct {
	ID        string  `json:"id"`
	Kategori  string  `json:"kategori"`
	Benar     float64 `json:"benar"`
	Salah     float64 `json:"salah"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type KategoriSoalListResponse struct {
	Data      []KategoriSoalResponse `json:"data"`
	Total     int64                  `json:"total"`
	Page      int                    `json:"page"`
	PageSize  int                    `json:"page_size"`
	TotalPage int                    `json:"total_page"`
}
