package link

import (
	"admin/internal/stat"

	"gorm.io/gorm"
)

const (
	ErrCreateLink    = "Failed to create link: %v"
	ErrUpdateLink    = "Failed to update link: %v"
	ErrDeleteLink    = "Failed to delete link: %v"
	ErrDeleteForever = "Failed to permanently delete link: %v"
	ErrLinkNotFound  = "Link not found: %v"
	ErrGetLinks      = "Failed to retrieve links: %v"
)

type Link struct {
	gorm.Model `swaggerignore:"true"`
	Url        string      `json:"url"`
	Hash       string      `json:"hash" gorm:"uniqueIndex"`
	Stats      []stat.Stat `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	//поставили каскадную связь между таблицами что не позволит просто так удалить ссылку, так как она может относиться ко множеству статистик
	//ограничения некритичны
}
