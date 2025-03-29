package migrations

import (
	"os"

	"admin/internal/brand"
	"admin/internal/category"
	"admin/internal/link"
	"admin/internal/product"
	"admin/internal/productVariant"
	"admin/internal/stat"
	"admin/internal/user"
	"admin/pkg/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CheckForMigrations() error {

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.Info("🚀 Starting migrations...")
		if err := RunMigrations(); err != nil {
			logger.Errorf("Error processing migrations: %v", err)
		}
		return nil
	}
	// if not "migrate" args[1]
	return nil
}

func RunMigrations() error {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		//DisableForeignKeyConstraintWhenMigrating: true, //временно игнорировать миграции в первый раз а потом их добавить
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{}, &product.Product{}, &category.Category{}, &brand.Brand{}, &productVariant.ProductVariant{})
	if err != nil {
		return err
	}

	logger.Info("✅")
	return nil
}
