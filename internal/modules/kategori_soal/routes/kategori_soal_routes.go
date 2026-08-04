package routes

import (
	"backend/internal/middleware"
	"backend/internal/modules/kategori_soal/controller"
	"backend/internal/modules/kategori_soal/repository"
	"backend/internal/modules/kategori_soal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupKategoriSoalRoutes(app *fiber.App, db *gorm.DB) {
	repo := repository.NewKategoriSoalRepository(db)
	svc := service.NewKategoriSoalService(repo)
	ctrl := controller.NewKategoriSoalController(svc)

	api := app.Group("/api")
	kategoriSoal := api.Group("/kategori-soal")

	kategoriSoal.Post("/", middleware.JWTAuth(), ctrl.CreateKategoriSoal)
	kategoriSoal.Get("/", ctrl.GetAllKategoriSoal)
	kategoriSoal.Get("/:id", ctrl.GetKategoriSoalByID)
	kategoriSoal.Put("/:id", middleware.JWTAuth(), ctrl.UpdateKategoriSoal)
	kategoriSoal.Delete("/:id", middleware.JWTAuth(), ctrl.DeleteKategoriSoal)
	kategoriSoal.Patch("/:id/restore", middleware.JWTAuth(), ctrl.RestoreKategoriSoal)
}
