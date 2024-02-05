package repositories

import (
	"clubhub/src/databases"
	"clubhub/src/databases/models"
)

type ownerInfo struct {
	Id         uint   `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	ZipCode    string `json:"zipCode"`
	LocationId uint   `json:"locationId"`
}

func CreateOwner(ownerInfo models.Owner) (uint, error) {
	result := databases.DB.Create(&ownerInfo)
	if result.Error != nil {
		return 0, result.Error
	}
	return ownerInfo.Id, nil
}
