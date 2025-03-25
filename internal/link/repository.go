package link

import (
	"admin/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(database *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: database,
	}
}
func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	result := repo.Database.DB.Create(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
	//для создания нам не нужно указывать таблицу линк потому что мы туда передаем структуру линк,и раз он имеет горм структуру, то создается он имеено в табличке линк
	// создание, получение результата по ссылкам. это всё как обертка над обычно db только с методами
}
func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.Database.DB.First(&link, "hash = ?", hash) // SQL QUERY BY CONDS
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}
func (repo *LinkRepository) Update(link *Link) (*Link, error) { // если поле в запросе не указано оно не обновляется и остается тем же
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) Delete(id uint, unscoped bool) error {
	query := repo.Database.DB
	if unscoped {
		query = query.Unscoped()
	}
	result := query.Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link                               // автоматическое lowercase и множественное число
	result := repo.Database.DB.First(&link, id) // SQL QUERY BY CONDS
	if result.Error != nil {
		return nil, result.Error
	}

	return &link, nil
}
func (repo *LinkRepository) Count(includeDeleted bool) int64 {
	var count int64
	query := repo.Database.Table("links")

	if !includeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	query.Count(&count)
	return count
}

func (repo *LinkRepository) GetAll(limit, offset int, includeDeleted bool) ([]Link, error) {
	var links []Link

	query := repo.Database.Table("links").Session(&gorm.Session{})

	if !includeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	err := query.
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&links).Error // Используем Find вместо Scan

	if err != nil {
		return nil, err
	}

	return links, nil
}
