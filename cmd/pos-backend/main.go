package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/joshuaolumoye/pos-backend/internal/handler"
	"github.com/joshuaolumoye/pos-backend/internal/infrastructure"
	"github.com/joshuaolumoye/pos-backend/internal/middleware"
	"github.com/joshuaolumoye/pos-backend/internal/repository"
	"github.com/joshuaolumoye/pos-backend/internal/usecase"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger and JWT
	utils.InitLogger()
	utils.InitJWTKey()

	// Connect to DB with GORM
	db := infrastructure.InitializeDB()

	// Run auto-migrations (GORM handles everything!)
	if err := infrastructure.AutoMigrate(db); err != nil {
		utils.Logger.Fatal("Migration failed", utils.ZapError(err))
	}

	// Initialize cache (optional)
	_, err := infrastructure.NewBigCache()
	if err != nil {
		utils.Logger.Warn("Cache init failed", utils.ZapError(err))
	}

	// Dependency Injection
	businessRepo := &repository.BusinessRepo{DB: db}
	branchRepo := &repository.BranchRepo{DB: db}
	staffRepo := &repository.StaffRepo{DB: db}

	// Initialize sqlx DB for ProductRepo
	dsn := "" // Use the same DSN as GORM
	dbType := "mysql"
	if v := os.Getenv("DB_TYPE"); v != "" {
		dbType = v
	}
	switch dbType {
	case "mysql":
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "root"
		}
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "127.0.0.1"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "3306"
		}
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "pos"
		}
		dsn = user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	case "postgres":
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "127.0.0.1"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "pos"
		}
		dsn = "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	case "sqlite":
		dbname := os.Getenv("DB_NAME")
		if dbname == "" {
			dbname = "pos.db"
		}
		dsn = dbname
	}
	sqlxDB, err := infrastructure.NewDB(dsn)
	if err != nil {
		utils.Logger.Fatal("Failed to connect sqlx DB", utils.ZapError(err))
	}
	productRepo := &repository.ProductRepo{DB: sqlxDB}

	authRepo := &repository.AuthRepo{}

	businessUC := &usecase.BusinessUsecase{BusinessRepo: businessRepo, BranchRepo: branchRepo}
	branchUC := &usecase.BranchUsecase{BranchRepo: branchRepo}
	staffUC := &usecase.StaffUsecase{StaffRepo: staffRepo}
	productUC := &usecase.ProductUsecase{ProductRepo: productRepo}

	handler.BusinessUC = businessUC
	handler.BranchUC = branchUC
	handler.StaffUC = staffUC
	handler.ProductUC = productUC
	handler.AuthRepo = authRepo

	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	// Public endpoints
	r.Post("/api/business/register", handler.RegisterBusinessHandler)
	r.Post("/api/business/login", handler.LoginHandler)
	r.Post("/api/staff/login", handler.StaffLoginHandler)

	// Protected endpoints
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)

		// POST endpoints - NO query params allowed
		protected.With(middleware.NoQueryParamsMiddleware).Post("/api/branch/create", handler.CreateBranchHandler)
		protected.With(middleware.NoQueryParamsMiddleware).Post("/api/staff/create", handler.CreateStaffHandler)
		protected.With(middleware.NoQueryParamsMiddleware).Post("/api/product/add", handler.AddProductHandler)
		protected.With(middleware.NoQueryParamsMiddleware).Post("/api/sync", handler.SyncDataHandler)
		protected.With(middleware.NoQueryParamsMiddleware).Post("/api/auth/refresh", handler.RefreshTokenHandler)

		// GET endpoints - query params allowed
		protected.Get("/api/branches", handler.GetBranchesHandler)
		protected.Get("/api/staff", handler.GetStaffListHandler)
		protected.Get("/api/staff/details", handler.GetStaffByIDHandler)
		protected.Get("/api/products", handler.GetProductsHandler)
		protected.Get("/api/product/{id}", handler.GetProductHandler)

		// PUT/DELETE endpoints - query params allowed for resource IDs
		protected.Put("/api/branch/update", handler.UpdateBranchHandler)
		protected.Delete("/api/branch/delete", handler.DeleteBranchHandler)
		protected.Put("/api/staff/update", handler.UpdateStaffHandler)
		protected.Put("/api/product/{id}", handler.UpdateProductHandler)
		protected.Delete("/api/product/delete", handler.DeleteProductHandler)
	})

	utils.Logger.Info("POS Backend starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		utils.Logger.Fatal("Server failed", utils.ZapError(err))
	}
}
