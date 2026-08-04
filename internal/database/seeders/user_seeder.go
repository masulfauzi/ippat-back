package seeders

import (
	"backend/internal/modules/user/model"
	"backend/internal/utils"

	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&model.User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("admin#123")
	if err != nil {
		return err
	}

	admin := model.User{
		Name:     "Admin",
		Username: "admin",
		Email:    "admin@ippat.local",
		Password: hashedPassword,
		Role:     "admin",
		Status:   "active",
	}

	return db.Create(&admin).Error
}
