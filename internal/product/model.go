package product

import (
	"admin/internal/brand"
	"admin/internal/category"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Price       int64  `gorm:"not null" json:"price"`
	Discount    int64  `gorm:"default:0" json:"discount"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// 🔹 Внешние ключи
	CategoryID uint              `gorm:"not null" json:"category_id"`
	Category   category.Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`

	BrandID uint        `gorm:"not null" json:"brand_id"`
	Brand   brand.Brand `gorm:"foreignKey:BrandID;constraint:OnDelete:CASCADE"`

	// 🔹 Дополнительные данные
	Images   string `gorm:"type:json" json:"images"`            // Храним ссылки на изображения JSON-массивом
	VideoURL string `gorm:"type:varchar(255)" json:"video_url"` // Видеообзор
}

//на поле discount,IsActive нужно делать слушателей (productVariants)

//для продукт варианта-//VendorCode   string  `gorm:"type:varchar(100);unique;not null"json:"vendor_code"`//артикул
//ReviewsCount int    `gorm:"default:0" json:"reviews_count"`     // Кол-во отзывов
//ProductGender gender = 12;
//ProductSeason season = 13;
