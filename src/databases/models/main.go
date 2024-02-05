package models

import "gorm.io/gorm"

func Sync(database *gorm.DB) {
	database.AutoMigrate(&Endpoints{})
	database.AutoMigrate(&Location{})
	database.AutoMigrate(&Owner{})
	database.AutoMigrate(&Company{})
	database.AutoMigrate(&Franchise{})

}
