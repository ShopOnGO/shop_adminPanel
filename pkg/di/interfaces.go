package di

import (
	"admin/internal/brand"
	"admin/internal/category"
	"admin/internal/product"
	"admin/internal/productVariant"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"
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
	GetFeaturedCategories(amount int, unscoped bool) ([]category.Category, error)
	FindByName(name string) (*category.Category, error)
	FindCategoryByID(id uint) (*category.Category, error) //done
	Update(category *category.Category) (*category.Category, error)
	Delete(name string, unscoped bool) error
}
type IBrandRepository interface {
	Create(category *brand.Brand) (*brand.Brand, error)
	GetFeaturedBrands(amount int, unscoped bool) ([]brand.Brand, error)
	FindByName(name string) (*brand.Brand, error)
	Update(brand *brand.Brand) (*brand.Brand, error)
	Delete(name string, unscoped bool) error
}

type ProductVariantRepositoryInterface interface {
	// Основные CRUD операции
	Create(variant *productVariant.ProductVariant) (*productVariant.ProductVariant, error)
	Update(variant *productVariant.ProductVariant) (*productVariant.ProductVariant, error)
	SoftDelete(id uint) error

	// Методы поиска
	GetBySKU(sku string) (*productVariant.ProductVariant, error)
	GetByBarcode(barcode string) (*productVariant.ProductVariant, error)
	GetByProductID(productID uint, includeInactive bool) ([]productVariant.ProductVariant, error)
	GetByFilters(filters map[string]interface{}, limit, offset int) ([]productVariant.ProductVariant, error)

	// Управление остатками
	UpdateStock(variantID uint, newStock uint32) error
	ReserveStock(variantID uint, quantity uint) error
	ReleaseStock(variantID uint, quantity uint) error
	GetAvailableStock(variantID uint) (uint32, error)
	BulkUpdateStock(variantStocks map[uint]uint32) error

	// Статусы и активность
	GetActive() ([]productVariant.ProductVariant, error)

	// Дополнительные методы
	GetFeaturedProducts(amount uint, random bool) ([]productVariant.ProductVariant, error) // Из исходного примера
}

// Дополнительные интерфейсы для конкретных реализаций
type StockManager interface {
	ReserveStock(variantID uint, quantity uint) error
	ReleaseStock(variantID uint, quantity uint) error
	GetAvailableStock(variantID uint) (uint32, error)
}

type SearchProvider interface {
	GetByFilters(filters map[string]interface{}, limit, offset int) ([]productVariant.ProductVariant, error)
	GetBySKU(sku string) (*productVariant.ProductVariant, error)
	GetByBarcode(barcode string) (*productVariant.ProductVariant, error)
}
