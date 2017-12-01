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
	Create(gallery *Gallery) error
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

// ensure that galleryGorm implements GalleryDB
var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

type gallerValidator struct {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
