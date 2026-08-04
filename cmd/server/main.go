package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/configs"
	_ "backend/docs"
	"backend/internal/database"
	"backend/internal/middleware"
	authroutes "backend/internal/modules/auth/routes"
	banksoalroutes "backend/internal/modules/bank_soal/routes"
	jadwalroutes "backend/internal/modules/jadwal/routes"
	jadwalkelasroutes "backend/internal/modules/jadwal_kelas/routes"
	jawabanroutes "backend/internal/modules/jawaban/routes"
	jurusanroutes "backend/internal/modules/jurusan/routes"
	kategorisoalroutes "backend/internal/modules/kategori_soal/routes"
	kelasroutes "backend/internal/modules/kelas/routes"
	mapelroutes "backend/internal/modules/mapel/routes"
	nilairoutes "backend/internal/modules/nilai/routes"
	pesertaroutes "backend/internal/modules/peserta/routes"
	soalroutes "backend/internal/modules/soal/routes"
	userroutes "backend/internal/modules/user/routes"
	"backend/internal/storage"

	"github.com/gofiber/fiber/v2"
	fiberswagger "github.com/swaggo/fiber-swagger"
)

// @title                       IPPAT Backend API
// @version                     1.0
// @description                 REST API untuk platform ujian online IPPAT (Fiber + GORM + PostgreSQL). Semua response dibungkus format {success, message, data, errors}.
// @description                 Endpoint yang butuh login (guru/admin/peserta) ditandai kunci gembok di UI ini — klik "Authorize" dan isi `Bearer <token>` yang didapat dari /api/auth/login atau /api/auth/register.
// @contact.name                IPPAT Backend Team
// @license.name                Proprietary
// @host                        localhost:3000
// @BasePath                    /api
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Ketik "Bearer" diikuti spasi dan JWT token. Contoh: "Bearer eyJhbGciOi..."

func main() {
	if err := configs.LoadEnv(); err != nil {
		log.Println("Warning: Error loading .env file:", err)
	}

	appConfig := configs.GetAppConfig()

	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.RunMigrations(database.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if err := storage.InitMinioClient(); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: appConfig.Name,
	})

	app.Use(middleware.CORS())
	app.Use(middleware.Logger())
	app.Use(middleware.Recovery())

	setupRoutes(app, appConfig)

	go func() {
		addr := fmt.Sprintf(":%d", appConfig.Port)
		log.Printf("Starting server on %s\n", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	_ = app.Shutdown()
	_ = database.Close()
	log.Println("Server shut down successfully")
}

func setupRoutes(app *fiber.App, appConfig *configs.AppConfig) {
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":    "ok",
			"service":   "Fiber Backend API",
			"server-no": appConfig.ServerNo,
		})
	})

	app.Static("/uploads", "./uploads")

	app.Get("/swagger/*", fiberswagger.WrapHandler)

	authroutes.SetupAuthRoutes(app, database.DB)
	userroutes.SetupUserRoutes(app, database.DB)
	mapelroutes.SetupMapelRoutes(app, database.DB)
	banksoalroutes.SetupBankSoalRoutes(app, database.DB)
	soalroutes.SetupSoalRoutes(app, database.DB)
	jurusanroutes.SetupJurusanRoutes(app, database.DB)
	kategorisoalroutes.SetupKategoriSoalRoutes(app, database.DB)
	kelasroutes.SetupKelasRoutes(app, database.DB)
	jadwalroutes.SetupJadwalRoutes(app, database.DB)
	jadwalkelasroutes.SetupJadwalKelasRoutes(app, database.DB)
	pesertaroutes.SetupPesertaRoutes(app, database.DB)
	nilairoutes.SetupNilaiRoutes(app, database.DB)
	jawabanroutes.SetupJawabanRoutes(app, database.DB)
}
