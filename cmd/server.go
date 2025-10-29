package main

import (
	"log"
	"net"

	"admin/configs"
	"admin/internal/brand"
	"admin/internal/category"
	"admin/internal/home"
	"admin/internal/link"
	"admin/internal/product"
	"admin/internal/productVariant"
	"admin/internal/stat"
	"admin/internal/user"
	"admin/migrations"
	"admin/pkg/db"
	"admin/pkg/dlq"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"

	"google.golang.org/grpc"
)

func AdminApp() *grpc.Server {

	conf := configs.LoadConfig()
	consoleLvl := conf.LogLevel
	fileLvl := conf.FileLogLevel
	logger.InitLogger(consoleLvl, fileLvl)
	logger.EnableFileLogging("TailorNado_admin-service")
	db := db.NewDB(conf)

	// Создаем новый gRPC-сервер
	grpcServer := grpc.NewServer()

	// repositories
	statRepository := stat.NewStatRepository(db)
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	brandRepository := brand.NewBrandRepository(db)
	productRepository := product.NewProductRepository(db)
	categoryRepository := category.NewCategoryRepository(db)
	productVariantRepository := productVariant.NewProductVariantRepository(db)

	//validators
	validator := &productVariant.ProductVariantValidator{}

	// services
	linkService := link.NewLinkService(linkRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{StatRepository: statRepository})
	homeService := home.NewHomeService(categoryRepository, productRepository, brandRepository)
	userService := user.NewUserService(userRepository)
	brandService := brand.NewBrandService(brandRepository)
	productService := product.NewProductServiceServer(productRepository)
	categoryService := category.NewCategoryService(categoryRepository)
	productVariantService := productVariant.NewVariantService(productVariantRepository, validator)

	// registration
	pb.RegisterUserServiceServer(grpcServer, userService)
	pb.RegisterBrandServiceServer(grpcServer, brandService)
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	pb.RegisterHomeServiceServer(grpcServer, homeService)
	pb.RegisterLinkServiceServer(grpcServer, linkService)
	pb.RegisterProductServiceServer(grpcServer, productService)
	pb.RegisterStatServiceServer(grpcServer, statService)
	pb.RegisterProductVariantServiceServer(grpcServer, productVariantService)

	log.Println("🚀 Запуск DLQ процессора...")
	dlq.StartDLQProcessor(conf)

	return grpcServer
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Errorf("Error due conn to tcp: %v", err)
		return
	}
	migrations.CheckForMigrations()
	logger.Info("gRPC server is running on :50051")

	grpcServer := AdminApp()

	if err := grpcServer.Serve(listener); err != nil {
		logger.Errorf("Error due starting the gRPC server: %v", err)
	}
}
