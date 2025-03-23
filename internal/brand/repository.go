package brand

import (
	"admin/pkg/db"
)

type BrandRepository struct {
	Database *db.Db
}

func NewBrandRepository(database *db.Db) *BrandRepository {
	return &BrandRepository{
		Database: database,
	}
}
func (repo *BrandRepository) Create(brand *Brand) (*Brand, error) {
	result := repo.Database.DB.Create(brand)
	if result.Error != nil {
		return nil, result.Error
	}
	return brand, nil
}
func (repo *BrandRepository) GetFeaturedBrands(amount int) ([]Brand, error) {
	var brand []Brand
	query := repo.Database.DB

	if amount > 0 {
		query = query.Limit(amount)
	}

	result := query.Find(&brand)
	if result.Error != nil {
		return nil, result.Error
	}

	return brand, nil
}

func (repo *BrandRepository) GetFeaturedWithDeletedBrands(amount int) ([]Brand, error) {
	var brand []Brand
	query := repo.Database.DB.Unscoped()

	if amount > 0 {
		query = query.Limit(amount)
	}

	result := query.Find(&brand)
	if result.Error != nil {
		return nil, result.Error
	}

	return brand, nil
}

func (repo *BrandRepository) FindBrandByID(id uint) (*Brand, error) {
	var brand Brand
	result := repo.Database.DB.First(&brand, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &brand, nil
}
func (repo *BrandRepository) FindByName(name string) (*Brand, error) {
	var brand Brand
	result := repo.Database.DB.First(&brand, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &brand, nil
}
func (repo *BrandRepository) Update(brand *Brand) (*Brand, error) {
	result := repo.Database.DB.Model(&Brand{}).Where("id = ?", brand.ID).Updates(brand)
	if result.Error != nil {
		return nil, result.Error
	}
	return brand, nil
}

func (repo *BrandRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Brand{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (repo *BrandRepository) DeleteForever(name string) error {
	// Принудительно удаляем все категории с таким именем
	// важно! find в методах поиска не даст удаленные мягко
	return repo.Database.DB.Unscoped().Where("name = ?", name).Delete(&Brand{}).Error
}
