package repositories

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
)

type ownerInfo struct {
	Id         uint
	FirstName  string
	LastName   string
	Email      string
	Phone      string
	Address    string
	ZipCode    string
	LocationId uint
}

func CreateOwner(ownerInfo models.Owner) (uint, error) {
	result := databases.DB.Create(&ownerInfo)
	if result.Error != nil {
		return 0, result.Error
	}
	return ownerInfo.Id, nil
}
