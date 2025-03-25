package category

import (
	"admin/pkg/db"
)

type CategoryRepository struct {
	Database *db.Db
}

func NewCategoryRepository(database *db.Db) *CategoryRepository {
	return &CategoryRepository{
		Database: database,
	}
}

func (repo *CategoryRepository) Create(category *Category) (*Category, error) {
	result := repo.Database.DB.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (repo *CategoryRepository) GetFeaturedCategories(amount int, unscoped bool) ([]Category, error) {
	var categories []Category
	query := repo.Database.DB
	if unscoped {
		query = query.Unscoped()
	}
	if amount > 0 {
		query = query.Limit(amount)
	}

	result := query.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (repo *CategoryRepository) FindByName(name string) (*Category, error) {
	var category Category
	result := repo.Database.DB.First(&category, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}
func (repo *CategoryRepository) Update(category *Category) (*Category, error) {
	result := repo.Database.DB.Model(&Category{}).Where("id = ?", category.ID).Updates(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (repo *CategoryRepository) FindCategoryByID(id uint) (*Category, error) {
	var category Category
	result := repo.Database.DB.First(&category, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

//	func (repo *CategoryRepository) Delete(id uint) error {
//		result := repo.Database.DB.Delete(&Category{}, id)
//		if result.Error != nil {
//			return result.Error
//		}
//		return nil
//	}
func (repo *CategoryRepository) Delete(name string, unscoped bool) error {
	// Принудительно удаляем все категории с таким именем
	// важно! find в методах поиска не даст удаленные мягко

	query := repo.Database.DB
	if unscoped {
		query = query.Unscoped()
	}
	return query.Where("name = ?", name).Delete(&Category{}).Error
}
