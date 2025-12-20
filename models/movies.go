package models

import (
	"time"
)

type Movie struct {
	ID           uint          `gorm:"primaryKey;autoIncrement;not null"`
	Name         string        `json:"name" binding:"required" gorm:"index:idx_movie,unique;not null"`
	Director     string        `json:"director" binding:"required" gorm:"not null"`
	Producer     string        `json:"producer" binding:"required" gorm:"not null"`
	ReleaseDate  string        `json:"release_date" binding:"required" gorm:"index:idx_movie,unique;not null"`
	Duration     time.Duration `json:"duration" binding:"required" gorm:"not null"`
	BannerImagem string        `json:"banner_image" binding:"required,base64" gorm:"not null"`
}
