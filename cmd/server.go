package main

import (
	"net"

	"admin/configs"
	"admin/internal/brand"
	"admin/internal/category"
	"admin/internal/home"
	"admin/internal/link"
	"admin/internal/product"
	"admin/internal/stat"
	"admin/internal/user"
	"admin/migrations"
	"admin/pkg/db"
	"admin/pkg/logger"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"

	"google.golang.org/grpc"
)

func AdminApp() *grpc.Server {

	conf := configs.LoadConfig()
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

	// services
	linkService := link.NewLinkService(linkRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{StatRepository: statRepository})
	homeService := home.NewHomeService(categoryRepository, productRepository, brandRepository)
	userService := user.NewUserService(userRepository)
	brandService := brand.NewBrandService(brandRepository)
	productService := product.NewProductServiceServer(productRepository)
	categoryService := category.NewCategoryService(categoryRepository)

	// registration
	pb.RegisterUserServiceServer(grpcServer, userService)
	pb.RegisterBrandServiceServer(grpcServer, brandService)
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	pb.RegisterHomeServiceServer(grpcServer, homeService)
	pb.RegisterLinkServiceServer(grpcServer, linkService)
	pb.RegisterProductServiceServer(grpcServer, productService)
	pb.RegisterStatServiceServer(grpcServer, statService)

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
