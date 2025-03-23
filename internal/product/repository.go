package product

import (
	"admin/pkg/db"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.Database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) GetByCategoryID(id uint) ([]Product, error) { //limit 20
	var products []Product
	result := repo.Database.DB.
		Where("category_id = ?", id).
		Limit(20).
		Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

// func (repo *ProductRepository) GetByVendorCode(code *uuid)
func (repo *ProductRepository) GetByName(name string) ([]Product, error) {
	var products []Product
	result := repo.Database.DB.
		Where("name = ?", name).
		Limit(20).
		Find(&products) // table should be named "products"
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}

func (repo *ProductRepository) GetFeaturedProducts(amount uint, random, includeDeleted bool) ([]Product, error) {
	var products []Product
	query := repo.Database.DB.Table("products") // Указываем таблицу явно

	// Фильтр по удаленным
	if !includeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// Сортировка
	if random {
		query = query.Order("RANDOM()")
	} else {
		query = query.Order("popularity DESC")
	}

	// Выполняем запрос
	result := query.Limit(int(amount)).Find(&products)

	return products, result.Error
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.Database.DB.Model(&Product{}).Where("id = ?", product.ID).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) Delete(id uint, unscoped bool) error {
	query := repo.Database.DB.Where("id = ?", id)

	if unscoped {
		return query.Unscoped().Delete(&Product{}).Error
	}
	return query.Delete(&Product{}).Error
}
