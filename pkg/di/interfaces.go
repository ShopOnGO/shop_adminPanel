package di

import (
	"admin/internal/brand"
	"admin/internal/category"
	"admin/internal/product"
	pb "admin/pkg/service"
)

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserRepository interface {
	Create(user *pb.User) (*pb.User, error)
	FindByEmail(email string) (*pb.User, error)
	Update(user *pb.User) (*pb.User, error)
	Delete(id uint, unscoped bool) error
}

type IProductRepository interface {
	Create(product *product.Product) (*product.Product, error)
	GetByCategoryID(id uint) ([]product.Product, error)
	GetByName(name string) ([]product.Product, error)
	GetFeaturedProducts(amount uint, random, IncludeDeleted bool) ([]product.Product, error)
	Update(product *product.Product) (*product.Product, error)
	Delete(id uint, unscoped bool) error
}

type ICategoryRepository interface {
	Create(category *category.Category) (*category.Category, error) //done
	GetFeaturedCategories(amount int) ([]category.Category, error)
	GetFeaturedWithDeletedCategories(amount int) ([]category.Category, error)
	FindByName(name string) (*category.Category, error)
	FindCategoryByID(id uint) (*category.Category, error) //done
	Update(category *category.Category) (*category.Category, error)
	Delete(id uint) error
	DeleteForever(name string) error
}
type IBrandRepository interface {
	Create(category *brand.Brand) (*brand.Brand, error)
	GetFeaturedBrands(amount int) ([]brand.Brand, error)
	GetFeaturedWithDeletedBrands(amount int) ([]brand.Brand, error)
	FindByName(name string) (*brand.Brand, error)
	Update(brand *brand.Brand) (*brand.Brand, error)
	Delete(id uint) error
	DeleteForever(name string) error
}
