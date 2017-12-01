package models

import "github.com/jinzhu/gorm"

// Gallery is the image container resource
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	Create(gallery *Gallery)
}

type galleryGorm struct {
	db *gorm.DB
}
